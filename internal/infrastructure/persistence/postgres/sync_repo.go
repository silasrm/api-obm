package postgres

import (
	"context"
	"fmt"
	"time"

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

func (r *SyncRepo) GetAllCMED(ctx context.Context) ([]map[string]interface{}, error) {
	const query = `
		SELECT c.CO_SEQ_ID, c.NU_SANREG, c.NO_PRODUTO, c.DS_SUBSTANCIA,
		       c.NO_LABORATORIO, c.NU_EAN1, c.NU_EAN2, c.NU_EAN3,
		       c.DS_CLASSE_TERAPEUTICA, c.DS_APRESENTACAO,
		       c.TP_PRODUTO, c.TP_REGIME_PRECO, c.DS_TARJA,
		       c.DT_REFERENCIA, c.VR_PF_SEM_IMPOSTOS, c.VR_PMC_SEM_IMPOSTOS,
		       c.NU_CNPJ
		FROM tb_cmed_conformidade c
		WHERE c.ST_REGISTRO_ATIVO = 'ACTIVE'
		ORDER BY c.CO_SEQ_ID`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query all cmed: %w", err)
	}
	defer rows.Close()
	var result []map[string]interface{}
	for rows.Next() {
		var coSeqID int64
		var nuSanreg *int64
		var noProduto, dsSubstancia, noLaboratorio *string
		var nuEan1, nuEan2, nuEan3, dsClasseTerapeutica, dsApresentacao *string
		var tpProduto, tpRegimePreco, dsTarja *string
		var dtReferencia *time.Time
		var vrPFSemImpostos, vrPMCSemImpostos *float64
		var nuCnpj *string
		if err := rows.Scan(&coSeqID, &nuSanreg, &noProduto, &dsSubstancia, &noLaboratorio,
			&nuEan1, &nuEan2, &nuEan3, &dsClasseTerapeutica, &dsApresentacao,
			&tpProduto, &tpRegimePreco, &dsTarja, &dtReferencia,
			&vrPFSemImpostos, &vrPMCSemImpostos, &nuCnpj); err != nil {
			return nil, fmt.Errorf("scan cmed: %w", err)
		}
		doc := map[string]interface{}{
			"co_seq_id": coSeqID,
		}
		if nuSanreg != nil {
			doc["nu_sanreg"] = fmt.Sprintf("%d", *nuSanreg)
		}
		if noProduto != nil {
			doc["no_produto"] = *noProduto
		}
		if dsSubstancia != nil {
			doc["ds_substancia"] = *dsSubstancia
		}
		if noLaboratorio != nil {
			doc["no_laboratorio"] = *noLaboratorio
		}
		if nuEan1 != nil {
			doc["nu_ean1"] = *nuEan1
		}
		if nuEan2 != nil {
			doc["nu_ean2"] = *nuEan2
		}
		if nuEan3 != nil {
			doc["nu_ean3"] = *nuEan3
		}
		if dsClasseTerapeutica != nil {
			doc["ds_classe_terapeutica"] = *dsClasseTerapeutica
		}
		if dsApresentacao != nil {
			doc["ds_apresentacao"] = *dsApresentacao
		}
		if tpProduto != nil {
			doc["tp_produto"] = *tpProduto
		}
		if tpRegimePreco != nil {
			doc["tp_regime_preco"] = *tpRegimePreco
		}
		if dsTarja != nil {
			doc["ds_tarja"] = *dsTarja
		}
	if dtReferencia != nil {
		doc["dt_referencia"] = dtReferencia.Format("2006-01-02")
	}
		if vrPFSemImpostos != nil {
			doc["vr_pf_sem_impostos"] = *vrPFSemImpostos
		}
		if vrPMCSemImpostos != nil {
			doc["vr_pmc_sem_impostos"] = *vrPMCSemImpostos
		}
		if nuCnpj != nil {
			doc["nu_cnpj"] = *nuCnpj
		}
		result = append(result, doc)
	}
	return result, rows.Err()
}
