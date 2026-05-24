package postgres

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/silasrm/api-obm/internal/domain/entity"
	"github.com/silasrm/api-obm/internal/domain/repository"
)

type VMPRepo struct {
	pool *pgxpool.Pool
}

func NewVMPRepo(pool *pgxpool.Pool) *VMPRepo {
	return &VMPRepo{pool: pool}
}

const vmpColumns = `CO_SEQ_ID, NU_VPID, DT_VPIDDT, CO_VPIDPREV, ST_INVALID,
	CO_VTMID, ST_COMBPRODCD, NO_NM, NU_SNOMED, CO_BASISCD,
	DT_NMDT, NO_NMPREV, CO_BASIS_PREVCD, CO_NMCHANGECD, NO_ABBREVNM,
	ST_SUG_F, ST_GLU_F, ST_PRES_F, ST_CFC_FREE, CO_PRES_STATCD,
	CO_NON_AVAILCD, DT_NON_AVAILDT, CO_DF_INDCD, QT_UDFS, CO_UDFS_UOMCD,
	CO_UNIT_DOSE_UOMCD, ST_FTN, CO_BRIMUNOCD, ST_HORUS, DS_NOTIFICATION,
	CO_TARGETCD, CO_ANVSCLSCD, NU_VERSION, ST_REGISTRO_ATIVO, ST_CHANGE_APPROVE`

func scanVMP(row pgx.Row) (*entity.VMP, error) {
	var v entity.VMP
	err := row.Scan(
		&v.COSeqID, &v.NUVpID, &v.DTVpIDDt, &v.COVpIDPrev, &v.STInvalid,
		&v.COVtmID, &v.STCombProdCd, &v.NONm, &v.NUSnomed, &v.COBasisCd,
		&v.DTNmDt, &v.NONmPrev, &v.COBasisPrevCd, &v.CONmChangeCd, &v.NOAbbrevNm,
		&v.STSugF, &v.STGluF, &v.STPresF, &v.STCfcFree, &v.COPresStatCd,
		&v.CONonAvailCd, &v.DTNonAvailDt, &v.CODFIndCd, &v.QTUdfs, &v.COUdfsUomCd,
		&v.COUnitDoseUomCd, &v.STFtn, &v.COBrimunoCd, &v.STHorus, &v.DSNotification,
		&v.COTargetCd, &v.COAnvsClsCd, &v.NUVersion, &v.STRegistroAtivo, &v.STChangeApprove,
	)
	return &v, err
}

func scanVMPRows(rows pgx.Rows) ([]entity.VMP, error) {
	var items []entity.VMP
	for rows.Next() {
		var v entity.VMP
		if err := rows.Scan(
			&v.COSeqID, &v.NUVpID, &v.DTVpIDDt, &v.COVpIDPrev, &v.STInvalid,
			&v.COVtmID, &v.STCombProdCd, &v.NONm, &v.NUSnomed, &v.COBasisCd,
			&v.DTNmDt, &v.NONmPrev, &v.COBasisPrevCd, &v.CONmChangeCd, &v.NOAbbrevNm,
			&v.STSugF, &v.STGluF, &v.STPresF, &v.STCfcFree, &v.COPresStatCd,
			&v.CONonAvailCd, &v.DTNonAvailDt, &v.CODFIndCd, &v.QTUdfs, &v.COUdfsUomCd,
			&v.COUnitDoseUomCd, &v.STFtn, &v.COBrimunoCd, &v.STHorus, &v.DSNotification,
			&v.COTargetCd, &v.COAnvsClsCd, &v.NUVersion, &v.STRegistroAtivo, &v.STChangeApprove,
		); err != nil {
			return nil, err
		}
		items = append(items, v)
	}
	return items, rows.Err()
}

func (r *VMPRepo) GetByID(ctx context.Context, id int64) (*entity.VMP, error) {
	query := fmt.Sprintf("SELECT %s FROM tb_vmp WHERE CO_SEQ_ID = $1 AND ST_REGISTRO_ATIVO = 'ACTIVE'", vmpColumns)
	v, err := scanVMP(r.pool.QueryRow(ctx, query, id))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query vmp: %w", err)
	}
	return v, nil
}

func (r *VMPRepo) GetDetailByID(ctx context.Context, id int64) (*entity.VMPDetail, error) {
	vmp, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if vmp == nil {
		return nil, nil
	}
	detail := &entity.VMPDetail{VMP: *vmp}

	detail.VTM, err = r.getVTM(ctx, vmp.COVtmID)
	if err != nil {
		return nil, err
	}
	detail.BasisOfName, err = r.getDomain(ctx, "td_basis_of_name", vmp.COBasisCd)
	if err != nil {
		return nil, err
	}
	detail.PresStatus, err = r.getDomain(ctx, "td_virtual_product_pres_status", vmp.COPresStatCd)
	if err != nil {
		return nil, err
	}
	detail.DFIndicator, err = r.getDomain(ctx, "td_df_indicator", vmp.CODFIndCd)
	if err != nil {
		return nil, err
	}
	detail.BrImunologico, err = r.getDomain(ctx, "td_brimunologico", vmp.COBrimunoCd)
	if err != nil {
		return nil, err
	}
	detail.AnvsClass, err = r.getDomain(ctx, "td_anvs_class_br", vmp.COAnvsClsCd)
	if err != nil {
		return nil, err
	}
	detail.Forms, err = r.getVMPRelatedDomains(ctx, "rl_vmp_form", "CO_FORMCD", "td_form", id)
	if err != nil {
		return nil, err
	}
	detail.Routes, err = r.getVMPRelatedDomains(ctx, "rl_vmp_route", "CO_ROUTECD", "td_route", id)
	if err != nil {
		return nil, err
	}
	detail.ATCClasses, err = r.getVMPRelatedDomains(ctx, "rl_vmp_atc_class_br", "CO_ATCCLSCD", "td_atc_class_br", id)
	if err != nil {
		return nil, err
	}
	detail.Catmats, err = r.getVMPRelatedDomains(ctx, "rl_vmp_catmat_br", "CO_CATMATCD", "td_catmat_br", id)
	if err != nil {
		return nil, err
	}
	detail.ControlDrugInfos, err = r.getVMPRelatedDomains(ctx, "rl_vmp_control_drug_info", "CO_CATCD", "td_control_drug_category", id)
	if err != nil {
		return nil, err
	}
	detail.RenamesBR, err = r.getVMPRelatedDomains(ctx, "rl_vmp_rename_br", "CO_RENAMEID", "td_rename_comp_br", id)
	if err != nil {
		return nil, err
	}
	detail.LocalAplicacao, err = r.getVMPRelatedDomains(ctx, "rl_vmp_local_aplicacao", "CO_LOCAL_APLICACAO_CD", "td_local_aplicacao", id)
	if err != nil {
		return nil, err
	}
	detail.Ingredients, err = r.getVMPIngredients(ctx, id)
	if err != nil {
		return nil, err
	}
	return detail, nil
}

func (r *VMPRepo) getVTM(ctx context.Context, vtmID int64) (*entity.VTM, error) {
	const query = `SELECT CO_SEQ_ID, NU_VTMID, NO_NM, NO_ABBREVNM,
		AU_AUDIT_CREATE_DATETIME, AU_AUDIT_CHANGE_DATETIME,
		AU_AUDIT_CREATE_USERNAME, AU_AUDIT_CHANGE_USERNAME,
		ST_REGISTRO_ATIVO, ST_CHANGE_APPROVE
		FROM tb_vtm WHERE CO_SEQ_ID = $1 AND ST_REGISTRO_ATIVO = 'ACTIVE'`
	var v entity.VTM
	err := r.pool.QueryRow(ctx, query, vtmID).Scan(
		&v.COSeqID, &v.NUVtmID, &v.NONm, &v.NOAbbrevNm,
		&v.AUAuditCreateDatetime, &v.AUAuditChangeDatetime,
		&v.AUAuditCreateUsername, &v.AUAuditChangeUsername,
		&v.STRegistroAtivo, &v.STChangeApprove,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query vtm: %w", err)
	}
	return &v, nil
}

func (r *VMPRepo) getDomain(ctx context.Context, table string, code sql.NullInt64) (*entity.Domain, error) {
	if !code.Valid {
		return nil, nil
	}
	return r.getDomainByID(ctx, table, code.Int64)
}

func (r *VMPRepo) getDomainByID(ctx context.Context, table string, id int64) (*entity.Domain, error) {
	query := fmt.Sprintf(`SELECT %s FROM %s WHERE CO_SEQ_ID = $1 AND ST_REGISTRO_ATIVO = 'ACTIVE'`, domainColumns, table)
	var d entity.Domain
	err := r.pool.QueryRow(ctx, query, id).Scan(domainScanFields(&d)...)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query domain %s: %w", table, err)
	}
	return &d, nil
}

func (r *VMPRepo) getVMPRelatedDomains(ctx context.Context, rlTable, fkCol, tdTable string, vmpID int64) ([]entity.Domain, error) {
	query := fmt.Sprintf(`SELECT d.%s FROM %s rl JOIN %s d ON d.CO_SEQ_ID = rl.%s
		WHERE rl.CO_VPID = $1 AND rl.ST_REGISTRO_ATIVO = 'ACTIVE' AND d.ST_REGISTRO_ATIVO = 'ACTIVE'`,
		domainColumns, rlTable, tdTable, fkCol)
	rows, err := r.pool.Query(ctx, query, vmpID)
	if err != nil {
		return nil, fmt.Errorf("query %s: %w", rlTable, err)
	}
	defer rows.Close()
	return scanDomainRows(rows)
}

func (r *VMPRepo) getVMPIngredients(ctx context.Context, vmpID int64) ([]entity.VMPIngredient, error) {
	const query = `SELECT rt.CO_SEQ_ID, rt.CO_VPID, rt.CO_ISCD, rt.CO_BS_SUBID, rt.CO_BASIS_STRNTCD,
		rt.QT_STRNT_NMRTR_VAL, rt.CO_STRNT_NMRTR_UOMCD, rt.QT_STRNT_DNMTR_VAL, rt.CO_STRNT_DNMTR_UOMCD,
		is.CO_SEQ_ID, is.NU_ISID, is.DT_ISIDDT, is.NO_NM, is.NO_NM_PT_BR,
		is.CO_ISRCID, is.CO_DCBCD, is.CO_NCAS, is.ST_REGISTRO_ATIVO, is.ST_CHANGE_APPROVE
		FROM rt_vmp_ingredient_subst rt
		LEFT JOIN td_ingredient_substances is ON is.CO_SEQ_ID = rt.CO_ISCD
		WHERE rt.CO_VPID = $1 AND rt.ST_REGISTRO_ATIVO = 'ACTIVE'`
	rows, err := r.pool.Query(ctx, query, vmpID)
	if err != nil {
		return nil, fmt.Errorf("query vmp ingredients: %w", err)
	}
	defer rows.Close()
	var result []entity.VMPIngredient
	for rows.Next() {
		var ing entity.VMPIngredient
		var subst entity.IngredientSubstance
		var substID sql.NullInt64
		err := rows.Scan(
			&ing.COSeqID, &ing.COVmpID, &ing.COIscd, &ing.COBsSubID, &ing.COBasisStrntCd,
			&ing.QTStrntNmrtrVal, &ing.COStrntNmrtrUomCd, &ing.QTStrntDnmtrVal, &ing.COStrntDnmtrUomCd,
			&substID, &subst.NUIsID, &subst.DTIsIDDt, &subst.NONm, &subst.NONmPtBr,
			&subst.COIsrCd, &subst.CODcbCd, &subst.CONCas, &subst.STRegistroAtivo, &subst.STChangeApprove,
		)
		if err != nil {
			return nil, fmt.Errorf("scan vmp ingredient: %w", err)
		}
		if substID.Valid {
			subst.COSeqID = substID.Int64
			ing.IngredientSubstance = &subst
		}
		result = append(result, ing)
	}
	return result, rows.Err()
}

func (r *VMPRepo) List(ctx context.Context, filters repository.FilterParams) (*entity.CursorPage[entity.VMP], error) {
	limit := filters.Limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	var where []string
	var args []interface{}
	argIdx := 1

	where = append(where, "ST_REGISTRO_ATIVO = 'ACTIVE'")

	if filters.Nome != "" {
		where = append(where, fmt.Sprintf("NO_NM ILIKE $%d", argIdx))
		args = append(args, "%"+filters.Nome+"%")
		argIdx++
	}
	if filters.Codigo != "" {
		where = append(where, fmt.Sprintf("NU_VPID = $%d", argIdx))
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

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM tb_vmp WHERE %s", whereClause)
	var total int64
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, fmt.Errorf("count vmp: %w", err)
	}

	dataQuery := fmt.Sprintf(
		"SELECT %s FROM tb_vmp WHERE %s ORDER BY CO_SEQ_ID LIMIT $%d OFFSET $%d",
		vmpColumns, whereClause, argIdx, argIdx+1,
	)
	args = append(args, limit+1, offset)

	rows, err := r.pool.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("query vmp list: %w", err)
	}
	defer rows.Close()

	items, err := scanVMPRows(rows)
	if err != nil {
		return nil, err
	}

	var nextCursor string
	if len(items) > limit {
		items = items[:limit]
		nextCursor = encodeCursor(offset + int64(limit))
	}

	return &entity.CursorPage[entity.VMP]{
		Items:  items,
		Cursor: nextCursor,
		Limit:  limit,
		Total:  total,
	}, nil
}

const domainColumns = `CO_SEQ_ID, NU_CD, NO_DESCR, NO_DESCR_PT_BR, ST_REGISTRO_ATIVO, ST_CHANGE_APPROVE,
	CO_PARENT, ST_HIERARCHY_LEVEL, NO_ABBREVNM, DS_SCOPE, DS_LIMITPREC,
	DS_QTYPRESC, DS_VALPRESC, CO_ALPHA_2, CO_ALPHA_3, CO_LATITUDE, CO_LONGITUDE,
	NO_SIGLA, PNI_DISPONIVEL, NU_CATMAT, DS_CATMAT, NO_HORUS, DT_CDDT, CO_FFACD, NU_VERSION`

func domainScanFields(d *entity.Domain) []interface{} {
	return []interface{}{
		&d.COSeqID, &d.NUCd, &d.NODescr, &d.NODescrPtBr, &d.STRegistroAtivo, &d.STChangeApprove,
		&d.COParent, &d.STHierarchyLevel, &d.NOAbbrevNm, &d.DSScope, &d.DSLimitPresc,
		&d.DSQtypPresc, &d.DSValPresc, &d.COAlpha2, &d.COAlpha3, &d.COLatitude, &d.COLongitude,
		&d.NOSigla, &d.PNIDisponivel, &d.NUCatmat, &d.DSCatmat, &d.NOHorus, &d.DTCdDt, &d.COFFacCd, &d.NUVersion,
	}
}

func scanDomainRows(rows pgx.Rows) ([]entity.Domain, error) {
	var result []entity.Domain
	for rows.Next() {
		var d entity.Domain
		if err := rows.Scan(domainScanFields(&d)...); err != nil {
			return nil, fmt.Errorf("scan domain: %w", err)
		}
		result = append(result, d)
	}
	return result, rows.Err()
}

func decodeCursor(cursor string) int64 {
	if cursor == "" {
		return 0
	}
	b, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return 0
	}
	offset, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		return 0
	}
	return offset
}

func encodeCursor(offset int64) string {
	return base64.StdEncoding.EncodeToString([]byte(strconv.FormatInt(offset, 10)))
}
