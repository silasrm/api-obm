package importer

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ValidationReport struct {
	TableCounts map[string]int64
	FKWarnings  []string
	Errors      []string
}

type Validator struct {
	pool *pgxpool.Pool
}

func NewValidator(pool *pgxpool.Pool) *Validator {
	return &Validator{pool: pool}
}

func (v *Validator) Validate(ctx context.Context) (*ValidationReport, error) {
	report := &ValidationReport{
		TableCounts: make(map[string]int64),
	}

	rows, err := v.pool.Query(ctx, "SELECT tablename FROM pg_tables WHERE schemaname = 'public' ORDER BY tablename")
	if err != nil {
		return nil, fmt.Errorf("querying tables: %w", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("scanning table name: %w", err)
		}
		tables = append(tables, name)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating tables: %w", err)
	}

	fmt.Println("Validação pós-importação:")
	for _, table := range tables {
		var count int64
		err := v.pool.QueryRow(ctx, fmt.Sprintf("SELECT COUNT(*) FROM %s", table)).Scan(&count)
		if err != nil {
			report.Errors = append(report.Errors, fmt.Sprintf("%s: %v", table, err))
			fmt.Printf("%s: ERRO (%v)\n", table, err)
			continue
		}
		report.TableCounts[table] = count
		fmt.Printf("%s: %d registros\n", table, count)
	}

	fmt.Println("Verificação de integridade referencial:")

	type fkCheck struct {
		label        string
		query        string
		requires     string
		warningFmt   string
	}
	checks := []fkCheck{
		{
			label:    "VMP->VTM",
			query:    "SELECT COUNT(*) FROM tb_vmp v WHERE v.CO_VTMID NOT IN (SELECT CO_SEQ_ID FROM tb_vtm) AND v.CO_VTMID != 0",
			requires: "tb_vtm",
			warningFmt: "VMP→VTM: %d órfãos ⚠",
		},
		{
			label:    "AMP->VMP",
			query:    "SELECT COUNT(*) FROM tb_amp a WHERE a.CO_VPID NOT IN (SELECT CO_SEQ_ID FROM tb_vmp) AND a.CO_VPID != 0",
			requires: "tb_vmp",
			warningFmt: "AMP→VMP: %d órfãos ⚠",
		},
		{
			label:    "AMP->Supplier",
			query:    "SELECT COUNT(*) FROM tb_amp a WHERE a.CO_SUPPCD NOT IN (SELECT CO_SEQ_ID FROM td_supplier) AND a.CO_SUPPCD != 0",
			requires: "td_supplier",
			warningFmt: "AMP→Fornecedor: %d órfãos ⚠",
		},
	}

	for _, check := range checks {
		if _, exists := report.TableCounts[check.requires]; !exists {
			fmt.Printf("%s: IGNORADO (tabela %s não importada)\n", check.label, check.requires)
			continue
		}

		var orphanCount int64
		err := v.pool.QueryRow(ctx, check.query).Scan(&orphanCount)
		if err != nil {
			report.Errors = append(report.Errors, fmt.Sprintf("%s: %v", check.label, err))
			fmt.Printf("%s: ERRO (%v)\n", check.label, err)
			continue
		}

		if orphanCount > 0 {
			warning := fmt.Sprintf(check.warningFmt, orphanCount)
			report.FKWarnings = append(report.FKWarnings, warning)
			fmt.Printf("%s: %d órfãos ⚠\n", check.label, orphanCount)
		} else {
			fmt.Printf("%s: OK (%d órfãos)\n", check.label, orphanCount)
		}
	}

	var totalRecords int64
	for _, c := range report.TableCounts {
		totalRecords += c
	}

	fmt.Printf("Resumo: %d tabelas, %d registros, %d erros, %d avisos\n",
		len(report.TableCounts), totalRecords, len(report.Errors), len(report.FKWarnings))

	return report, nil
}
