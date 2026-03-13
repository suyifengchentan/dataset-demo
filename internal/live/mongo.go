package live

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func analyzeMongo(parent context.Context, cfg MongoConfig, timeout time.Duration) LiveSection {
	ctx, cancel := context.WithTimeout(parent, timeout)
	defer cancel()

	section := LiveSection{
		ID:     "mongodb-live",
		Name:   "MongoDB 实时分析",
		Status: "error",
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.URI))
	if err != nil {
		section.Error = err.Error()
		section.Summary = "MongoDB 连接创建失败。"
		return section
	}
	defer func() {
		disconnectCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = client.Disconnect(disconnectCtx)
	}()

	if err := client.Ping(ctx, readpref.PrimaryPreferred()); err != nil {
		section.Error = err.Error()
		section.Summary = "MongoDB 连通性检查失败。"
		return section
	}

	dbName := cfg.Database
	if dbName == "" {
		dbName = "admin"
	}
	db := client.Database(dbName)

	buildInfo := bson.M{}
	if err := db.RunCommand(ctx, bson.D{{Key: "buildInfo", Value: 1}}).Decode(&buildInfo); err != nil {
		section.Warnings = append(section.Warnings, "未能读取 buildInfo: "+err.Error())
	}
	hello := bson.M{}
	if err := db.RunCommand(ctx, bson.D{{Key: "hello", Value: 1}}).Decode(&hello); err != nil {
		section.Warnings = append(section.Warnings, "未能读取 hello: "+err.Error())
	}

	section.Connected = true
	section.Status = "ok"
	section.Summary = "已连接 MongoDB，展示真实实例的文档集合、复制角色和示例文档。"
	section.Metrics = append(section.Metrics,
		Metric{Label: "Version", Value: stringify(buildInfo["version"])},
		Metric{Label: "Database", Value: dbName},
		Metric{Label: "Writable Primary", Value: stringify(hello["isWritablePrimary"])},
		Metric{Label: "Replica Set", Value: valueOr(stringify(hello["setName"]), "(standalone)")},
	)

	names, err := db.ListCollectionNames(ctx, bson.D{})
	if err == nil {
		rows := make([][]string, 0, len(names))
		for _, name := range names {
			rows = append(rows, []string{name})
		}
		section.Tables = append(section.Tables, DataTable{
			Title:   "集合列表",
			Columns: []string{"collection"},
			Rows:    rows,
		})
	} else {
		section.Warnings = append(section.Warnings, "未能读取集合列表: "+err.Error())
	}

	if cfg.Collection != "" {
		coll := db.Collection(cfg.Collection)

		var stats bson.M
		if err := db.RunCommand(ctx, bson.D{{Key: "collStats", Value: cfg.Collection}}).Decode(&stats); err == nil {
			section.Tables = append(section.Tables, DataTable{
				Title:   "集合统计",
				Columns: []string{"字段", "值"},
				Rows: [][]string{
					{"count", stringify(stats["count"])},
					{"size", stringify(stats["size"])},
					{"avgObjSize", stringify(stats["avgObjSize"])},
					{"storageSize", stringify(stats["storageSize"])},
					{"nindexes", stringify(stats["nindexes"])},
				},
			})
		} else {
			section.Warnings = append(section.Warnings, "未能读取 collStats: "+err.Error())
		}

		var sample bson.M
		if err := coll.FindOne(ctx, bson.D{}).Decode(&sample); err == nil {
			section.Snippets = append(section.Snippets, Snippet{
				Title:    "示例文档",
				Language: "json",
				Content:  extJSON(sample),
			})
		} else {
			section.Warnings = append(section.Warnings, "未能读取示例文档: "+err.Error())
		}

		cursor, err := coll.Aggregate(ctx, mongo.Pipeline{
			bson.D{{Key: "$limit", Value: 5}},
			bson.D{{Key: "$project", Value: bson.D{{Key: "_id", Value: 1}, {Key: "status", Value: 1}, {Key: "tags", Value: 1}}}},
		})
		if err == nil {
			defer cursor.Close(ctx)
			rows := [][]string{}
			for cursor.Next(ctx) {
				doc := bson.M{}
				if err := cursor.Decode(&doc); err != nil {
					continue
				}
				rows = append(rows, []string{
					stringify(doc["_id"]),
					stringify(doc["status"]),
					stringify(doc["tags"]),
				})
			}
			if len(rows) > 0 {
				section.Tables = append(section.Tables, DataTable{
					Title:   "聚合管道结果摘录",
					Columns: []string{"_id", "status", "tags"},
					Rows:    rows,
				})
			}
		} else {
			section.Warnings = append(section.Warnings, "未能执行示例聚合: "+err.Error())
		}

		section.Summary = fmt.Sprintf("已连接 MongoDB，正在分析 `%s.%s` 的文档与集合统计。", dbName, cfg.Collection)
	}

	section.Snippets = append(section.Snippets, Snippet{
		Title:    "hello 返回摘录",
		Language: "json",
		Content:  extJSON(hello),
	})

	return section
}

func extJSON(value any) string {
	data, err := bson.MarshalExtJSON(value, true, true)
	if err == nil {
		return string(data)
	}

	converted := normalizeMongoValue(value)
	return compactJSON(converted)
}

func normalizeMongoValue(value any) any {
	switch v := value.(type) {
	case primitive.DateTime:
		return v.Time().Format(time.RFC3339)
	default:
		return value
	}
}
