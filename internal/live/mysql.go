package live

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func analyzeMySQL(parent context.Context, cfg MySQLConfig, timeout time.Duration) LiveSection {
	ctx, cancel := context.WithTimeout(parent, timeout)
	defer cancel()

	section := LiveSection{
		ID:     "mysql-live",
		Name:   "MySQL 实时分析",
		Status: "error",
	}

	db, err := sql.Open("mysql", cfg.DSN)
	if err != nil {
		section.Error = err.Error()
		section.Summary = "MySQL 连接创建失败。"
		return section
	}
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		section.Error = err.Error()
		section.Summary = "MySQL 连通性检查失败。"
		return section
	}

	schema := cfg.Schema
	if schema == "" {
		if detected, err := queryString(ctx, db, "SELECT DATABASE()"); err == nil {
			schema = detected
		}
	}

	version, _ := queryString(ctx, db, "SELECT VERSION()")
	section.Connected = true
	section.Status = "ok"
	section.Summary = "已连接 MySQL，下面的内容来自真实实例查询而非静态演示。"
	section.Metrics = append(section.Metrics,
		Metric{Label: "Version", Value: valueOr(version, "unknown")},
		Metric{Label: "Schema", Value: valueOr(schema, "(未选择)")},
	)

	engineTable, err := queryRows(ctx, db, `
		SELECT ENGINE, SUPPORT, TRANSACTIONS, XA, SAVEPOINTS
		FROM information_schema.ENGINES
		WHERE ENGINE IN ('InnoDB', 'MyISAM', 'MEMORY')
		ORDER BY ENGINE`)
	if err == nil {
		engineTable.Title = "存储引擎支持情况"
		section.Tables = append(section.Tables, engineTable)
	} else {
		section.Warnings = append(section.Warnings, "未能读取 information_schema.ENGINES: "+err.Error())
	}

	if schema != "" && cfg.Table != "" {
		tableMeta, err := queryRows(ctx, db, `
			SELECT TABLE_NAME, ENGINE, TABLE_ROWS, CREATE_TIME, UPDATE_TIME
			FROM information_schema.TABLES
			WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?`, schema, cfg.Table)
		if err == nil {
			tableMeta.Title = "目标表元数据"
			section.Tables = append(section.Tables, tableMeta)
		} else {
			section.Warnings = append(section.Warnings, "未能读取目标表元数据: "+err.Error())
		}

		indexes, err := queryRows(ctx, db, `
			SELECT INDEX_NAME, NON_UNIQUE, SEQ_IN_INDEX, COLUMN_NAME, INDEX_TYPE
			FROM information_schema.STATISTICS
			WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?
			ORDER BY INDEX_NAME, SEQ_IN_INDEX`, schema, cfg.Table)
		if err == nil {
			indexes.Title = "目标表索引"
			section.Tables = append(section.Tables, indexes)
		} else {
			section.Warnings = append(section.Warnings, "未能读取目标表索引: "+err.Error())
		}
	}

	lockQuery := `
		SELECT
			dl.OBJECT_SCHEMA,
			dl.OBJECT_NAME,
			dl.INDEX_NAME,
			dl.LOCK_TYPE,
			dl.LOCK_MODE,
			dl.LOCK_STATUS,
			COALESCE(dw.REQUESTING_ENGINE_LOCK_ID, '') AS REQUESTING_LOCK,
			COALESCE(dw.BLOCKING_ENGINE_LOCK_ID, '') AS BLOCKING_LOCK
		FROM performance_schema.data_locks dl
		LEFT JOIN performance_schema.data_lock_waits dw
			ON dl.ENGINE_LOCK_ID = dw.REQUESTING_ENGINE_LOCK_ID
			OR dl.ENGINE_LOCK_ID = dw.BLOCKING_ENGINE_LOCK_ID`
	args := []any{}
	if schema != "" {
		lockQuery += " WHERE dl.OBJECT_SCHEMA = ?"
		args = append(args, schema)
	}
	lockQuery += " ORDER BY dl.OBJECT_NAME, dl.INDEX_NAME, dl.LOCK_MODE"
	locks, err := queryRows(ctx, db, lockQuery, args...)
	if err == nil {
		locks.Title = "当前活动锁"
		if len(locks.Rows) == 0 {
			locks.Rows = [][]string{{"(none)", "(none)", "(none)", "(none)", "(none)", "当前未观察到活动记录锁", "", ""}}
		}
		section.Tables = append(section.Tables, locks)
	} else {
		section.Warnings = append(section.Warnings, "未能读取 performance_schema.data_locks，通常是权限不足或实例未开启相关采集: "+err.Error())
	}

	if strings.TrimSpace(cfg.ExplainQuery) != "" {
		var explain sql.NullString
		err := db.QueryRowContext(ctx, "EXPLAIN FORMAT=JSON "+cfg.ExplainQuery).Scan(&explain)
		if err == nil && explain.String != "" {
			section.Snippets = append(section.Snippets, Snippet{
				Title:    "EXPLAIN FORMAT=JSON",
				Language: "json",
				Content:  explain.String,
			})
		} else if err != nil {
			section.Warnings = append(section.Warnings, "未能执行 EXPLAIN FORMAT=JSON: "+err.Error())
		}
	}

	section.Metrics = append(section.Metrics,
		Metric{Label: "Lock View", Value: "performance_schema.data_locks", Hint: "用于读取实例当前锁状态"},
		Metric{Label: "Index Source", Value: "information_schema.STATISTICS", Hint: "索引信息来自真实元数据"},
	)

	if cfg.Table != "" {
		section.Summary = fmt.Sprintf("已连接 MySQL，正在分析 `%s.%s` 的引擎、索引和当前锁视图。", schema, cfg.Table)
	}

	return section
}

func valueOr(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}
