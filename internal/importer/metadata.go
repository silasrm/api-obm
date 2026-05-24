package importer

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type MetadataManager struct {
	pool *pgxpool.Pool
}

func NewMetadataManager(pool *pgxpool.Pool) *MetadataManager {
	return &MetadataManager{pool: pool}
}

func (m *MetadataManager) EnsureTable(ctx context.Context) error {
	_, err := m.pool.Exec(ctx, `CREATE TABLE IF NOT EXISTS obm_metadata (key VARCHAR(100) PRIMARY KEY, value TEXT NOT NULL, updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW())`)
	return err
}

func (m *MetadataManager) Set(ctx context.Context, key, value string) error {
	_, err := m.pool.Exec(ctx, `INSERT INTO obm_metadata (key, value, updated_at) VALUES ($1, $2, NOW()) ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW()`, key, value)
	return err
}

func (m *MetadataManager) Get(ctx context.Context, key string) (string, error) {
	var value string
	err := m.pool.QueryRow(ctx, `SELECT value FROM obm_metadata WHERE key = $1`, key).Scan(&value)
	if err != nil {
		return "", nil
	}
	return value, nil
}

func (m *MetadataManager) RecordImport(ctx context.Context, version, sourceFile string, tableCounts map[string]int64) error {
	if err := m.EnsureTable(ctx); err != nil {
		return err
	}

	if err := m.Set(ctx, "obm_version", version); err != nil {
		return err
	}

	if err := m.Set(ctx, "import_date", time.Now().Format(time.RFC3339)); err != nil {
		return err
	}

	if err := m.Set(ctx, "source_file", sourceFile); err != nil {
		return err
	}

	parts := ""
	first := true
	for table, count := range tableCounts {
		if first {
			parts += fmt.Sprintf("%s:%d", table, count)
			first = false
		} else {
			parts += fmt.Sprintf(", %s:%d", table, count)
		}
	}
	recordCounts := fmt.Sprintf("{%s}", parts)

	if err := m.Set(ctx, "record_counts", recordCounts); err != nil {
		return err
	}

	if err := m.Set(ctx, "tables_imported", fmt.Sprintf("%d", len(tableCounts))); err != nil {
		return err
	}

	return nil
}
