package demo

type DemoPage struct {
	Title    string        `json:"title"`
	Subtitle string        `json:"subtitle"`
	Sections []DemoSection `json:"sections"`
}

type DemoSection struct {
	ID             string          `json:"id"`
	Name           string          `json:"name"`
	Overview       string          `json:"overview"`
	Highlights     []string        `json:"highlights"`
	Visualizations []Visualization `json:"visualizations"`
	Comparisons    []ComparisonRow `json:"comparisons,omitempty"`
	Callouts       []Callout       `json:"callouts,omitempty"`
}

type Visualization struct {
	Type        string            `json:"type"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Items       []VisualItem      `json:"items,omitempty"`
	Steps       []StepItem        `json:"steps,omitempty"`
	Columns     []string          `json:"columns,omitempty"`
	Rows        [][]string        `json:"rows,omitempty"`
	Meta        map[string]string `json:"meta,omitempty"`
}

type VisualItem struct {
	Title    string   `json:"title"`
	Label    string   `json:"label,omitempty"`
	Detail   string   `json:"detail,omitempty"`
	Accent   string   `json:"accent,omitempty"`
	Children []string `json:"children,omitempty"`
}

type StepItem struct {
	Title  string `json:"title"`
	Detail string `json:"detail"`
	State  string `json:"state,omitempty"`
}

type ComparisonRow struct {
	Aspect      string `json:"aspect"`
	Traditional string `json:"traditional"`
	Target      string `json:"target"`
}

type Callout struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func Build() DemoPage {
	return DemoPage{
		Title:    "数据库核心机制可视化实验台",
		Subtitle: "用 Go 提供演示接口，用前端图形化展示 MySQL、PostgreSQL、Redis 与 MongoDB 的关键机制。",
		Sections: []DemoSection{
			{
				ID:       "mysql",
				Name:     "MySQL：锁、索引与存储引擎",
				Overview: "这一部分聚焦并发控制、B+Tree 索引访问路径以及 InnoDB / MyISAM / Memory 的能力差异。",
				Highlights: []string{
					"展示共享锁、排他锁、间隙锁和 Next-Key Lock 的效果",
					"展示 B+Tree 从根节点到叶子节点的命中路径",
					"对比 InnoDB、MyISAM、Memory 的事务、锁粒度与崩溃恢复能力",
				},
				Visualizations: []Visualization{
					{
						Type:        "lock-flow",
						Title:       "事务加锁时序",
						Description: "模拟两个事务访问同一范围时，锁如何阻塞或放行。",
						Steps: []StepItem{
							{Title: "T1 开始事务", Detail: "读取 id BETWEEN 10 AND 20，InnoDB 施加 Next-Key Lock。", State: "active"},
							{Title: "T2 尝试插入", Detail: "插入 id=15 时命中间隙锁，被阻塞等待。", State: "waiting"},
							{Title: "T1 提交", Detail: "提交后释放记录锁和间隙锁。", State: "done"},
							{Title: "T2 获得执行", Detail: "等待结束，插入操作完成。", State: "done"},
						},
						Meta: map[string]string{
							"focus": "锁冲突与可重复读下的幻读防护",
						},
					},
					{
						Type:        "tree",
						Title:       "B+Tree 索引命中路径",
						Description: "以联合索引 (dept_id, age, name) 为例，说明最左前缀与范围扫描。",
						Items: []VisualItem{
							{Title: "Root", Detail: "按 dept_id 分段", Accent: "root", Children: []string{"Branch A", "Branch B", "Branch C"}},
							{Title: "Branch A", Label: "dept_id = 10", Detail: "继续按 age 排序", Accent: "hit", Children: []string{"Leaf 10-21", "Leaf 10-34"}},
							{Title: "Branch B", Label: "dept_id = 20", Detail: "查询未命中分支", Accent: "idle", Children: []string{"Leaf 20-25", "Leaf 20-41"}},
							{Title: "Leaf 10-34", Label: "age >= 30", Detail: "定位到叶子页，再顺序扫描 name。", Accent: "scan"},
						},
						Meta: map[string]string{
							"query": "WHERE dept_id = 10 AND age >= 30 ORDER BY name",
						},
					},
					{
						Type:        "table",
						Title:       "存储引擎能力对比",
						Description: "把常见引擎放在一个表中，便于课堂展示或面试讲解。",
						Columns:     []string{"引擎", "事务", "锁粒度", "索引结构", "崩溃恢复"},
						Rows: [][]string{
							{"InnoDB", "支持", "行锁 + 间隙锁", "聚簇索引 + 二级索引", "支持"},
							{"MyISAM", "不支持", "表锁", "非聚簇索引", "较弱"},
							{"Memory", "不支持", "表锁", "Hash / BTREE", "重启丢失"},
						},
					},
				},
				Callouts: []Callout{
					{Title: "独立判断", Content: "锁并不是越细越好。行锁提升并发，但会增加死锁检测、事务管理与索引维护开销。"},
				},
			},
			{
				ID:       "postgresql",
				Name:     "PostgreSQL：区别于传统数据库的能力",
				Overview: "这里强调 PostgreSQL 在 MVCC、扩展性、复杂数据类型和分析能力上的长期优势。",
				Highlights: []string{
					"MVCC 通过 tuple 版本减少读写阻塞",
					"支持 JSONB、数组、自定义类型与扩展",
					"强大的 CTE、窗口函数、全文检索与逻辑复制能力",
				},
				Visualizations: []Visualization{
					{
						Type:        "comparison",
						Title:       "传统关系库 vs PostgreSQL",
						Description: "把核心差异抽成几个高频讨论维度。",
					},
					{
						Type:        "stepper",
						Title:       "MVCC 版本链",
						Description: "读取旧版本不阻塞写入，写入通过新 tuple 版本推进事务可见性。",
						Steps: []StepItem{
							{Title: "Tuple v1", Detail: "xmin=100, xmax=0，事务快照可见。", State: "done"},
							{Title: "事务 200 更新", Detail: "生成 v2，并把 v1 标记 xmax=200。", State: "active"},
							{Title: "旧快照读取", Detail: "仍可见 v1，不必等待写事务完成。", State: "done"},
							{Title: "VACUUM 清理", Detail: "当旧版本不再可见时回收空间。", State: "idle"},
						},
					},
					{
						Type:        "cards",
						Title:       "生态与扩展能力",
						Description: "这些能力让 PostgreSQL 不只是传统 OLTP 数据库。",
						Items: []VisualItem{
							{Title: "PostGIS", Detail: "原生支持地理空间计算", Accent: "hit"},
							{Title: "JSONB", Detail: "支持半结构化数据与 GIN 索引", Accent: "root"},
							{Title: "Logical Replication", Detail: "更灵活的数据分发与解耦", Accent: "scan"},
							{Title: "Custom Type / Operator", Detail: "支持深度扩展数据模型", Accent: "idle"},
						},
					},
				},
				Comparisons: []ComparisonRow{
					{Aspect: "并发控制", Traditional: "更多依赖锁协调读写", Target: "MVCC 让读写冲突显著下降"},
					{Aspect: "数据模型", Traditional: "以固定表结构为主", Target: "关系模型上叠加 JSONB、数组、范围类型"},
					{Aspect: "可扩展性", Traditional: "扩展点有限", Target: "可扩展类型、函数、索引、插件生态丰富"},
					{Aspect: "分析能力", Traditional: "复杂分析常依赖外部系统", Target: "窗口函数、CTE、全文检索能力较强"},
				},
				Callouts: []Callout{
					{Title: "平衡结论", Content: "PostgreSQL 的优势主要体现在一致性、扩展性和复杂查询能力，但其调优复杂度也通常高于轻量级部署方案。"},
				},
			},
			{
				ID:       "redis",
				Name:     "Redis：原子性与 epoll 事件驱动",
				Overview: "用命令执行序列和 IO 多路复用图解释 Redis 为什么能在单线程执行模型下保持高吞吐。",
				Highlights: []string{
					"单线程命令执行保证单条命令原子性",
					"MULTI/EXEC 与 Lua 脚本扩展复合原子操作",
					"epoll 负责高效监听大量 socket 事件",
				},
				Visualizations: []Visualization{
					{
						Type:        "stepper",
						Title:       "原子命令执行序列",
						Description: "展示 INCR、Lua、事务队列在事件循环中的排队执行。",
						Steps: []StepItem{
							{Title: "Client A -> INCR stock", Detail: "命令进入队列，主线程独占执行。", State: "done"},
							{Title: "Client B -> DECR stock", Detail: "必须等待上一条命令完成后执行。", State: "waiting"},
							{Title: "Lua Script", Detail: "脚本在执行期间不会被其他命令打断。", State: "active"},
							{Title: "回写结果", Detail: "执行完成后统一返回响应。", State: "done"},
						},
					},
					{
						Type:        "epoll",
						Title:       "epoll 驱动的事件循环",
						Description: "把连接注册、就绪通知、命令执行和响应发送串成一条路径。",
						Items: []VisualItem{
							{Title: "监听 socket", Detail: "连接进入 epoll 红黑树", Accent: "root"},
							{Title: "就绪列表", Detail: "活跃 fd 被送入 ready list", Accent: "hit"},
							{Title: "事件循环", Detail: "取出事件、解析 RESP、执行命令", Accent: "scan"},
							{Title: "回写客户端", Detail: "写事件就绪后返回结果", Accent: "idle"},
						},
						Meta: map[string]string{
							"core": "IO 等待交给内核，多数 CPU 时间留给命令执行",
						},
					},
				},
				Callouts: []Callout{
					{Title: "准确表述", Content: "Redis 并非“完全单线程”。经典说法主要指命令执行线程单线程，而网络读写、持久化等在新版本里可由额外线程辅助。"},
				},
			},
			{
				ID:       "mongodb",
				Name:     "MongoDB：文档模型与分布式特性",
				Overview: "强调文档型数据模型、聚合管道、复制集和分片机制，说明它适合快速演进的数据场景。",
				Highlights: []string{
					"灵活 schema 适合迭代快的数据建模",
					"聚合管道让复杂数据变换更直观",
					"复制集与分片帮助提升可用性和水平扩展性",
				},
				Visualizations: []Visualization{
					{
						Type:        "cards",
						Title:       "文档模型特性",
						Description: "展示一条订单文档内联嵌套的典型结构。",
						Items: []VisualItem{
							{Title: "Order", Detail: "顶层文档可直接包含用户、商品、地址等字段", Accent: "root"},
							{Title: "items[]", Detail: "数组内嵌多个商品明细", Accent: "hit"},
							{Title: "shipping.address", Detail: "对象嵌套减少跨表 JOIN", Accent: "scan"},
							{Title: "tags", Detail: "字段可按业务逐步演化", Accent: "idle"},
						},
					},
					{
						Type:        "stepper",
						Title:       "聚合管道",
						Description: "从文档过滤到分组汇总的典型处理链。",
						Steps: []StepItem{
							{Title: "$match", Detail: "过滤 30 天内的订单", State: "done"},
							{Title: "$unwind", Detail: "展开 items 数组", State: "done"},
							{Title: "$group", Detail: "按品类汇总销售额", State: "active"},
							{Title: "$sort", Detail: "按总额降序输出结果", State: "idle"},
						},
					},
					{
						Type:        "cluster",
						Title:       "复制集与分片",
						Description: "把可用性和扩容能力放到同一张图里展示。",
						Items: []VisualItem{
							{Title: "Primary", Detail: "接收写请求并复制 oplog", Accent: "root"},
							{Title: "Secondary A", Detail: "异步复制，支持故障切换", Accent: "hit"},
							{Title: "Secondary B", Detail: "可用于读扩展", Accent: "scan"},
							{Title: "Shard Router", Detail: "按 shard key 路由到不同分片", Accent: "idle"},
						},
					},
				},
				Callouts: []Callout{
					{Title: "客观边界", Content: "MongoDB 的灵活 schema 有利于快速开发，但如果缺少约束治理，也更容易出现字段不一致和查询复杂化问题。"},
				},
			},
		},
	}
}
