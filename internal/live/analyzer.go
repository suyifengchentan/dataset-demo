package live

import (
	"context"
	"time"
)

type Analyzer struct{}

func NewAnalyzer() Analyzer {
	return Analyzer{}
}

func (a Analyzer) Analyze(ctx context.Context, req AnalyzeRequest) AnalyzeResponse {
	timeout := 6 * time.Second
	if req.TimeoutSec > 0 && req.TimeoutSec <= 30 {
		timeout = time.Duration(req.TimeoutSec) * time.Second
	}

	sections := make([]LiveSection, 0, 4)

	if cfg := req.MySQL; cfg != nil && cfg.Enabled && cfg.DSN != "" {
		sections = append(sections, analyzeMySQL(ctx, *cfg, timeout))
	}
	if cfg := req.Postgres; cfg != nil && cfg.Enabled && cfg.DSN != "" {
		sections = append(sections, analyzePostgres(ctx, *cfg, timeout))
	}
	if cfg := req.Redis; cfg != nil && cfg.Enabled && cfg.Addr != "" {
		sections = append(sections, analyzeRedis(ctx, *cfg, timeout))
	}
	if cfg := req.MongoDB; cfg != nil && cfg.Enabled && cfg.URI != "" {
		sections = append(sections, analyzeMongo(ctx, *cfg, timeout))
	}

	return AnalyzeResponse{
		GeneratedAt: time.Now().Format(time.RFC3339),
		Sections:    sections,
	}
}
