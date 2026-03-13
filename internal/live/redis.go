package live

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func analyzeRedis(parent context.Context, cfg RedisConfig, timeout time.Duration) LiveSection {
	ctx, cancel := context.WithTimeout(parent, timeout)
	defer cancel()

	section := LiveSection{
		ID:     "redis-live",
		Name:   "Redis 实时分析",
		Status: "error",
	}

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Username: cfg.Username,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	defer client.Close()

	if err := client.Ping(ctx).Err(); err != nil {
		section.Error = err.Error()
		section.Summary = "Redis 连通性检查失败。"
		return section
	}

	serverInfo, _ := client.Info(ctx, "server").Result()
	statsInfo, _ := client.Info(ctx, "stats").Result()
	memoryInfo, _ := client.Info(ctx, "memory").Result()
	clientsInfo, _ := client.Info(ctx, "clients").Result()
	commandInfo, _ := client.Info(ctx, "commandstats").Result()

	server := parseRedisInfo(serverInfo)
	stats := parseRedisInfo(statsInfo)
	memory := parseRedisInfo(memoryInfo)
	clients := parseRedisInfo(clientsInfo)
	commandStats := parseRedisInfo(commandInfo)

	section.Connected = true
	section.Status = "ok"
	section.Summary = "已连接 Redis，原子性探针和 IO 多路复用信息都来自真实实例。"
	section.Metrics = append(section.Metrics,
		Metric{Label: "Version", Value: valueOr(server["redis_version"], "unknown")},
		Metric{Label: "Multiplexing API", Value: valueOr(server["multiplexing_api"], "unknown"), Hint: "Linux 上通常为 epoll，macOS 常见为 kqueue"},
		Metric{Label: "Ops/s", Value: valueOr(stats["instantaneous_ops_per_sec"], "0")},
		Metric{Label: "Connected Clients", Value: valueOr(clients["connected_clients"], "0")},
		Metric{Label: "Memory", Value: valueOr(memory["used_memory_human"], valueOr(memory["used_memory"], "unknown"))},
	)

	prefix := cfg.KeyPrefix
	if prefix == "" {
		prefix = "dbdemo:atomic"
	}
	probeKey := fmt.Sprintf("%s:%d", prefix, time.Now().UnixNano())
	incrValue, incrErr := client.Incr(ctx, probeKey).Result()
	luaScript := redis.NewScript(`
local value = redis.call("INCRBY", KEYS[1], ARGV[1])
redis.call("PEXPIRE", KEYS[1], ARGV[2])
return value`)
	luaValue, luaErr := luaScript.Run(ctx, client, []string{probeKey}, 5, 120000).Result()
	finalValue, getErr := client.Get(ctx, probeKey).Result()
	_ = client.Expire(ctx, probeKey, 2*time.Minute).Err()

	probeRows := [][]string{
		{"INCR", stringify(incrValue), valueOfErr(incrErr)},
		{"Lua INCRBY", stringify(luaValue), valueOfErr(luaErr)},
		{"GET final", finalValue, valueOfErr(getErr)},
	}
	section.Tables = append(section.Tables, DataTable{
		Title:   "原子性探针",
		Columns: []string{"步骤", "返回值", "错误"},
		Rows:    probeRows,
	})

	eventsTable := DataTable{
		Title:   "事件循环与运行时概况",
		Columns: []string{"字段", "值"},
		Rows: [][]string{
			{"multiplexing_api", valueOr(server["multiplexing_api"], "unknown")},
			{"io_threads_active", valueOr(server["io_threads_active"], "unknown")},
			{"total_commands_processed", valueOr(stats["total_commands_processed"], "unknown")},
			{"rejected_connections", valueOr(stats["rejected_connections"], "0")},
		},
	}
	section.Tables = append(section.Tables, eventsTable)

	topCommands := make([][]string, 0, 6)
	for _, key := range sortedKeys(commandStats) {
		if len(topCommands) >= 6 {
			break
		}
		topCommands = append(topCommands, []string{key, commandStats[key]})
	}
	if len(topCommands) > 0 {
		section.Tables = append(section.Tables, DataTable{
			Title:   "命令统计摘录",
			Columns: []string{"命令", "统计"},
			Rows:    topCommands,
		})
	}

	section.Snippets = append(section.Snippets, Snippet{
		Title:    "Redis 实例摘录",
		Language: "ini",
		Content:  prettyMapLines(server, "redis_version", "os", "arch_bits", "multiplexing_api", "io_threads_active"),
	})

	return section
}

func valueOfErr(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
