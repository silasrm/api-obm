package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/silasrm/api-obm/internal/domain/entity"
	"github.com/silasrm/api-obm/internal/domain/repository"
)

type AMPRepo struct {
	pool *pgxpool.Pool
}

func NewAMPRepo(pool *pgxpool.Pool) *AMPRepo {
	return &AMPRepo{pool: pool}
}

const ampColumns = `CO_SEQ_ID, NU_APID, CO_VPID, ST_COMBPRODCD, NO_NM, DS_DESCR, NO_ABBREVNM,
	CO_SUPPCD, CO_FLAVOURCD, CO_LIC_AUTHCHANGECD, ST_PARALLEL_IMPORT, CO_LIC_AUTHCD,
	CO_AVAIL_RESTRICTCD, CO_MEDCLSCD, CO_MONITORINGREASONCD, CO_ENTERALADMINSTATUSCD,
	DS_ENTERALTUBESADMINOBS, NU_NREG, NU_NPROC, NU_VENCREG, NU_VALIDITY, CO_VALIDITYUNIT,
	ST_REGISTRO_ATIVO, ST_CHANGE_APPROVE`

func scanAMP(row pgx.Row) (*entity.AMP, error) {
	var a entity.AMP
	err := row.Scan(
		&a.COSeqID, &a.NUApID, &a.COVpID, &a.STCombProdCd, &a.NONm, &a.DSDescr, &a.NOAbbrevNm,
		&a.COSuppCd, &a.COFlavourCd, &a.COLicAuthChangeCd, &a.STParallelImport, &a.COLicAuthCd,
		&a.COAvailRestrictCd, &a.COMedClsCd, &a.COMonitoringReasonCd, &a.COEnteralAdminStatusCd,
		&a.DSEnteralTubesAdminObs, &a.NUNReg, &a.NUNProc, &a.NUVencReg, &a.NUValidity, &a.COValidityUnit,
		&a.STRegistroAtivo, &a.STChangeApprove,
	)
	return &a, err
}

func (r *AMPRepo) GetByID(ctx context.Context, id int64) (*entity.AMP, error) {
	query := fmt.Sprintf("SELECT %s FROM tb_amp WHERE CO_SEQ_ID = $1 AND ST_REGISTRO_ATIVO = 'ACTIVE'", ampColumns)
	a, err := scanAMP(r.pool.QueryRow(ctx, query, id))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query amp: %w", err)
	}
	return a, nil
}

func (r *AMPRepo) GetDetailByID(ctx context.Context, id int64) (*entity.AMPDetail, error) {
	amp, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if amp == nil {
		return nil, nil
	}
	detail := &entity.AMPDetail{AMP: *amp}

	vmp, err := scanVMP(r.pool.QueryRow(ctx,
		fmt.Sprintf("SELECT %s FROM tb_vmp WHERE CO_SEQ_ID = $1 AND ST_REGISTRO_ATIVO = 'ACTIVE'", vmpColumns),
		amp.COVpID))
	if err == nil {
		detail.VMP = vmp
	}

	var supp entity.Supplier
	err = r.pool.QueryRow(ctx, `SELECT CO_SEQ_ID, NU_CD, DT_CDDT, NO_DESCR, NU_CNPJ, NU_NAUTH, CO_COUNTRYCD,
		ST_REGISTRO_ATIVO, ST_CHANGE_APPROVE FROM td_supplier WHERE CO_SEQ_ID = $1 AND ST_REGISTRO_ATIVO = 'ACTIVE'`, amp.COSuppCd).Scan(
		&supp.COSeqID, &supp.NUCd, &supp.DTCdDt, &supp.NODescr, &supp.NUCnpj, &supp.NUNAuth, &supp.COCountryCd,
		&supp.STRegistroAtivo, &supp.STChangeApprove,
	)
	if err == nil {
		detail.Supplier = &supp
	}

	detail.Flavour, _ = r.getAMPDomain(ctx, "td_flavour", amp.COFlavourCd)
	detail.LicAuth, _ = r.getAMPDomain(ctx, "td_licensing_authority", amp.COLicAuthCd)
	detail.MedClass, _ = r.getAMPDomain(ctx, "td_med_class_br", amp.COMedClsCd)
	detail.AvailRestriction, _ = r.getAMPDomain(ctx, "td_availability_restriction", amp.COAvailRestrictCd)

	detail.Routes, err = r.getAMPRelatedDomains(ctx, "rl_amp_route", "CO_ROUTECD", "td_route", "CO_APID", id)
	if err != nil {
		return nil, err
	}
	detail.PreservConds, err = r.getAMPRelatedDomains(ctx, "rl_amp_preserv_cond_br", "CO_PRESERVCONDCD", "td_preserv_cond_br", "CO_APID", id)
	if err != nil {
		return nil, err
	}
	detail.Ingredients, err = r.getAMPIngredients(ctx, id)
	if err != nil {
		return nil, err
	}

	return detail, nil
}

func (r *AMPRepo) getAMPDomain(ctx context.Context, table string, code sql.NullInt64) (*entity.Domain, error) {
	if !code.Valid {
		return nil, nil
	}
	query := fmt.Sprintf(`SELECT %s FROM %s WHERE CO_SEQ_ID = $1 AND ST_REGISTRO_ATIVO = 'ACTIVE'`, domainColumns, table)
	var d entity.Domain
	err := r.pool.QueryRow(ctx, query, code.Int64).Scan(domainScanFields(&d)...)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query domain %s: %w", table, err)
	}
	return &d, nil
}

func (r *AMPRepo) getAMPRelatedDomains(ctx context.Context, rlTable, fkCol, tdTable, rlFKCol string, ampID int64) ([]entity.Domain, error) {
	query := fmt.Sprintf(`SELECT d.%s FROM %s rl JOIN %s d ON d.CO_SEQ_ID = rl.%s
		WHERE rl.%s = $1 AND rl.ST_REGISTRO_ATIVO = 'ACTIVE' AND d.ST_REGISTRO_ATIVO = 'ACTIVE'`,
		domainColumns, rlTable, tdTable, fkCol, rlFKCol)
	rows, err := r.pool.Query(ctx, query, ampID)
	if err != nil {
		return nil, fmt.Errorf("query %s: %w", rlTable, err)
	}
	defer rows.Close()
	return scanDomainRows(rows)
}

func (r *AMPRepo) getAMPIngredients(ctx context.Context, ampID int64) ([]entity.AMPIngredient, error) {
	const query = `SELECT rt.CO_SEQ_ID, rt.CO_ISID, rt.CO_APID, rt.QT_STRNTH, rt.CO_UOMCD,
		is.CO_SEQ_ID, is.NU_ISID, is.DT_ISIDDT, is.NO_NM, is.NO_NM_PT_BR,
		is.CO_ISRCID, is.CO_DCBCD, is.CO_NCAS, is.ST_REGISTRO_ATIVO, is.ST_CHANGE_APPROVE
		FROM rt_amp_ingredient_subst rt
		LEFT JOIN td_ingredient_substances is ON is.CO_SEQ_ID = rt.CO_ISID
		WHERE rt.CO_APID = $1 AND rt.ST_REGISTRO_ATIVO = 'ACTIVE'`
	rows, err := r.pool.Query(ctx, query, ampID)
	if err != nil {
		return nil, fmt.Errorf("query amp ingredients: %w", err)
	}
	defer rows.Close()
	var result []entity.AMPIngredient
	for rows.Next() {
		var ing entity.AMPIngredient
		var subst entity.IngredientSubstance
		var substID sql.NullInt64
		err := rows.Scan(
			&ing.COSeqID, &ing.COIsID, &ing.COApID, &ing.QTStrnth, &ing.COUomCd,
			&substID, &subst.NUIsID, &subst.DTIsIDDt, &subst.NONm, &subst.NONmPtBr,
			&subst.COIsrCd, &subst.CODcbCd, &subst.CONCas, &subst.STRegistroAtivo, &subst.STChangeApprove,
		)
		if err != nil {
			return nil, fmt.Errorf("scan amp ingredient: %w", err)
		}
		if substID.Valid {
			subst.COSeqID = substID.Int64
			ing.IngredientSubstance = &subst
		}
		result = append(result, ing)
	}
	return result, rows.Err()
}

func (r *AMPRepo) List(ctx context.Context, filters repository.FilterParams) (*entity.CursorPage[entity.AMP], error) {
	limit := filters.Limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	var where []string
	var args []interface{}
	argIdx := 1

	where = append(where, "a.ST_REGISTRO_ATIVO = 'ACTIVE'")

	if filters.Nome != "" {
		where = append(where, fmt.Sprintf("a.NO_NM ILIKE $%d", argIdx))
		args = append(args, "%"+filters.Nome+"%")
		argIdx++
	}
	if filters.Codigo != "" {
		where = append(where, fmt.Sprintf("a.NU_APID = $%d", argIdx))
		args = append(args, filters.Codigo)
		argIdx++
	}
	if filters.Fabricante != "" {
		where = append(where, fmt.Sprintf("s.NO_DESCR ILIKE $%d", argIdx))
		args = append(args, "%"+filters.Fabricante+"%")
		argIdx++
	}
	if filters.Ativo != nil {
		if *filters.Ativo {
			where = append(where, "a.ST_REGISTRO_ATIVO = 'ACTIVE'")
		} else {
			where = append(where, "a.ST_REGISTRO_ATIVO = 'INACTIVE'")
		}
	}

	offset := decodeCursor(filters.Cursor)

	joinClause := ""
	if filters.Fabricante != "" {
		joinClause = " JOIN td_supplier s ON s.CO_SEQ_ID = a.CO_SUPPCD"
	}

	whereClause := strings.Join(where, " AND ")
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM tb_amp a%s WHERE %s", joinClause, whereClause)
	var total int64
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, fmt.Errorf("count amp: %w", err)
	}

	dataQuery := fmt.Sprintf(
		"SELECT a.%s FROM tb_amp a%s WHERE %s ORDER BY a.CO_SEQ_ID LIMIT $%d OFFSET $%d",
		ampColumns, joinClause, whereClause, argIdx, argIdx+1,
	)
	args = append(args, limit+1, offset)

	rows, err := r.pool.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("query amp list: %w", err)
	}
	defer rows.Close()

	var items []entity.AMP
	for rows.Next() {
		var a entity.AMP
		if err := rows.Scan(
			&a.COSeqID, &a.NUApID, &a.COVpID, &a.STCombProdCd, &a.NONm, &a.DSDescr, &a.NOAbbrevNm,
			&a.COSuppCd, &a.COFlavourCd, &a.COLicAuthChangeCd, &a.STParallelImport, &a.COLicAuthCd,
			&a.COAvailRestrictCd, &a.COMedClsCd, &a.COMonitoringReasonCd, &a.COEnteralAdminStatusCd,
			&a.DSEnteralTubesAdminObs, &a.NUNReg, &a.NUNProc, &a.NUVencReg, &a.NUValidity, &a.COValidityUnit,
			&a.STRegistroAtivo, &a.STChangeApprove,
		); err != nil {
			return nil, fmt.Errorf("scan amp: %w", err)
		}
		items = append(items, a)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	var nextCursor string
	if len(items) > limit {
		items = items[:limit]
		nextCursor = encodeCursor(offset + int64(limit))
	}

	return &entity.CursorPage[entity.AMP]{
		Items:  items,
		Cursor: nextCursor,
		Limit:  limit,
		Total:  total,
	}, nil
}
