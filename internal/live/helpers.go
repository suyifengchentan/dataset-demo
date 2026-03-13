package live

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

func queryRows(ctx context.Context, db *sql.DB, query string, args ...any) (DataTable, error) {
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return DataTable{}, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return DataTable{}, err
	}

	table := DataTable{Columns: columns}
	for rows.Next() {
		values := make([]any, len(columns))
		dest := make([]any, len(columns))
		for i := range values {
			dest[i] = &values[i]
		}
		if err := rows.Scan(dest...); err != nil {
			return DataTable{}, err
		}

		row := make([]string, len(columns))
		for i, value := range values {
			row[i] = stringify(value)
		}
		table.Rows = append(table.Rows, row)
	}

	return table, rows.Err()
}

func queryString(ctx context.Context, db *sql.DB, query string, args ...any) (string, error) {
	var value sql.NullString
	if err := db.QueryRowContext(ctx, query, args...).Scan(&value); err != nil {
		return "", err
	}
	return value.String, nil
}

func stringify(value any) string {
	switch v := value.(type) {
	case nil:
		return ""
	case []byte:
		return string(v)
	default:
		return fmt.Sprint(v)
	}
}

func compactJSON(value any) string {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return fmt.Sprint(value)
	}
	return string(data)
}

func prettyMapLines(values map[string]string, keys ...string) string {
	lines := make([]string, 0, len(keys))
	for _, key := range keys {
		if value, ok := values[key]; ok && value != "" {
			lines = append(lines, fmt.Sprintf("%s = %s", key, value))
		}
	}
	return strings.Join(lines, "\n")
}

func parseRedisInfo(raw string) map[string]string {
	info := map[string]string{}
	for _, line := range strings.Split(raw, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			info[parts[0]] = parts[1]
		}
	}
	return info
}

func sortedKeys(values map[string]string) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
