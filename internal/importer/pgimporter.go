package importer

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PGImporter struct {
	pool *pgxpool.Pool
}

type ImportStats struct {
	TablesCreated int
	RowsInserted  map[string]int64
	Duration      time.Duration
	Errors        []error
}

func NewPGImporter(pool *pgxpool.Pool) *PGImporter {
	return &PGImporter{pool: pool}
}

func (p *PGImporter) Import(ctx context.Context, reader io.Reader) (*ImportStats, error) {
	stats := &ImportStats{
		RowsInserted: make(map[string]int64),
	}
	start := time.Now()

	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 0, 32*1024*1024), 32*1024*1024)

	var builder strings.Builder
	var currentTable string

	for scanner.Scan() {
		line := scanner.Text()
		builder.WriteString(line)
		builder.WriteByte('\n')

		if !strings.HasSuffix(strings.TrimSpace(line), ";") {
			continue
		}

		statement := strings.TrimSpace(builder.String())
		builder.Reset()

		if statement == "" {
			continue
		}

		if strings.HasPrefix(statement, "CREATE TABLE IF NOT EXISTS") {
			parts := strings.Fields(statement)
			if len(parts) >= 5 {
				currentTable = strings.Trim(parts[4], `"`)
				fmt.Printf("Criando tabela: %s\n", currentTable)
				stats.TablesCreated++
			}
		}

		_, err := p.pool.Exec(ctx, statement)
		if err != nil {
			stats.Errors = append(stats.Errors, err)
			continue
		}

		if strings.HasPrefix(statement, "INSERT") {
			if currentTable != "" {
				stats.RowsInserted[currentTable]++
			}
		}
	}

	if err := scanner.Err(); err != nil {
		stats.Errors = append(stats.Errors, err)
	}

	stats.Duration = time.Since(start)

	totalRows := int64(0)
	for _, count := range stats.RowsInserted {
		totalRows += count
	}

	fmt.Printf("Importação concluída em %s: %d tabelas, %d registros, %d erros\n", stats.Duration, stats.TablesCreated, totalRows, len(stats.Errors))

	return stats, nil
}
