package demo

type DemoPage struct {
	Title     string        `json:"title"`
	Subtitle  string        `json:"subtitle"`
	HeroTitle string        `json:"heroTitle"`
	HeroBody  string        `json:"heroBody"`
	HeroStats []HeroStat    `json:"heroStats"`
	Sections  []DemoSection `json:"sections"`
}

type HeroStat struct {
	Value  string `json:"value"`
	Label  string `json:"label"`
	Detail string `json:"detail,omitempty"`
}

type DemoSection struct {
	ID             string          `json:"id"`
	Name           string          `json:"name"`
	Overview       string          `json:"overview"`
	Scenario       string          `json:"scenario,omitempty"`
	Goal           string          `json:"goal,omitempty"`
	Outcomes       []string        `json:"outcomes,omitempty"`
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
	Language    string            `json:"language,omitempty"`
	Content     string            `json:"content,omitempty"`
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
		Title:     "数据库业务场景可视化实验台",
		Subtitle:  "用一套页面同时讲清楚业务链路、查询设计和存储引擎行为，覆盖 MySQL、PostgreSQL、Redis、MongoDB 与 Milvus。",
		HeroTitle: "把数据库原理翻译成业务团队看得懂的流程图、索引路径和查询样例",
		HeroBody:  "页面保留实时分析入口，同时将静态教学内容升级为业务场景导向：MySQL 的锁和索引、PostgreSQL 的递归检索与全文搜索、Redis 的 IO 多路复用和 Lua 原子脚本、MongoDB 的文档检索，以及 Milvus 的向量存储和相似性搜索。",
		HeroStats: []HeroStat{
			{Value: "5", Label: "数据库主题", Detail: "新增 Milvus 向量检索模块"},
			{Value: "11", Label: "业务场景", Detail: "覆盖电商、知识库、客服与推荐"},
			{Value: "20+", Label: "可视化面板", Detail: "流程图、树图、表格和代码示例"},
			{Value: "4", Label: "实时连接器", Detail: "MySQL、PostgreSQL、Redis、MongoDB"},
		},
		Sections: []DemoSection{
			{
				ID:       "mysql",
				Name:     "MySQL：订单与库存场景中的锁和索引",
				Overview: "这一部分围绕秒杀库存预占和运营订单列表查询，解释 InnoDB 如何使用锁、索引访问路径和执行计划来保证正确性与性能。",
				Scenario: "场景：电商系统在流量高峰时需要同时完成支付确认、库存预占、仓库补货，以及运营后台的热点订单分页查询。",
				Goal:     "目标：让产品、后端和 DBA 团队能在一个故事里看懂锁等待、范围保护和联合索引的价值。",
				Outcomes: []string{
					"库存预占冲突可以直观看成行锁、间隙锁和 Next-Key Lock 的竞争过程。",
					"运营后台查单可以映射到联合索引最左匹配和回表成本。",
					"执行计划和索引设计可以用业务语言解释，而不只是数据库术语。",
				},
				Highlights: []string{
					"展示两个支付确认事务为什么会在同一段库存范围上等待。",
					"可视化 `shop_id + status + created_at` 查询的 B+Tree 命中路径。",
					"把读取速度、写入成本和锁范围放到同一个业务决策视图里。",
				},
				Visualizations: []Visualization{
					{
						Type:        "lock-flow",
						Title:       "库存预占的锁等待时间线",
						Description: "支付确认流和补货流竞争同一段库存范围。这个时间线展示 InnoDB 如何把写冲突转成可控等待，而不是超卖。",
						Steps: []StepItem{
							{Title: "T1 锁定库存范围", Detail: "支付流程对仓库 17 的 sku 9001 执行 `SELECT ... FOR UPDATE`，锁住匹配记录及相关间隙。", State: "active"},
							{Title: "T2 插入新补货批次", Detail: "补货流程想往同一范围插入数据，因为命中间隙保护而被阻塞。", State: "waiting"},
							{Title: "T1 提交预占结果", Detail: "预占数量更新完成后，记录锁和间隙锁被释放。", State: "done"},
							{Title: "T2 恢复执行", Detail: "补货记录成功落库，后续事务可以重新评估可售库存。", State: "done"},
						},
						Meta: map[string]string{
							"业务动作": "支付确认与库存预占",
							"核心机制": "记录锁、间隙锁、Next-Key Lock",
						},
					},
					{
						Type:        "tree",
						Title:       "订单看板查询的联合索引命中路径",
						Description: "运营查询按门店、订单状态和时间窗口过滤。树图展示索引如何在回表前先缩小扫描范围。",
						Items: []VisualItem{
							{Title: "根页", Label: "idx_shop_status_created_at", Detail: "索引的第一层先按门店拆分订单。", Accent: "root", Children: []string{"旗舰店", "直营网店", "跨境店铺"}},
							{Title: "分支：旗舰店", Label: "shop_id = 1024", Detail: "进入单个门店后，再按订单状态和时间顺序组织叶子页。", Accent: "hit", Children: []string{"paid", "packed", "shipped"}},
							{Title: "叶子页：paid", Label: "09:00 到 12:00", Detail: "存储引擎落到正确叶子页后，按时间范围顺序扫描。", Accent: "scan"},
							{Title: "回表阶段", Label: "select *", Detail: "如果要读金额、收件信息或地址字段，就要回到聚簇索引。", Accent: "idle"},
						},
						Meta: map[string]string{
							"查询模板": "WHERE shop_id = ? AND status = 'paid' AND created_at >= ? ORDER BY created_at DESC",
						},
					},
					{
						Type:        "table",
						Title:       "业务查询与索引策略矩阵",
						Description: "用一个简明矩阵把高频业务查询、索引方案和需要关注的取舍放在一起。",
						Columns:     []string{"业务查询", "推荐索引", "收益", "代价或风险"},
						Rows: [][]string{
							{"按门店分页查已支付订单", "(shop_id, status, created_at)", "缩小排序和扫描范围", "增加写放大"},
							{"查看单个用户历史订单", "(user_id, created_at)", "快速命中用户时间线", "与运营看板不是同一路径"},
							{"按 sku 和仓库更新库存", "(sku_id, warehouse_id)", "缩小锁定范围", "缺少索引时锁影响会扩大"},
							{"按支付单号查单", "UNIQUE(pay_no)", "快速精确定位", "需要强唯一约束"},
						},
					},
					{
						Type:        "code",
						Title:       "事务与 EXPLAIN 示例",
						Description: "把业务写操作和诊断 SQL 放在一起，页面就能同时解释正确性和性能。",
						Language:    "sql",
						Content:     "BEGIN;\nSELECT batch_id, available_qty\nFROM inventory_batches\nWHERE sku_id = 9001 AND warehouse_id = 17\nORDER BY expire_at\nFOR UPDATE;\n\nUPDATE inventory_batches\nSET reserved_qty = reserved_qty + 2\nWHERE batch_id = 88102;\nCOMMIT;\n\nEXPLAIN FORMAT=JSON\nSELECT id, created_at, total_amount\nFROM orders\nWHERE shop_id = 1024 AND status = 'paid'\n  AND created_at >= '2026-03-17 09:00:00'\nORDER BY created_at DESC\nLIMIT 20;",
					},
				},
				Callouts: []Callout{
					{Title: "业务解读", Content: "在 MySQL 里，锁不是抽象概念。它直接决定会不会超卖，以及第二个事务要等多久。索引也不是越多越好，每新增一个索引都会增加写入维护成本。"},
				},
			},
			{
				ID:       "postgresql",
				Name:     "PostgreSQL：递归查询、全文检索与 JSONB 混合过滤",
				Overview: "这一部分用知识库和组织树场景说明 PostgreSQL 如何在一条查询中同时完成递归遍历、全文排序和 JSONB 条件过滤。",
				Scenario: "场景：企业知识库既要按关键词搜索文档，又要按组织树限制可见范围，还要根据 JSONB 元数据中的地区、等级和标签进行过滤。",
				Goal:     "目标：把 PostgreSQL 展示成复杂检索引擎，而不只是传统关系数据库。",
				Outcomes: []string{
					"递归 CTE 负责展开组织树或分类树中的可见部分。",
					"全文检索负责把自然语言关键字转成带排序的文档匹配结果。",
					"JSONB 让弹性元数据可以继续留在数据库内被高效过滤。",
				},
				Highlights: []string{
					"用 WITH RECURSIVE 解释组织层级和类目展开。",
					"用 tsvector 和 tsquery 解释知识搜索的召回与排序。",
					"用 JSONB 和 GIN 解释固定结构与弹性属性如何共存。",
				},
				Visualizations: []Visualization{
					{
						Type:        "stepper",
						Title:       "递归查询展开流程",
						Description: "先找到起始节点，再递归展开可见树，最后将结果交给全文排序和 JSONB 过滤。",
						Steps: []StepItem{
							{Title: "锚点节点", Detail: "先定位当前用户对应的组织单元或文档分类根节点。", State: "done"},
							{Title: "递归展开", Detail: "通过 WITH RECURSIVE 持续沿着 `parent_id` 收集所有可达子节点。", State: "active"},
							{Title: "全文排序", Detail: "文档标题和正文向量与 `websearch_to_tsquery` 匹配并计算排序分值。", State: "done"},
							{Title: "JSONB 过滤", Detail: "对 `doc_meta` 中的地区、等级和标签做进一步筛选。", State: "idle"},
						},
						Meta: map[string]string{
							"业务对象": "知识树与组织可见性",
							"输出结果": "带排序的可见文档",
						},
					},
					{
						Type:        "table",
						Title:       "混合查询中的索引分工",
						Description: "一条业务查询通常要同时依赖多种索引，这张表把每个阶段的责任拆开。",
						Columns:     []string{"能力", "典型 SQL", "索引或机制", "业务价值"},
						Rows: [][]string{
							{"递归树遍历", "WITH RECURSIVE org_tree AS (...)", "主键与 parent_id 索引", "快速找出所有可见部门"},
							{"全文匹配", "search_vector @@ websearch_to_tsquery(...)", "GIN(search_vector)", "按语义召回文档而非仅精确词"},
							{"JSONB 元数据过滤", "doc_meta @> '{\"region\":\"cn\"}'", "GIN(doc_meta jsonb_path_ops)", "按弹性属性过滤"},
							{"排序与分页", "ORDER BY rank DESC, updated_at DESC", "Top-N 策略与排序支持", "保证搜索结果稳定可读"},
						},
					},
					{
						Type:        "code",
						Title:       "混合检索 SQL 示例",
						Description: "这段示例把递归可见性、全文排序和 JSONB 过滤放进同一条业务查询里。",
						Language:    "sql",
						Content:     "WITH RECURSIVE org_tree AS (\n    SELECT id\n    FROM org_units\n    WHERE id = 42\n  UNION ALL\n    SELECT child.id\n    FROM org_units child\n    JOIN org_tree parent ON child.parent_id = parent.id\n), ranked_docs AS (\n    SELECT d.id,\n           d.title,\n           ts_rank_cd(d.search_vector, websearch_to_tsquery('simple', 'refund approval flow')) AS rank\n    FROM knowledge_docs d\n    JOIN org_tree t ON d.org_unit_id = t.id\n    WHERE d.search_vector @@ websearch_to_tsquery('simple', 'refund approval flow')\n      AND d.doc_meta @> '{\"region\":\"cn\",\"level\":\"internal\"}'::jsonb\n      AND d.doc_meta -> 'tags' ? 'finance'\n)\nSELECT id, title, rank\nFROM ranked_docs\nORDER BY rank DESC, id DESC\nLIMIT 10;",
					},
					{
						Type:        "cards",
						Title:       "最适合落地的业务模块",
						Description: "这些模块里，递归 SQL、全文检索和 JSONB 经常能替掉多套拆散的能力。",
						Items: []VisualItem{
							{Title: "知识库搜索", Detail: "搜索和组织可见性在一个系统里完成。", Accent: "hit"},
							{Title: "目录中心", Detail: "递归类目加灵活属性字段。", Accent: "root"},
							{Title: "流程模板中心", Detail: "可同时按团队、版本、标签和关键字过滤。", Accent: "scan"},
							{Title: "审计控制台", Detail: "保留关系表，同时承载不断演进的审核元数据。", Accent: "idle"},
						},
					},
				},
				Callouts: []Callout{
					{Title: "讲解角度", Content: "对大多数业务团队来说，PostgreSQL 的价值不只是 MVCC，而是那些原本需要拆成搜索引擎、权限树服务和元数据过滤器的问题，往往可以先用一条 SQL 交付出来。"},
				},
			},
			{
				ID:       "redis",
				Name:     "Redis：IO 多路复用与多命令 Lua 脚本",
				Overview: "这一部分使用秒杀网关故事解释 Redis 为什么能在大量连接下保持高性能，以及为什么 Lua 比客户端拼接多条命令更适合原子业务动作。",
				Scenario: "场景：秒杀活动期间，请求必须经过限流、重复用户校验、库存扣减和事件投递，同时不能超卖，也不能让计数器状态不一致。",
				Goal:     "目标：展示 IO 多路复用、事件循环和服务端脚本如何在一个真实业务流中协同工作。",
				Outcomes: []string{
					"epoll 或 kqueue 让网络等待从业务执行路径中移开。",
					"事件循环把命令执行串行化，避免共享内存锁竞争。",
					"Lua 把多次校验和写入打包成一个原子操作。",
				},
				Highlights: []string{
					"把 socket 就绪、命令解析、命令执行和响应写回连成一条可见链路。",
					"展示一个同时完成资格校验、扣库存和推送事件的业务脚本。",
					"从业务一致性角度对比客户端拼接命令、MULTI EXEC 和 Lua。",
				},
				Visualizations: []Visualization{
					{
						Type:        "epoll",
						Title:       "Redis 事件循环中的 IO 多路复用",
						Description: "Redis 先判断哪些文件描述符准备好了，再把 CPU 时间花在命令执行上，从而避免大量空闲连接消耗主执行路径。",
						Items: []VisualItem{
							{Title: "连接注册", Label: "accept + fd", Detail: "新客户端 socket 注册到事件循环中，等待可读或可写事件。", Accent: "root"},
							{Title: "就绪列表", Label: "epoll_wait", Detail: "内核只返回活跃 socket，避免遍历全部连接。", Accent: "hit"},
							{Title: "命令执行", Label: "parse -> execute", Detail: "主线程按顺序执行 GET、INCR、EVALSHA 等命令。", Accent: "scan"},
							{Title: "响应回写", Label: "send reply", Detail: "连接可写时把结果发回客户端，然后进入下一轮事件循环。", Accent: "idle"},
						},
						Meta: map[string]string{
							"业务收益":  "大量活跃连接下保持稳定时延",
							"机制关键词": "事件循环、非阻塞 IO、就绪列表",
						},
					},
					{
						Type:        "stepper",
						Title:       "用 Lua 实现秒杀原子动作",
						Description: "脚本把多个业务规则合并成一个不可被其他命令打断的执行单元。",
						Steps: []StepItem{
							{Title: "检查用户资格", Detail: "先拒绝重复购买或黑名单用户，避免浪费库存。", State: "done"},
							{Title: "校验库存", Detail: "读取秒杀库存，当计数器已经归零时立即停止。", State: "active"},
							{Title: "扣减并记录", Detail: "在同一脚本里完成库存扣减并记录成功用户。", State: "done"},
							{Title: "推送订单事件", Detail: "追加一条下游事件供异步订单服务消费。", State: "idle"},
						},
						Meta: map[string]string{
							"适用场景": "秒杀、领券、令牌桶、幂等控制",
						},
					},
					{
						Type:        "code",
						Title:       "资格校验、扣库存与事件投递 Lua 示例",
						Description: "和客户端多次往返相比，这段脚本可以保证业务状态迁移是原子的。",
						Language:    "lua",
						Content:     "local stockKey = KEYS[1]\nlocal successKey = KEYS[2]\nlocal streamKey = KEYS[3]\nlocal userId = ARGV[1]\nlocal orderId = ARGV[2]\n\nif redis.call('SISMEMBER', successKey, userId) == 1 then\n  return {err = 'DUPLICATE_USER'}\nend\n\nlocal stock = tonumber(redis.call('GET', stockKey) or '0')\nif stock <= 0 then\n  return {err = 'OUT_OF_STOCK'}\nend\n\nredis.call('DECR', stockKey)\nredis.call('SADD', successKey, userId)\nredis.call('XADD', streamKey, '*', 'user_id', userId, 'order_id', orderId)\nreturn {'OK', tostring(stock - 1)}",
					},
					{
						Type:        "table",
						Title:       "客户端命令、MULTI EXEC 与 Lua 的对比",
						Description: "同样是四步业务动作，逻辑放在哪里，表现完全不同。",
						Columns:     []string{"方案", "执行方式", "优势", "风险"},
						Rows: [][]string{
							{"客户端拼接命令", "GET 再 DECR 再 SADD 再 XADD", "上手简单", "中途失败容易留下不一致状态"},
							{"MULTI EXEC", "事务队列执行", "可以保证命令顺序", "条件分支能力较弱"},
							{"Lua 脚本", "服务端单次执行", "复杂规则下原子性更强", "脚本过长会阻塞主线程"},
						},
					},
				},
				Callouts: []Callout{
					{Title: "系统边界", Content: "Redis 最适合做高并发入口的快速决策层。它应该尽快做出判断、保护热点资源，然后把长期状态交给主事务存储系统。"},
				},
			},
			{
				ID:       "mongodb",
				Name:     "MongoDB：客服与内容工作台的文档检索",
				Overview: "这一部分说明文档模型如何承载嵌套业务对象、数组过滤和读优化检索，适合客服和运营工作台类场景。",
				Scenario: "场景：客服平台希望一个用户画像文档中直接包含订单摘要、设备信息、最近工单、标签和渠道元数据，方便坐席一站式检索和查看。",
				Goal:     "目标：把 MongoDB 展示为实用的文档检索层，而不是简单的 JSON 存储。",
				Outcomes: []string{
					"一个用户画像文档可以容纳多个嵌套业务对象，减少大量 JOIN 读取。",
					"数组字段和嵌入对象可以直接支持标签、渠道和工单状态过滤。",
					"聚合与投影可以生成工作台真正需要的返回结构。",
				},
				Highlights: []string{
					"用客户画像文档解释嵌套对象、数组和可演进字段。",
					"展示 match、elemMatch、project、sort、limit 如何组成一次检索流水线。",
					"把叙事重心放在读模型和以文档为中心的访问模式上。",
				},
				Visualizations: []Visualization{
					{
						Type:        "cards",
						Title:       "客户画像文档结构",
						Description: "一条文档就可以携带客服通话时最关心的大部分上下文信息。",
						Items: []VisualItem{
							{Title: "profile", Detail: "身份信息、会员等级和获客渠道。", Accent: "root"},
							{Title: "devices[]", Detail: "最近设备、风险标记和地理上下文。", Accent: "hit"},
							{Title: "tickets[]", Detail: "退款或服务工单的状态与升级信息。", Accent: "scan"},
							{Title: "tags 与 preferences", Detail: "营销和产品信号会随业务演进而变化。", Accent: "idle"},
						},
					},
					{
						Type:        "stepper",
						Title:       "文档检索流水线",
						Description: "客服工作台通常先做针对性过滤，再返回一个轻量且可直接渲染的结构。",
						Steps: []StepItem{
							{Title: "主条件匹配", Detail: "先按地区、会员等级、标签和设备系统过滤。", State: "done"},
							{Title: "工单过滤", Detail: "只保留包含待升级工单的文档。", State: "active"},
							{Title: "字段投影", Detail: "只返回工作台真正需要渲染的字段。", State: "done"},
							{Title: "排序与限制", Detail: "优先展示最近活跃的客户。", State: "idle"},
						},
						Meta: map[string]string{
							"业务对象": "客服用户画像",
						},
					},
					{
						Type:        "code",
						Title:       "客服工作台查询示例",
						Description: "嵌套数组、标签和内嵌字段都可以一起参与一次检索。",
						Language:    "javascript",
						Content:     "db.customer_profiles.find(\n  {\n    region: 'east-cn',\n    vip_level: { $gte: 3 },\n    tags: 'high-value',\n    devices: { $elemMatch: { os: 'iOS', risk_flag: false } },\n    tickets: { $elemMatch: { status: 'pending-escalation', category: 'refund' } }\n  },\n  {\n    user_id: 1,\n    nickname: 1,\n    vip_level: 1,\n    tags: 1,\n    'tickets.$': 1,\n    last_active_at: 1\n  }\n).sort({ last_active_at: -1 }).limit(20);",
					},
					{
						Type:        "table",
						Title:       "文档检索的索引建议",
						Description: "即使文档结构灵活，高频页面仍然需要稳定的访问路径。",
						Columns:     []string{"过滤模式", "建议索引", "用途", "注意点"},
						Rows: [][]string{
							{"region + vip_level", "{ region: 1, vip_level: -1 }", "快速定位高价值客户", "排序方向要匹配页面需求"},
							{"tags", "{ tags: 1 }", "支持标签圈选", "数组索引会增加索引体积"},
							{"tickets.status", "{ 'tickets.status': 1 }", "按工单状态检索", "数组字段会形成 multikey 索引"},
							{"last_active_at", "{ last_active_at: -1 }", "快速排序客服队列", "通常适合放入复合索引"},
						},
					},
				},
				Callouts: []Callout{
					{Title: "设计提醒", Content: "MongoDB 对文档型读取非常强大，但并不意味着所有关系都应该永远内嵌。文档过大和 schema 漂移失控，会让更新、索引和运维行为都变得更重。"},
				},
			},
			{
				ID:       "milvus",
				Name:     "Milvus：高维向量存储与相似性搜索",
				Overview: "这一部分用语义搜索和推荐场景解释 embedding 如何被写入、分段、建立索引并在向量数据库内完成搜索。",
				Scenario: "场景：企业把客服知识文章、商品图片或用户行为转成向量，希望对自然语言问题或推荐请求召回最相似的内容。",
				Goal:     "目标：通过写入链路、ANN 索引和搜索流程，让向量数据库的行为变得具体，而不只是停留在算法名词。",
				Outcomes: []string{
					"Embedding 生成只是开始，真正影响时延和召回质量的是 segment 存储和 ANN 索引设计。",
					"向量搜索通常承担第一阶段召回，再交给重排或业务规则细化结果。",
					"Milvus 解决的是语义相似问题，它应与事务数据库协同，而不是替代事务数据库。",
				},
				Highlights: []string{
					"可视化向量如何与标量元数据一起存入 collection。",
					"展示 query vector 到 ANN 召回，再到标量过滤和最终 Top-K 的路径。",
					"用客服、知识搜索和推荐系统的语言解释向量检索。",
				},
				Visualizations: []Visualization{
					{
						Type:        "stepper",
						Title:       "向量写入与索引构建流程",
						Description: "业务内容先被编码成 embedding，再连同元数据写入、持久化为 segment，之后建立 ANN 索引。",
						Steps: []StepItem{
							{Title: "生成 embedding", Detail: "文本、图片或行为特征被编码成 768 维或 1536 维向量。", State: "done"},
							{Title: "写入 collection", Detail: "每个向量会和主键、租户、渠道、更新时间等字段一起存储。", State: "active"},
							{Title: "持久化为 segment", Detail: "Milvus 将数据组织成 segment，方便 flush、compact 和索引。", State: "done"},
							{Title: "建立 ANN 索引", Detail: "向量字段上建立 IVF、HNSW 等近似最近邻索引。", State: "idle"},
						},
						Meta: map[string]string{
							"业务对象": "文档向量、商品向量、用户兴趣向量",
						},
					},
					{
						Type:        "cluster",
						Title:       "相似性搜索路径",
						Description: "一次语义搜索通常遵循：编码查询、召回候选、应用标量过滤，再重排或组装最终答案。",
						Items: []VisualItem{
							{Title: "查询编码器", Label: "文本转向量", Detail: "用户问题先被实时编码成 query vector。", Accent: "root"},
							{Title: "ANN 召回", Label: "Top-N 候选", Detail: "向量索引快速返回最近邻候选，而不是全量精确比对。", Accent: "hit"},
							{Title: "标量过滤", Label: "租户、渠道、时效", Detail: "业务元数据负责剔除语义相似但业务上不合法的结果。", Accent: "scan"},
							{Title: "重排与组装", Label: "最终 Top-K", Detail: "重排模型或业务规则再产生最终列表或答案片段。", Accent: "idle"},
						},
						Meta: map[string]string{
							"关键指标": "召回率、时延、最终答案质量",
						},
					},
					{
						Type:        "table",
						Title:       "常见 Milvus 索引选择",
						Description: "不同 ANN 策略在内存、时延和召回之间有不同取舍。",
						Columns:     []string{"索引类型", "最适合场景", "特点", "业务提示"},
						Rows: [][]string{
							{"IVF_FLAT", "中等规模检索", "训练简单，基线稳定", "适合原型和首版上线"},
							{"HNSW", "高召回、低时延", "在线效果好，但内存成本高", "适合搜索和客服召回"},
							{"IVF_PQ", "压缩型大规模搜索", "节省内存，但有一定精度折中", "适合向量量级快速增长场景"},
							{"DiskANN 类方案", "超大规模低成本检索", "更依赖存储层级", "适合容量优先的场景"},
						},
					},
					{
						Type:        "code",
						Title:       "语义搜索请求结构",
						Description: "这个示例展示向量相似度和业务标量过滤如何在一次请求里配合。",
						Language:    "json",
						Content:     "{\n  \"collectionName\": \"kb_chunks\",\n  \"vector\": \"[0.018, -0.227, ... 768 dims]\",\n  \"annsField\": \"embedding\",\n  \"metricType\": \"COSINE\",\n  \"limit\": 5,\n  \"filter\": \"tenant_id == 'acme' and channel in ['refund', 'logistics']\",\n  \"outputFields\": [\"doc_id\", \"chunk_text\", \"channel\", \"updated_at\"]\n}",
					},
				},
				Callouts: []Callout{
					{Title: "定位建议", Content: "Milvus 最适合被解释成语义召回层。它负责快速找到相似内容，而权限、精确过滤、事务更新和最终业务状态仍应由原本负责这些职责的系统来处理。"},
				},
			},
		},
	}
}
