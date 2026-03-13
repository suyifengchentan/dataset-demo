# 数据库核心机制可视化实验台

这是一个 Go 实现的数据库教学与诊断项目，支持两种使用方式：

- 静态教学模式：直接展示 MySQL、PostgreSQL、Redis、MongoDB 的核心机制可视化
- 实时数据库模式：连接真实数据库实例，读取元数据、运行统计、聚合结果和诊断信息

## 已实现能力

### MySQL

- 读取真实实例版本、当前 schema、目标表元数据
- 读取真实索引信息 `information_schema.STATISTICS`
- 读取存储引擎支持情况 `information_schema.ENGINES`
- 尝试读取当前活动锁 `performance_schema.data_locks`
- 支持用户输入 `EXPLAIN FORMAT=JSON` 查询

### PostgreSQL

- 读取版本、数据库名、`wal_level`
- 展示 JSONB、数组、范围类型等原生能力
- 读取扩展信息 `pg_available_extensions`
- 读取逻辑复制 publication
- 读取 `pg_stat_user_tables` 中的 MVCC / vacuum 统计

### Redis

- 读取真实实例版本、内存、OPS、客户端数量
- 读取 IO 多路复用实现 `multiplexing_api`
- 执行真实原子性探针：`INCR` + Lua `INCRBY`
- 展示命令统计和事件循环相关信息

### MongoDB

- 读取真实实例版本、复制角色、集合列表
- 读取集合统计 `collStats`
- 展示真实样例文档
- 执行示例聚合管道并展示结果摘录

## 运行项目

```bash
go run .
```

默认监听：

```text
http://localhost:8080
```

## 使用自己的数据库

启动服务后，在页面顶部填写连接参数：

- MySQL 使用 DSN，例如 `root:root@tcp(127.0.0.1:3306)/visual_lab`
- PostgreSQL 使用 DSN，例如 `postgres://postgres:postgres@127.0.0.1:5432/visual_lab?sslmode=disable`
- Redis 使用 `host:port`
- MongoDB 使用 URI，例如 `mongodb://127.0.0.1:27017`

如果某个数据库暂时不需要连接，可以取消对应的“启用”复选框。

## 一键启动示例环境

项目提供了一套 Docker Compose 样例环境，包含四类数据库和初始化数据：

```bash
docker compose up -d
```

启动后，前端点击“填充示例环境”即可自动写入默认连接参数。

### 示例环境端口

- MySQL: `127.0.0.1:3306`
- PostgreSQL: `127.0.0.1:5432`
- Redis: `127.0.0.1:6379`
- MongoDB: `127.0.0.1:27017`

### 示例环境默认账号

- MySQL: `root / root`
- PostgreSQL: `postgres / postgres`
- Redis: 无密码
- MongoDB: 默认未启用鉴权

## 项目结构

```text
.
├── docker-compose.yml
├── docker/
│   ├── mongo/init/01_seed.js
│   ├── mysql/init/01_schema.sql
│   └── postgres/init/01_schema.sql
├── internal/
│   ├── demo/
│   │   └── data.go
│   └── live/
│       ├── analyzer.go
│       ├── helpers.go
│       ├── mongo.go
│       ├── mysql.go
│       ├── postgres.go
│       ├── redis.go
│       └── types.go
├── main.go
└── static/
    ├── app.js
    ├── index.html
    └── style.css
```

## 设计判断

这个项目现在已经从“纯讲解页面”升级为“半教学、半观测”的工具，但仍然有几个边界需要明确：

- MySQL 的物理 B+Tree 页结构通常不能通过通用 SQL 直接完整读取，所以当前展示的是基于真实索引元数据和 `EXPLAIN` 的逻辑访问视图，而不是存储页十六进制解析
- Redis 的 `multiplexing_api` 取决于实例运行平台；如果你在 macOS 本机运行 Redis，页面看到的可能是 `kqueue` 而不是 `epoll`
- MongoDB 的聚合示例是“可运行的通用示例”，它不一定正好等于你的业务分析口径
- PostgreSQL 的 MVCC 机制本质上比简单统计表更复杂，当前页面更偏向观测和解释，而不是事务级调试器

## 后续可继续扩展

- 增加数据库连接池状态、慢查询、锁等待时间线
- 为 MySQL 和 PostgreSQL 增加事务实验按钮
- 为 Redis 增加 `MULTI/EXEC` 与 Lua 的对照实验
- 为 MongoDB 增加更智能的 schema 推断和 pipeline 生成器
