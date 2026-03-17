package live

type AnalyzeRequest struct {
	MySQL      *MySQLConfig    `json:"mysql,omitempty"`
	Postgres   *PostgresConfig `json:"postgres,omitempty"`
	Redis      *RedisConfig    `json:"redis,omitempty"`
	MongoDB    *MongoConfig    `json:"mongodb,omitempty"`
	TimeoutSec int             `json:"timeoutSec,omitempty"`
}

type MySQLConfig struct {
	Enabled      bool   `json:"enabled"`
	DSN          string `json:"dsn"`
	Schema       string `json:"schema,omitempty"`
	Table        string `json:"table,omitempty"`
	ExplainQuery string `json:"explainQuery,omitempty"`
}

type PostgresConfig struct {
	Enabled bool   `json:"enabled"`
	DSN     string `json:"dsn"`
	Schema  string `json:"schema,omitempty"`
	Table   string `json:"table,omitempty"`
}

type RedisConfig struct {
	Enabled   bool   `json:"enabled"`
	Addr      string `json:"addr"`
	Username  string `json:"username,omitempty"`
	Password  string `json:"password,omitempty"`
	DB        int    `json:"db,omitempty"`
	KeyPrefix string `json:"keyPrefix,omitempty"`
}

type MongoConfig struct {
	Enabled    bool   `json:"enabled"`
	URI        string `json:"uri"`
	Database   string `json:"database,omitempty"`
	Collection string `json:"collection,omitempty"`
}

type AnalyzeResponse struct {
	GeneratedAt string        `json:"generatedAt"`
	Sections    []LiveSection `json:"sections"`
}

type SimulateResponse struct {
	GeneratedAt string        `json:"generatedAt"`
	Summary     string        `json:"summary"`
	Sections    []LiveSection `json:"sections"`
}

type LiveSection struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	Status    string      `json:"status"`
	Summary   string      `json:"summary"`
	Error     string      `json:"error,omitempty"`
	Metrics   []Metric    `json:"metrics,omitempty"`
	Tables    []DataTable `json:"tables,omitempty"`
	Snippets  []Snippet   `json:"snippets,omitempty"`
	Warnings  []string    `json:"warnings,omitempty"`
	Connected bool        `json:"connected"`
}

type Metric struct {
	Label string `json:"label"`
	Value string `json:"value"`
	Hint  string `json:"hint,omitempty"`
	Tone  string `json:"tone,omitempty"`
}

type DataTable struct {
	Title   string     `json:"title"`
	Columns []string   `json:"columns"`
	Rows    [][]string `json:"rows"`
}

type Snippet struct {
	Title    string `json:"title"`
	Language string `json:"language,omitempty"`
	Content  string `json:"content"`
}
