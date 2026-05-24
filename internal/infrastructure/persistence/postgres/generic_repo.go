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

type GenericRepo struct {
	pool    *pgxpool.Pool
	table   string
	columns string
}

func NewVTMRepo(pool *pgxpool.Pool) *GenericRepo {
	return &GenericRepo{pool: pool, table: "tb_vtm",
		columns: "CO_SEQ_ID, NU_VTMID, NO_NM, NO_ABBREVNM, AU_AUDIT_CREATE_DATETIME, AU_AUDIT_CHANGE_DATETIME, AU_AUDIT_CREATE_USERNAME, AU_AUDIT_CHANGE_USERNAME, ST_REGISTRO_ATIVO, ST_CHANGE_APPROVE"}
}

func NewVMPPRepo(pool *pgxpool.Pool) *GenericRepo {
	return &GenericRepo{pool: pool, table: "tb_vmpp",
		columns: "CO_SEQ_ID, NU_VPPID, CO_VPID, NO_NM, ST_COMBPACKCD, QT_QTYVAL_FORNECIMENTO, CO_QTY_UOMCD_FORNECIMENTO, QT_QTYVAL, CO_QTY_UOMCD, ST_REGISTRO_ATIVO, ST_CHANGE_APPROVE"}
}

func NewAMPPRepo(pool *pgxpool.Pool) *GenericRepo {
	return &GenericRepo{pool: pool, table: "tb_ampp",
		columns: "CO_SEQ_ID, NU_APPID, CO_VPPID, NO_NM, CO_APID, DS_SUBP, ST_COMBPACKCD, CO_LEGAL_CATCD, CO_DISCCD, DT_DISCDT, NU_SANREG, DT_REGPUBLIC, ST_HOSP, ST_REGISTRO_ATIVO, CO_PRIMARYPCK, CO_SECONDARYPCK, CO_INDFARMPOPCD, NU_EAN13a, NU_EAN13b, NU_EAN13c, ST_CHANGE_APPROVE"}
}

func NewDCBRepo(pool *pgxpool.Pool) *GenericRepo {
	return &GenericRepo{pool: pool, table: "tb_dcb",
		columns: "CO_SEQ_ID, NU_DCB, DS_DCB, CO_CAS, DS_CLASSIF, ST_REGISTRO_ATIVO, ST_CHANGE_APPROVE"}
}

func NewSupplierRepo(pool *pgxpool.Pool) *GenericRepo {
	return &GenericRepo{pool: pool, table: "td_supplier",
		columns: "CO_SEQ_ID, NU_CD, DT_CDDT, NO_DESCR, NU_CNPJ, NU_NAUTH, CO_COUNTRYCD, ST_REGISTRO_ATIVO, ST_CHANGE_APPROVE"}
}

func NewIngredientRepo(pool *pgxpool.Pool) *GenericRepo {
	return &GenericRepo{pool: pool, table: "td_ingredient_substances",
		columns: "CO_SEQ_ID, NU_ISID, DT_ISIDDT, NO_NM, NO_NM_PT_BR, CO_ISRCID, CO_DCBCD, CO_NCAS, ST_REGISTRO_ATIVO, ST_CHANGE_APPROVE"}
}

func (r *GenericRepo) GetByID(ctx context.Context, id int64) (interface{}, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE CO_SEQ_ID = $1 AND ST_REGISTRO_ATIVO = 'ACTIVE'", r.columns, r.table)
	row := r.pool.QueryRow(ctx, query, id)
	switch r.table {
	case "tb_vtm":
		var v entity.VTM
		if err := row.Scan(&v.COSeqID, &v.NUVtmID, &v.NONm, &v.NOAbbrevNm,
			&v.AUAuditCreateDatetime, &v.AUAuditChangeDatetime,
			&v.AUAuditCreateUsername, &v.AUAuditChangeUsername,
			&v.STRegistroAtivo, &v.STChangeApprove); err != nil {
			if err == pgx.ErrNoRows {
				return nil, nil
			}
			return nil, err
		}
		return &v, nil
	case "tb_vmpp":
		var v entity.VMPP
		if err := row.Scan(&v.COSeqID, &v.NUVppID, &v.COVpID, &v.NONm, &v.STCombPackCd,
			&v.QTQtyvalFornecimento, &v.COQtyUomCdFornecimento, &v.QTQtyval, &v.COQtyUomCd,
			&v.STRegistroAtivo, &v.STChangeApprove); err != nil {
			if err == pgx.ErrNoRows {
				return nil, nil
			}
			return nil, err
		}
		return &v, nil
	case "tb_ampp":
		var v entity.AMPP
		if err := row.Scan(&v.COSeqID, &v.NUAppID, &v.COVppID, &v.NONm, &v.COApID, &v.DSSubP, &v.STCombPackCd,
			&v.COLegalCatCd, &v.CODiscCd, &v.DTDiscDt, &v.NUSanReg, &v.DTRegPublic, &v.STHosp,
			&v.STRegistroAtivo, &v.COPrimaryPck, &v.COSecondaryPck, &v.COIndFarmPopCd,
			&v.NUEAN13a, &v.NUEAN13b, &v.NUEAN13c, &v.STChangeApprove); err != nil {
			if err == pgx.ErrNoRows {
				return nil, nil
			}
			return nil, err
		}
		return &v, nil
	case "tb_dcb":
		var v entity.DCB
		if err := row.Scan(&v.COSeqID, &v.NUDcb, &v.DSDcb, &v.COCas, &v.DSClassif, &v.STRegistroAtivo, &v.STChangeApprove); err != nil {
			if err == pgx.ErrNoRows {
				return nil, nil
			}
			return nil, err
		}
		return &v, nil
	case "td_supplier":
		var v entity.Supplier
		if err := row.Scan(&v.COSeqID, &v.NUCd, &v.DTCdDt, &v.NODescr, &v.NUCnpj, &v.NUNAuth, &v.COCountryCd,
			&v.STRegistroAtivo, &v.STChangeApprove); err != nil {
			if err == pgx.ErrNoRows {
				return nil, nil
			}
			return nil, err
		}
		return &v, nil
	case "td_ingredient_substances":
		var v entity.IngredientSubstance
		if err := row.Scan(&v.COSeqID, &v.NUIsID, &v.DTIsIDDt, &v.NONm, &v.NONmPtBr,
			&v.COIsrCd, &v.CODcbCd, &v.CONCas, &v.STRegistroAtivo, &v.STChangeApprove); err != nil {
			if err == pgx.ErrNoRows {
				return nil, nil
			}
			return nil, err
		}
		return &v, nil
	default:
		return nil, fmt.Errorf("unknown table: %s", r.table)
	}
}

func (r *GenericRepo) List(ctx context.Context, filters repository.FilterParams) (interface{}, error) {
	limit := filters.Limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	var where []string
	var args []interface{}
	argIdx := 1

	where = append(where, "ST_REGISTRO_ATIVO = 'ACTIVE'")

	if filters.Nome != "" {
		nameCol := r.getNameColumn()
		where = append(where, fmt.Sprintf("%s ILIKE $%d", nameCol, argIdx))
		args = append(args, "%"+filters.Nome+"%")
		argIdx++
	}
	if filters.Codigo != "" {
		codeCol := r.getCodeColumn()
		where = append(where, fmt.Sprintf("%s = $%d", codeCol, argIdx))
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

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s", r.table, whereClause)
	var total int64
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, fmt.Errorf("count %s: %w", r.table, err)
	}

	dataQuery := fmt.Sprintf(
		"SELECT %s FROM %s WHERE %s ORDER BY CO_SEQ_ID LIMIT $%d OFFSET $%d",
		r.columns, r.table, whereClause, argIdx, argIdx+1,
	)
	args = append(args, limit+1, offset)

	rows, err := r.pool.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("query %s list: %w", r.table, err)
	}
	defer rows.Close()

	return r.scanRows(rows, limit, offset, total)
}

func (r *GenericRepo) getNameColumn() string {
	switch r.table {
	case "tb_dcb":
		return "DS_DCB"
	case "td_supplier":
		return "NO_DESCR"
	default:
		return "NO_NM"
	}
}

func (r *GenericRepo) getCodeColumn() string {
	switch r.table {
	case "tb_vtm":
		return "NU_VTMID"
	case "tb_vmpp":
		return "NU_VPPID"
	case "tb_ampp":
		return "NU_APPID"
	case "tb_dcb":
		return "NU_DCB"
	case "td_supplier":
		return "NU_CD"
	case "td_ingredient_substances":
		return "NU_ISID"
	default:
		return "NU_CD"
	}
}

func (r *GenericRepo) scanRows(rows pgx.Rows, limit int, offset int64, total int64) (interface{}, error) {
	switch r.table {
	case "tb_vtm":
		return r.scanVTMRows(rows, limit, offset, total)
	case "tb_vmpp":
		return r.scanVMPPRows(rows, limit, offset, total)
	case "tb_ampp":
		return r.scanAMPPRows(rows, limit, offset, total)
	case "tb_dcb":
		return r.scanDCBRows(rows, limit, offset, total)
	case "td_supplier":
		return r.scanSupplierRows(rows, limit, offset, total)
	case "td_ingredient_substances":
		return r.scanIngredientRows(rows, limit, offset, total)
	default:
		return nil, fmt.Errorf("unknown table: %s", r.table)
	}
}

func (r *GenericRepo) scanVTMRows(rows pgx.Rows, limit int, offset int64, total int64) (*entity.CursorPage[entity.VTM], error) {
	var items []entity.VTM
	for rows.Next() {
		var v entity.VTM
		if err := rows.Scan(&v.COSeqID, &v.NUVtmID, &v.NONm, &v.NOAbbrevNm,
			&v.AUAuditCreateDatetime, &v.AUAuditChangeDatetime,
			&v.AUAuditCreateUsername, &v.AUAuditChangeUsername,
			&v.STRegistroAtivo, &v.STChangeApprove); err != nil {
			return nil, err
		}
		items = append(items, v)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	var nextCursor string
	if len(items) > limit {
		items = items[:limit]
		nextCursor = encodeCursor(offset + int64(limit))
	}
	return &entity.CursorPage[entity.VTM]{Items: items, Cursor: nextCursor, Limit: limit, Total: total}, nil
}

func (r *GenericRepo) scanVMPPRows(rows pgx.Rows, limit int, offset int64, total int64) (*entity.CursorPage[entity.VMPP], error) {
	var items []entity.VMPP
	for rows.Next() {
		var v entity.VMPP
		if err := rows.Scan(&v.COSeqID, &v.NUVppID, &v.COVpID, &v.NONm, &v.STCombPackCd,
			&v.QTQtyvalFornecimento, &v.COQtyUomCdFornecimento, &v.QTQtyval, &v.COQtyUomCd,
			&v.STRegistroAtivo, &v.STChangeApprove); err != nil {
			return nil, err
		}
		items = append(items, v)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	var nextCursor string
	if len(items) > limit {
		items = items[:limit]
		nextCursor = encodeCursor(offset + int64(limit))
	}
	return &entity.CursorPage[entity.VMPP]{Items: items, Cursor: nextCursor, Limit: limit, Total: total}, nil
}

func (r *GenericRepo) scanAMPPRows(rows pgx.Rows, limit int, offset int64, total int64) (*entity.CursorPage[entity.AMPP], error) {
	var items []entity.AMPP
	for rows.Next() {
		var v entity.AMPP
		if err := rows.Scan(&v.COSeqID, &v.NUAppID, &v.COVppID, &v.NONm, &v.COApID, &v.DSSubP, &v.STCombPackCd,
			&v.COLegalCatCd, &v.CODiscCd, &v.DTDiscDt, &v.NUSanReg, &v.DTRegPublic, &v.STHosp,
			&v.STRegistroAtivo, &v.COPrimaryPck, &v.COSecondaryPck, &v.COIndFarmPopCd,
			&v.NUEAN13a, &v.NUEAN13b, &v.NUEAN13c, &v.STChangeApprove); err != nil {
			return nil, err
		}
		items = append(items, v)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	var nextCursor string
	if len(items) > limit {
		items = items[:limit]
		nextCursor = encodeCursor(offset + int64(limit))
	}
	return &entity.CursorPage[entity.AMPP]{Items: items, Cursor: nextCursor, Limit: limit, Total: total}, nil
}

func (r *GenericRepo) scanDCBRows(rows pgx.Rows, limit int, offset int64, total int64) (*entity.CursorPage[entity.DCB], error) {
	var items []entity.DCB
	for rows.Next() {
		var v entity.DCB
		if err := rows.Scan(&v.COSeqID, &v.NUDcb, &v.DSDcb, &v.COCas, &v.DSClassif, &v.STRegistroAtivo, &v.STChangeApprove); err != nil {
			return nil, err
		}
		items = append(items, v)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	var nextCursor string
	if len(items) > limit {
		items = items[:limit]
		nextCursor = encodeCursor(offset + int64(limit))
	}
	return &entity.CursorPage[entity.DCB]{Items: items, Cursor: nextCursor, Limit: limit, Total: total}, nil
}

func (r *GenericRepo) scanSupplierRows(rows pgx.Rows, limit int, offset int64, total int64) (*entity.CursorPage[entity.Supplier], error) {
	var items []entity.Supplier
	for rows.Next() {
		var v entity.Supplier
		if err := rows.Scan(&v.COSeqID, &v.NUCd, &v.DTCdDt, &v.NODescr, &v.NUCnpj, &v.NUNAuth, &v.COCountryCd,
			&v.STRegistroAtivo, &v.STChangeApprove); err != nil {
			return nil, err
		}
		items = append(items, v)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	var nextCursor string
	if len(items) > limit {
		items = items[:limit]
		nextCursor = encodeCursor(offset + int64(limit))
	}
	return &entity.CursorPage[entity.Supplier]{Items: items, Cursor: nextCursor, Limit: limit, Total: total}, nil
}

func (r *GenericRepo) scanIngredientRows(rows pgx.Rows, limit int, offset int64, total int64) (*entity.CursorPage[entity.IngredientSubstance], error) {
	var items []entity.IngredientSubstance
	for rows.Next() {
		var v entity.IngredientSubstance
		if err := rows.Scan(&v.COSeqID, &v.NUIsID, &v.DTIsIDDt, &v.NONm, &v.NONmPtBr,
			&v.COIsrCd, &v.CODcbCd, &v.CONCas, &v.STRegistroAtivo, &v.STChangeApprove); err != nil {
			return nil, err
		}
		items = append(items, v)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	var nextCursor string
	if len(items) > limit {
		items = items[:limit]
		nextCursor = encodeCursor(offset + int64(limit))
	}
	return &entity.CursorPage[entity.IngredientSubstance]{Items: items, Cursor: nextCursor, Limit: limit, Total: total}, nil
}
