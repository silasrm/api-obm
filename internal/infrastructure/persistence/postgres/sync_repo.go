package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SyncRepo struct {
	pool *pgxpool.Pool
}

func NewSyncRepo(pool *pgxpool.Pool) *SyncRepo {
	return &SyncRepo{pool: pool}
}

func (r *SyncRepo) GetAllVMPs(ctx context.Context) ([]map[string]interface{}, error) {
	const query = `
		SELECT v.CO_SEQ_ID, v.NU_VPID, v.NO_NM, v.CO_VTMID,
			COALESCE(vt.NO_NM, '') AS vtm_name,
			COALESCE(bn.NO_DESCR, '') AS basis_of_name,
			COALESCE(ps.NO_DESCR, '') AS pres_status,
			COALESCE(df.NO_DESCR, '') AS df_indicator
		FROM tb_vmp v
		LEFT JOIN tb_vtm vt ON vt.CO_SEQ_ID = v.CO_VTMID AND vt.ST_REGISTRO_ATIVO = 'ACTIVE'
		LEFT JOIN td_basis_of_name bn ON bn.CO_SEQ_ID = v.CO_BASISCD AND bn.ST_REGISTRO_ATIVO = 'ACTIVE'
		LEFT JOIN td_virtual_product_pres_status ps ON ps.CO_SEQ_ID = v.CO_PRES_STATCD AND ps.ST_REGISTRO_ATIVO = 'ACTIVE'
		LEFT JOIN td_df_indicator df ON df.CO_SEQ_ID = v.CO_DF_INDCD AND df.ST_REGISTRO_ATIVO = 'ACTIVE'
		WHERE v.ST_REGISTRO_ATIVO = 'ACTIVE'
		ORDER BY v.CO_SEQ_ID`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query all vmps: %w", err)
	}
	defer rows.Close()
	var result []map[string]interface{}
	for rows.Next() {
		var coSeqID, coVtmid int64
		var nuVpid, noNm, vtmName, basisOfName, presStatus, dfIndicator string
		if err := rows.Scan(&coSeqID, &nuVpid, &noNm, &coVtmid, &vtmName, &basisOfName, &presStatus, &dfIndicator); err != nil {
			return nil, fmt.Errorf("scan vmp: %w", err)
		}
		result = append(result, map[string]interface{}{
			"co_seq_id":     coSeqID,
			"nu_vpid":       nuVpid,
			"no_nm":         noNm,
			"co_vtmid":      coVtmid,
			"vtm_name":      vtmName,
			"basis_of_name": basisOfName,
			"pres_status":   presStatus,
			"df_indicator":  dfIndicator,
		})
	}
	return result, rows.Err()
}

func (r *SyncRepo) GetAllAMPs(ctx context.Context) ([]map[string]interface{}, error) {
	const query = `
		SELECT a.CO_SEQ_ID, a.NU_APID, a.NO_NM, a.DS_DESCR, a.CO_VPID,
			COALESCE(v.NO_NM, '') AS vmp_name,
			COALESCE(s.NO_DESCR, '') AS supplier_name
		FROM tb_amp a
		LEFT JOIN tb_vmp v ON v.CO_SEQ_ID = a.CO_VPID AND v.ST_REGISTRO_ATIVO = 'ACTIVE'
		LEFT JOIN td_supplier s ON s.CO_SEQ_ID = a.CO_SUPPCD AND s.ST_REGISTRO_ATIVO = 'ACTIVE'
		WHERE a.ST_REGISTRO_ATIVO = 'ACTIVE'
		ORDER BY a.CO_SEQ_ID`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query all amps: %w", err)
	}
	defer rows.Close()
	var result []map[string]interface{}
	for rows.Next() {
		var coSeqID, coVpid int64
		var nuApid, noNm string
		var dsDescr, vmpName, supplierName *string
		if err := rows.Scan(&coSeqID, &nuApid, &noNm, &dsDescr, &coVpid, &vmpName, &supplierName); err != nil {
			return nil, fmt.Errorf("scan amp: %w", err)
		}
		doc := map[string]interface{}{
			"co_seq_id": coSeqID,
			"nu_apid":   nuApid,
			"no_nm":     noNm,
			"co_vpid":   coVpid,
		}
		if dsDescr != nil {
			doc["ds_descr"] = *dsDescr
		}
		if vmpName != nil {
			doc["vmp_name"] = *vmpName
		}
		if supplierName != nil {
			doc["supplier_name"] = *supplierName
		}
		result = append(result, doc)
	}
	return result, rows.Err()
}

func (r *SyncRepo) GetAllSuppliers(ctx context.Context) ([]map[string]interface{}, error) {
	const query = `
		SELECT CO_SEQ_ID, NU_CD, NO_DESCR, NU_CNPJ
		FROM td_supplier
		WHERE ST_REGISTRO_ATIVO = 'ACTIVE'
		ORDER BY CO_SEQ_ID`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query all suppliers: %w", err)
	}
	defer rows.Close()
	var result []map[string]interface{}
	for rows.Next() {
		var coSeqID int64
		var nuCd, noDescr string
		var nuCnpj *string
		if err := rows.Scan(&coSeqID, &nuCd, &noDescr, &nuCnpj); err != nil {
			return nil, fmt.Errorf("scan supplier: %w", err)
		}
		doc := map[string]interface{}{
			"co_seq_id": coSeqID,
			"nu_cd":     nuCd,
			"no_descr":  noDescr,
		}
		if nuCnpj != nil {
			doc["nu_cnpj"] = *nuCnpj
		}
		result = append(result, doc)
	}
	return result, rows.Err()
}
