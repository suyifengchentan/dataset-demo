async function loadDemo() {
  const response = await fetch("/api/demo");
  if (!response.ok) {
    throw new Error("failed to load demo");
  }
  return response.json();
}

async function analyzeLive(payload) {
  const response = await fetch("/api/live/analyze", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  });
  if (!response.ok) {
    throw new Error(await response.text());
  }
  return response.json();
}

async function simulateBusinessWrite(payload) {
  const response = await fetch("/api/live/simulate", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  });
  if (!response.ok) {
    throw new Error(await response.text());
  }
  return response.json();
}

function createElement(tag, className, text) {
  const element = document.createElement(tag);
  if (className) {
    element.className = className;
  }
  if (text !== undefined) {
    element.textContent = text;
  }
  return element;
}

function renderHero(data) {
  document.title = data.title || "数据库业务场景可视化实验台";
  document.getElementById("page-title").textContent = data.title || "";
  document.getElementById("page-subtitle").textContent = data.subtitle || "";
  document.getElementById("hero-heading").textContent = data.heroTitle || "";
  document.getElementById("hero-description").textContent = data.heroBody || "";

  const heroGrid = document.getElementById("hero-grid");
  heroGrid.innerHTML = "";

  const fallbackStats = [
    { value: String(data.sections?.length || 0), label: "数据库主题" },
    {
      value: String(
        (data.sections || []).reduce((count, section) => count + (section.visualizations?.length || 0), 0),
      ),
      label: "可视化面板",
    },
  ];

  const stats = data.heroStats?.length ? data.heroStats : fallbackStats;
  stats.forEach((item) => {
    const card = createElement("div", "mini-card");
    card.appendChild(createElement("span", "chip", item.label));
    card.appendChild(createElement("strong", "", item.value));
    if (item.detail) {
      card.appendChild(createElement("p", "", item.detail));
    }
    heroGrid.appendChild(card);
  });
}

function renderNav(sections) {
  const nav = document.getElementById("nav");
  nav.innerHTML = "";

  const liveButton = createElement("button", "", "实时数据库输出");
  liveButton.addEventListener("click", () => {
    document.getElementById("live-content")?.scrollIntoView({ behavior: "smooth", block: "start" });
  });
  nav.appendChild(liveButton);

  sections.forEach((section, index) => {
    const button = createElement("button", index === 0 ? "active" : "", section.name);
    button.addEventListener("click", () => {
      document.querySelectorAll(".nav button").forEach((node) => node.classList.remove("active"));
      button.classList.add("active");
      document.getElementById(section.id)?.scrollIntoView({ behavior: "smooth", block: "start" });
    });
    nav.appendChild(button);
  });
}

function renderHighlights(items) {
  const list = createElement("ul", "highlight-list");
  items.forEach((item) => {
    list.appendChild(createElement("li", "", item));
  });
  return list;
}

function renderOutcomeList(items) {
  const list = createElement("div", "outcome-list");
  items.forEach((item) => {
    const node = createElement("div", "outcome-item");
    node.appendChild(createElement("span", "chip", "业务价值"));
    node.appendChild(createElement("p", "", item));
    list.appendChild(node);
  });
  return list;
}

function renderScenario(section) {
  if (!section.scenario && !section.goal && !section.outcomes?.length) {
    return null;
  }

  const wrap = createElement("div", "scenario-banner");

  if (section.scenario) {
    const card = createElement("div", "scenario-card");
    card.appendChild(createElement("span", "chip", "业务场景"));
    card.appendChild(createElement("p", "", section.scenario));
    wrap.appendChild(card);
  }

  if (section.goal) {
    const card = createElement("div", "scenario-card");
    card.appendChild(createElement("span", "chip", "讲解目标"));
    card.appendChild(createElement("p", "", section.goal));
    wrap.appendChild(card);
  }

  if (section.outcomes?.length) {
    wrap.appendChild(renderOutcomeList(section.outcomes));
  }

  return wrap;
}

function renderStepper(viz) {
  const wrap = createElement("div", "stepper");
  viz.steps.forEach((step, idx) => {
    const node = createElement("div", `step ${step.state || ""}`.trim());
    node.dataset.index = String(idx + 1);
    node.appendChild(createElement("div", "step-title", step.title));
    node.appendChild(createElement("div", "", step.detail));
    wrap.appendChild(node);
  });
  return wrap;
}

function renderItems(items, className) {
  const wrap = createElement("div", className);
  items.forEach((item) => {
    const node = createElement("div", `node ${item.accent || ""}`.trim());
    node.appendChild(createElement("h5", "", item.title));
    if (item.label) {
      node.appendChild(createElement("span", "node-label", item.label));
    }
    node.appendChild(createElement("div", "", item.detail || ""));
    if (item.children?.length) {
      node.appendChild(createElement("div", "node-children", item.children.join("  |  ")));
    }
    wrap.appendChild(node);
  });
  return wrap;
}

function renderCards(items) {
  const wrap = createElement("div", "cards-grid");
  items.forEach((item) => {
    const card = createElement("div", `fact-card ${item.accent || ""}`.trim());
    card.appendChild(createElement("h5", "", item.title));
    card.appendChild(createElement("p", "", item.detail || ""));
    wrap.appendChild(card);
  });
  return wrap;
}

function renderTable(columns, rows, className) {
  const table = createElement("table", className);
  const thead = document.createElement("thead");
  const headerRow = document.createElement("tr");
  columns.forEach((column) => headerRow.appendChild(createElement("th", "", column)));
  thead.appendChild(headerRow);
  table.appendChild(thead);

  const tbody = document.createElement("tbody");
  rows.forEach((row) => {
    const tr = document.createElement("tr");
    row.forEach((cell) => tr.appendChild(createElement("td", "", cell)));
    tbody.appendChild(tr);
  });
  table.appendChild(tbody);
  return table;
}

function renderCodeVisualization(viz) {
  const wrap = createElement("div", "code-panel");
  const pre = createElement("pre", "code-block");
  const code = createElement("code", viz.language ? `language-${viz.language}` : "", viz.content || "");
  pre.appendChild(code);
  wrap.appendChild(pre);
  return wrap;
}

function renderMeta(meta) {
  if (!meta || !Object.keys(meta).length) {
    return null;
  }

  const row = createElement("div", "meta-row");
  Object.entries(meta).forEach(([key, value]) => {
    row.appendChild(createElement("span", "meta-chip", `${key}: ${value}`));
  });
  return row;
}

function renderVisualization(viz, comparisons) {
  const panel = createElement("article", "panel");
  panel.appendChild(createElement("h4", "", viz.title));
  panel.appendChild(createElement("p", "", viz.description));

  const meta = renderMeta(viz.meta);
  if (meta) {
    panel.appendChild(meta);
  }

  switch (viz.type) {
    case "lock-flow":
    case "stepper":
      panel.appendChild(renderStepper(viz));
      break;
    case "tree":
      panel.appendChild(renderItems(viz.items || [], "tree"));
      break;
    case "epoll":
      panel.appendChild(renderItems(viz.items || [], "epoll-flow"));
      break;
    case "cluster":
      panel.appendChild(renderItems(viz.items || [], "cluster-diagram"));
      break;
    case "cards":
      panel.appendChild(renderCards(viz.items || []));
      break;
    case "table":
      panel.appendChild(renderTable(viz.columns || [], viz.rows || [], "data-table"));
      break;
    case "comparison":
      if (comparisons?.length) {
        panel.appendChild(
          renderTable(
            ["维度", "传统方案", "目标方案"],
            comparisons.map((item) => [item.aspect, item.traditional, item.target]),
            "compare-table",
          ),
        );
      }
      break;
    case "code":
      panel.appendChild(renderCodeVisualization(viz));
      break;
    default:
      panel.appendChild(createElement("p", "", "该类型的可视化暂未实现。"));
  }

  return panel;
}

function renderCallouts(callouts) {
  const wrap = createElement("div", "callout-grid");
  callouts.forEach((item) => {
    const node = createElement("div", "callout");
    node.appendChild(createElement("h5", "", item.title));
    node.appendChild(createElement("p", "", item.content));
    wrap.appendChild(node);
  });
  return wrap;
}

function renderStaticSections(sections) {
  const content = document.getElementById("content");
  content.innerHTML = "";

  sections.forEach((section) => {
    const block = createElement("section", "section");
    block.id = section.id;

    const header = createElement("div", "section-header");
    const copy = document.createElement("div");
    copy.appendChild(createElement("div", "chip", section.name));
    copy.appendChild(createElement("h3", "", section.name));
    copy.appendChild(createElement("p", "", section.overview));
    header.appendChild(copy);
    block.appendChild(header);

    const scenario = renderScenario(section);
    if (scenario) {
      block.appendChild(scenario);
    }

    block.appendChild(renderHighlights(section.highlights || []));

    const vizGrid = createElement("div", "viz-grid");
    (section.visualizations || []).forEach((viz) => {
      vizGrid.appendChild(renderVisualization(viz, section.comparisons));
    });
    block.appendChild(vizGrid);

    if (section.callouts?.length) {
      block.appendChild(renderCallouts(section.callouts));
    }

    content.appendChild(block);
  });
}

function renderMetricCard(metric) {
  const card = createElement("div", `metric-card ${metric.tone || ""}`.trim());
  card.appendChild(createElement("span", "chip", metric.label));
  card.appendChild(createElement("strong", "", metric.value || "-"));
  if (metric.hint) {
    card.appendChild(createElement("p", "", metric.hint));
  }
  return card;
}

function renderLiveTable(table) {
  const panel = createElement("article", "panel");
  panel.appendChild(createElement("h4", "", table.title));
  panel.appendChild(renderTable(table.columns || [], table.rows || [], "data-table"));
  return panel;
}

function renderSnippet(snippet) {
  const panel = createElement("article", "panel");
  panel.appendChild(createElement("h4", "", snippet.title));
  const pre = createElement("pre", "code-block");
  const code = createElement(
    "code",
    snippet.language ? `language-${snippet.language}` : "",
    snippet.content || "",
  );
  pre.appendChild(code);
  panel.appendChild(pre);
  return panel;
}

function renderReportSections(report, options) {
  const target = document.getElementById(options.contentId);
  const meta = document.getElementById(options.metaId);
  target.innerHTML = "";

  if (!report.sections?.length) {
    meta.textContent = options.emptyText;
    return;
  }

  const summary = report.summary ? ` | ${report.summary}` : "";
  meta.textContent = `${options.timePrefix}${report.generatedAt}${summary}`;

  report.sections.forEach((section) => {
    const block = createElement("section", "section");

    const header = createElement("div", "section-header");
    const copy = document.createElement("div");
    copy.appendChild(createElement("div", "chip", section.name));
    copy.appendChild(createElement("h3", "", section.name));
    copy.appendChild(createElement("p", "", section.summary || ""));
    header.appendChild(copy);

    const badge = createElement(
      "span",
      `status-badge ${section.status || "idle"}`.trim(),
      section.connected ? "已连接" : "错误",
    );
    header.appendChild(badge);
    block.appendChild(header);

    if (section.error) {
      const errorNode = createElement("div", "callout");
      errorNode.appendChild(createElement("h5", "", "连接错误"));
      errorNode.appendChild(createElement("p", "", section.error));
      block.appendChild(errorNode);
    }

    if (section.metrics?.length) {
      const metricGrid = createElement("div", "metric-grid");
      section.metrics.forEach((metric) => metricGrid.appendChild(renderMetricCard(metric)));
      block.appendChild(metricGrid);
    }

    if (section.warnings?.length) {
      const warningGrid = createElement("div", "callout-grid");
      section.warnings.forEach((warning) => {
        const node = createElement("div", "callout");
        node.appendChild(createElement("h5", "", "告警"));
        node.appendChild(createElement("p", "", warning));
        warningGrid.appendChild(node);
      });
      block.appendChild(warningGrid);
    }

    if (section.tables?.length) {
      const tableGrid = createElement("div", "viz-grid");
      section.tables.forEach((table) => tableGrid.appendChild(renderLiveTable(table)));
      block.appendChild(tableGrid);
    }

    if (section.snippets?.length) {
      const snippetGrid = createElement("div", "viz-grid");
      section.snippets.forEach((snippet) => snippetGrid.appendChild(renderSnippet(snippet)));
      block.appendChild(snippetGrid);
    }

    target.appendChild(block);
  });
}

function fieldValue(id) {
  return document.getElementById(id)?.value?.trim() || "";
}

function checked(id) {
  return Boolean(document.getElementById(id)?.checked);
}

const exampleConfig = {
  mysql: {
    enabled: true,
    dsn: "root:root@tcp(127.0.0.1:3306)/visual_lab",
    schema: "visual_lab",
    table: "employees",
    explainQuery: "SELECT * FROM employees WHERE dept_id = 10 AND age >= 30 ORDER BY name;",
  },
  postgres: {
    enabled: true,
    dsn: "postgres://postgres:postgres@127.0.0.1:5432/visual_lab?sslmode=disable",
    schema: "public",
    table: "orders",
  },
  redis: {
    enabled: true,
    addr: "127.0.0.1:6379",
    username: "",
    password: "",
    db: 0,
    keyPrefix: "dbdemo:atomic",
  },
  mongodb: {
    enabled: true,
    uri: "mongodb://127.0.0.1:27017",
    database: "visual_lab",
    collection: "orders",
  },
};

function saveConfig(payload) {
  localStorage.setItem("db-visual-lab-config", JSON.stringify(payload));
}

function loadStoredConfig() {
  const raw = localStorage.getItem("db-visual-lab-config");
  if (!raw) {
    return null;
  }
  try {
    return JSON.parse(raw);
  } catch {
    return null;
  }
}

function populateForm(config) {
  const mysql = config.mysql || {};
  document.getElementById("mysql-enabled").checked = Boolean(mysql.enabled);
  document.getElementById("mysql-dsn").value = mysql.dsn || "";
  document.getElementById("mysql-schema").value = mysql.schema || "";
  document.getElementById("mysql-table").value = mysql.table || "";
  document.getElementById("mysql-explain").value = mysql.explainQuery || "";

  const postgres = config.postgres || {};
  document.getElementById("postgres-enabled").checked = Boolean(postgres.enabled);
  document.getElementById("postgres-dsn").value = postgres.dsn || "";
  document.getElementById("postgres-schema").value = postgres.schema || "";
  document.getElementById("postgres-table").value = postgres.table || "";

  const redis = config.redis || {};
  document.getElementById("redis-enabled").checked = Boolean(redis.enabled);
  document.getElementById("redis-addr").value = redis.addr || "";
  document.getElementById("redis-username").value = redis.username || "";
  document.getElementById("redis-password").value = redis.password || "";
  document.getElementById("redis-db").value = redis.db ?? 0;
  document.getElementById("redis-prefix").value = redis.keyPrefix || "";

  const mongo = config.mongodb || {};
  document.getElementById("mongo-enabled").checked = Boolean(mongo.enabled);
  document.getElementById("mongo-uri").value = mongo.uri || "";
  document.getElementById("mongo-database").value = mongo.database || "";
  document.getElementById("mongo-collection").value = mongo.collection || "";
}

function collectPayload() {
  return {
    timeoutSec: 8,
    mysql: {
      enabled: checked("mysql-enabled"),
      dsn: fieldValue("mysql-dsn"),
      schema: fieldValue("mysql-schema"),
      table: fieldValue("mysql-table"),
      explainQuery: fieldValue("mysql-explain"),
    },
    postgres: {
      enabled: checked("postgres-enabled"),
      dsn: fieldValue("postgres-dsn"),
      schema: fieldValue("postgres-schema"),
      table: fieldValue("postgres-table"),
    },
    redis: {
      enabled: checked("redis-enabled"),
      addr: fieldValue("redis-addr"),
      username: fieldValue("redis-username"),
      password: fieldValue("redis-password"),
      db: Number(document.getElementById("redis-db")?.value || 0),
      keyPrefix: fieldValue("redis-prefix"),
    },
    mongodb: {
      enabled: checked("mongo-enabled"),
      uri: fieldValue("mongo-uri"),
      database: fieldValue("mongo-database"),
      collection: fieldValue("mongo-collection"),
    },
  };
}

function setBanner(message, state) {
  const banner = document.getElementById("live-banner");
  banner.className = `live-banner ${state || ""}`.trim();
  banner.innerHTML = "";
  const label =
    state === "error" ? "失败" : state === "loading" ? "连接中" : state === "success" ? "完成" : "提示";
  banner.appendChild(createElement("strong", "", label));
  banner.appendChild(createElement("span", "", message));
}

function bindActions() {
  document.getElementById("fill-example").addEventListener("click", () => {
    populateForm(exampleConfig);
    setBanner("已填充 Docker Compose 示例连接参数。", "info");
  });

  document.getElementById("clear-config").addEventListener("click", () => {
    populateForm({
      mysql: { enabled: false },
      postgres: { enabled: false },
      redis: { enabled: false, db: 0 },
      mongodb: { enabled: false },
    });
    localStorage.removeItem("db-visual-lab-config");
    document.getElementById("simulation-content").innerHTML = "";
    document.getElementById("simulation-meta").textContent = "等待触发模拟。";
    document.getElementById("live-content").innerHTML = "";
    document.getElementById("live-meta").textContent = "等待连接。";
    setBanner("已清空本地保存的连接配置。", "info");
  });

  document.getElementById("simulate-write").addEventListener("click", async () => {
    const payload = collectPayload();
    saveConfig(payload);
    setBanner("正在执行业务写入模拟：会尝试向已启用数据库写入示例数据。", "loading");

    try {
      const report = await simulateBusinessWrite(payload);
      renderReportSections(report, {
        contentId: "simulation-content",
        metaId: "simulation-meta",
        emptyText: "当前没有启用任何可模拟写入的数据库连接。",
        timePrefix: "最近一次模拟时间：",
      });
      setBanner("业务写入模拟完成。你可以重复点击观察 Redis 库存、事件流和各库写入结果变化。", "success");
    } catch (error) {
      setBanner(`业务写入模拟失败：${error.message}`, "error");
    }
  });

  document.getElementById("analyze").addEventListener("click", async () => {
    const payload = collectPayload();
    saveConfig(payload);
    setBanner("正在连接已启用的数据库并采集实时指标。", "loading");

    try {
      const report = await analyzeLive(payload);
      renderReportSections(report, {
        contentId: "live-content",
        metaId: "live-meta",
        emptyText: "当前没有启用任何实时数据库连接。",
        timePrefix: "最近一次实时分析时间：",
      });
      setBanner("实时分析完成。你可以调整连接参数后再次执行。", "success");
    } catch (error) {
      setBanner(`实时分析失败：${error.message}`, "error");
    }
  });
}

function initialize() {
  const stored = loadStoredConfig();
  populateForm(stored || exampleConfig);
  bindActions();
}

loadDemo()
  .then((data) => {
    initialize();
    renderHero(data);
    renderNav(data.sections || []);
    renderStaticSections(data.sections || []);
  })
  .catch((error) => {
    const content = document.getElementById("content");
    const block = createElement("section", "section");
    block.appendChild(createElement("h3", "", "页面加载失败"));
    block.appendChild(createElement("p", "", `错误信息：${error.message}`));
    content.appendChild(block);
  });
