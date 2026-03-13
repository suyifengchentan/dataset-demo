package live

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func analyzePostgres(parent context.Context, cfg PostgresConfig, timeout time.Duration) LiveSection {
	ctx, cancel := context.WithTimeout(parent, timeout)
	defer cancel()

	section := LiveSection{
		ID:     "postgres-live",
		Name:   "PostgreSQL 实时分析",
		Status: "error",
	}

	db, err := sql.Open("pgx", cfg.DSN)
	if err != nil {
		section.Error = err.Error()
		section.Summary = "PostgreSQL 连接创建失败。"
		return section
	}
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		section.Error = err.Error()
		section.Summary = "PostgreSQL 连通性检查失败。"
		return section
	}

	version, _ := queryString(ctx, db, "SHOW server_version")
	currentDB, _ := queryString(ctx, db, "SELECT current_database()")
	walLevel, _ := queryString(ctx, db, "SHOW wal_level")
	section.Connected = true
	section.Status = "ok"
	section.Summary = "已连接 PostgreSQL，输出聚焦 MVCC、扩展能力和分析特性。"
	section.Metrics = append(section.Metrics,
		Metric{Label: "Version", Value: valueOr(version, "unknown")},
		Metric{Label: "Database", Value: valueOr(currentDB, "(unknown)")},
		Metric{Label: "wal_level", Value: valueOr(walLevel, "(unknown)")},
	)

	featureTable, err := queryRows(ctx, db, `
		SELECT feature, availability
		FROM (
			SELECT 'jsonb' AS feature, EXISTS (SELECT 1 FROM pg_type WHERE typname = 'jsonb')::text AS availability
			UNION ALL
			SELECT 'array', EXISTS (SELECT 1 FROM pg_type WHERE typname = '_text')::text
			UNION ALL
			SELECT 'range', EXISTS (SELECT 1 FROM pg_type WHERE typname = 'int4range')::text
			UNION ALL
			SELECT 'window_function', 'true'
		) t`)
	if err == nil {
		featureTable.Title = "区别于传统关系库的原生能力"
		section.Tables = append(section.Tables, featureTable)
	} else {
		section.Warnings = append(section.Warnings, "未能读取能力概览: "+err.Error())
	}

	extensions, err := queryRows(ctx, db, `
		SELECT name, default_version, COALESCE(installed_version, '(not installed)') AS installed_version
		FROM pg_available_extensions
		WHERE name IN ('pg_trgm', 'uuid-ossp', 'postgis')
		ORDER BY name`)
	if err == nil {
		extensions.Title = "扩展生态"
		section.Tables = append(section.Tables, extensions)
	} else {
		section.Warnings = append(section.Warnings, "未能读取扩展列表: "+err.Error())
	}

	publications, err := queryRows(ctx, db, `
		SELECT pubname, puballtables, pubinsert, pubupdate, pubdelete
		FROM pg_publication
		ORDER BY pubname`)
	if err == nil {
		publications.Title = "逻辑复制 Publication"
		if len(publications.Rows) > 0 {
			section.Tables = append(section.Tables, publications)
		}
	} else {
		section.Warnings = append(section.Warnings, "未能读取逻辑复制 publication: "+err.Error())
	}

	if cfg.Table != "" {
		schema := cfg.Schema
		if schema == "" {
			schema = "public"
		}

		mvccTable, err := queryRows(ctx, db, `
			SELECT relname, n_live_tup, n_dead_tup, n_tup_ins, n_tup_upd, n_tup_del, vacuum_count, autovacuum_count
			FROM pg_stat_user_tables
			WHERE schemaname = $1 AND relname = $2`, schema, cfg.Table)
		if err == nil {
			mvccTable.Title = "MVCC / Vacuum 统计"
			section.Tables = append(section.Tables, mvccTable)
		} else {
			section.Warnings = append(section.Warnings, "未能读取 pg_stat_user_tables: "+err.Error())
		}

		definition, err := queryRows(ctx, db, `
			SELECT column_name, data_type, is_nullable
			FROM information_schema.columns
			WHERE table_schema = $1 AND table_name = $2
			ORDER BY ordinal_position`, schema, cfg.Table)
		if err == nil {
			definition.Title = "目标表结构"
			section.Tables = append(section.Tables, definition)
		} else {
			section.Warnings = append(section.Warnings, "未能读取目标表结构: "+err.Error())
		}

		section.Summary = fmt.Sprintf("已连接 PostgreSQL，正在分析 `%s.%s` 的 MVCC 统计和扩展能力。", schema, cfg.Table)
	}

	settings, err := queryRows(ctx, db, `
		SELECT name, setting
		FROM pg_settings
		WHERE name IN ('shared_buffers', 'max_connections', 'max_worker_processes')
		ORDER BY name`)
	if err == nil {
		settings.Title = "关键运行参数"
		section.Tables = append(section.Tables, settings)
	}

	section.Snippets = append(section.Snippets, Snippet{
		Title:    "推荐观测 SQL",
		Language: "sql",
		Content:  "SELECT relname, n_live_tup, n_dead_tup FROM pg_stat_user_tables ORDER BY n_dead_tup DESC LIMIT 10;",
	})

	return section
}
