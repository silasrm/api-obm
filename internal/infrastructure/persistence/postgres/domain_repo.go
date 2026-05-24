package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/silasrm/api-obm/internal/domain/entity"
	"github.com/silasrm/api-obm/internal/domain/repository"
)

type DomainRepo struct {
	pool *pgxpool.Pool
}

var domainTables = map[string]string{
	"anvs-class":             "td_anvs_class_br",
	"atc-class":              "td_atc_class_br",
	"catmat":                 "td_catmat_br",
	"form":                   "td_form",
	"route":                  "td_route",
	"control-drug":           "td_control_drug_category",
	"unit-of-measure":        "td_unit_of_measure",
	"package":                "td_package",
	"legal-category":         "td_legal_category",
	"country":                "td_country",
	"flavour":                "td_flavour",
	"basis-of-name":          "td_basis_of_name",
	"basis-of-strnth":        "td_basis_of_strnth",
	"brimunologico":          "td_brimunologico",
	"df-indicator":           "td_df_indicator",
	"discontinued-ind":       "td_discontinued_ind",
	"pres-status":            "td_virtual_product_pres_status",
	"non-avail":              "td_virtual_product_non_avail",
	"avail-restriction":      "td_availability_restriction",
	"licensing-authority":    "td_licensing_authority",
	"lic-auth-change-reason": "td_lic_auth_change_reason",
	"med-class":              "td_med_class_br",
	"name-change-reason":     "td_name_change_reason",
	"ont-form-route":         "td_ont_form_route",
	"preserv-cond":           "td_preserv_cond_br",
	"rename-comp":            "td_rename_comp_br",
	"ingredient-source":      "td_ingredient_source_br",
	"healthcare-prof":        "td_healthcare_prof_br",
	"indicacao-farmpop":      "td_indicacao_farmpop_br",
	"monitoring-reason":      "td_monitoring_reason_br",
	"phpid":                  "td_phpid",
	"local-aplicacao":        "td_local_aplicacao",
}

func NewDomainRepo(pool *pgxpool.Pool) *DomainRepo {
	return &DomainRepo{pool: pool}
}

func (r *DomainRepo) GetByID(ctx context.Context, domainType string, id int64) (*entity.Domain, error) {
	table, ok := domainTables[domainType]
	if !ok {
		return nil, fmt.Errorf("unknown domain type: %s", domainType)
	}
	query := fmt.Sprintf("SELECT %s FROM %s WHERE CO_SEQ_ID = $1 AND ST_REGISTRO_ATIVO = 'ACTIVE'", domainColumns, table)
	var d entity.Domain
	err := r.pool.QueryRow(ctx, query, id).Scan(domainScanFields(&d)...)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query domain %s: %w", domainType, err)
	}
	return &d, nil
}

func (r *DomainRepo) List(ctx context.Context, domainType string, filters repository.FilterParams) (*entity.CursorPage[entity.Domain], error) {
	table, ok := domainTables[domainType]
	if !ok {
		return nil, fmt.Errorf("unknown domain type: %s", domainType)
	}

	limit := filters.Limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	var where []string
	var args []interface{}
	argIdx := 1

	where = append(where, "ST_REGISTRO_ATIVO = 'ACTIVE'")

	if filters.Nome != "" {
		where = append(where, fmt.Sprintf("(NO_DESCR ILIKE $%d OR NO_DESCR_PT_BR ILIKE $%d)", argIdx, argIdx))
		args = append(args, "%"+filters.Nome+"%")
		argIdx++
	}
	if filters.Codigo != "" {
		where = append(where, fmt.Sprintf("NU_CD = $%d", argIdx))
		args = append(args, filters.Codigo)
		argIdx++
	}
	if filters.Ativo != nil {
		if *filters.Ativo {
			where = append(where, "ST_REGISTRO_ATIVO = 'ACTIVE'")
		} else {
			where = append(where, "ST_REGISTRO_ATIVO = 'INACTIVE'")
		}
	}

	offset := decodeCursor(filters.Cursor)
	whereClause := strings.Join(where, " AND ")

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s", table, whereClause)
	var total int64
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, fmt.Errorf("count domain %s: %w", domainType, err)
	}

	dataQuery := fmt.Sprintf(
		"SELECT %s FROM %s WHERE %s ORDER BY CO_SEQ_ID LIMIT $%d OFFSET $%d",
		domainColumns, table, whereClause, argIdx, argIdx+1,
	)
	args = append(args, limit+1, offset)

	rows, err := r.pool.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("query domain %s list: %w", domainType, err)
	}
	defer rows.Close()

	items, err := scanDomainRows(rows)
	if err != nil {
		return nil, err
	}

	var nextCursor string
	if len(items) > limit {
		items = items[:limit]
		nextCursor = encodeCursor(offset + int64(limit))
	}

	return &entity.CursorPage[entity.Domain]{
		Items:  items,
		Cursor: nextCursor,
		Limit:  limit,
		Total:  total,
	}, nil
}
