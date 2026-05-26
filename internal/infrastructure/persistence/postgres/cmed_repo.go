package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/silasrm/api-obm/internal/domain/entity"
	"github.com/silasrm/api-obm/internal/domain/repository"
)

type CMEDRepo struct {
	pool *pgxpool.Pool
}

func NewCMEDRepo(pool *pgxpool.Pool) *CMEDRepo {
	return &CMEDRepo{pool: pool}
}

var _ repository.CMEDRepository = (*CMEDRepo)(nil)

const cmedColumns = `co_seq_id, nu_sanreg, nu_ggrem, ds_substancia, nu_cnpj, no_laboratorio,
no_produto, ds_apresentacao, ds_classe_terapeutica, tp_produto, tp_regime_preco,
nu_ean1, nu_ean2, nu_ean3,
vr_pf_sem_impostos, vr_pf_0, vr_pf_12, vr_pf_17, vr_pf_18, vr_pf_20,
vr_pmc_sem_impostos, vr_pmc_0, vr_pmc_12, vr_pmc_17, vr_pmc_18, vr_pmc_20,
js_precos_pf, js_precos_pmc,
st_restricao_hospitalar, st_cap, st_confaz_87, st_icms_0, ds_analise_recural,
ds_lista_pis_cofins, st_comercializacao, ds_tarja, ds_destinacao_comercial,
dt_referencia::text, st_registro_ativo`

func scanCMED(row pgx.Row) (*entity.CMEDConformidade, error) {
	var c entity.CMEDConformidade
	err := row.Scan(
		&c.COSeqID, &c.NUSanReg, &c.NUGgrem, &c.DSSubstancia, &c.NUCnpj, &c.NOLaboratorio,
		&c.NOProduto, &c.DSApresentacao, &c.DSClasseTerapeutica, &c.TPProduto, &c.TPRegimePreco,
		&c.NUEAN1, &c.NUEAN2, &c.NUEAN3,
		&c.VRPFSemImpostos, &c.VRPF0, &c.VRPF12, &c.VRPF17, &c.VRPF18, &c.VRPF20,
		&c.VRPMCSemImpostos, &c.VRPMC0, &c.VRPMC12, &c.VRPMC17, &c.VRPMC18, &c.VRPMC20,
		&c.JSPrecosPF, &c.JSPrecosPMC,
		&c.STRestricaoHospitalar, &c.STCap, &c.STConfaz87, &c.STIcms0, &c.DSAnaliseRecural,
		&c.DSListaPisCofins, &c.STComercializacao, &c.DSTarja, &c.DSDestinacaoComercial,
		&c.DTReferencia, &c.STRegistroAtivo,
	)
	return &c, err
}

func scanCMEDRows(rows pgx.Rows) ([]entity.CMEDConformidade, error) {
	var items []entity.CMEDConformidade
	for rows.Next() {
		var c entity.CMEDConformidade
		if err := rows.Scan(
			&c.COSeqID, &c.NUSanReg, &c.NUGgrem, &c.DSSubstancia, &c.NUCnpj, &c.NOLaboratorio,
			&c.NOProduto, &c.DSApresentacao, &c.DSClasseTerapeutica, &c.TPProduto, &c.TPRegimePreco,
			&c.NUEAN1, &c.NUEAN2, &c.NUEAN3,
			&c.VRPFSemImpostos, &c.VRPF0, &c.VRPF12, &c.VRPF17, &c.VRPF18, &c.VRPF20,
			&c.VRPMCSemImpostos, &c.VRPMC0, &c.VRPMC12, &c.VRPMC17, &c.VRPMC18, &c.VRPMC20,
			&c.JSPrecosPF, &c.JSPrecosPMC,
			&c.STRestricaoHospitalar, &c.STCap, &c.STConfaz87, &c.STIcms0, &c.DSAnaliseRecural,
			&c.DSListaPisCofins, &c.STComercializacao, &c.DSTarja, &c.DSDestinacaoComercial,
			&c.DTReferencia, &c.STRegistroAtivo,
		); err != nil {
			return nil, err
		}
		items = append(items, c)
	}
	return items, rows.Err()
}

func (r *CMEDRepo) GetByID(ctx context.Context, id int64) (*entity.CMEDConformidade, error) {
	query := fmt.Sprintf("SELECT %s FROM tb_cmed_conformidade WHERE co_seq_id = $1", cmedColumns)
	c, err := scanCMED(r.pool.QueryRow(ctx, query, id))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, pgx.ErrNoRows
		}
		return nil, fmt.Errorf("query cmed: %w", err)
	}
	return c, nil
}

func (r *CMEDRepo) GetByNuSanReg(ctx context.Context, nuSanReg int64, dtReferencia string) (*entity.CMEDConformidade, error) {
	var query string
	var args []interface{}
	if dtReferencia != "" {
		query = fmt.Sprintf("SELECT %s FROM tb_cmed_conformidade WHERE nu_sanreg = $1 AND dt_referencia = $2 AND st_registro_ativo = 'ACTIVE'", cmedColumns)
		args = []interface{}{nuSanReg, dtReferencia}
	} else {
		query = fmt.Sprintf("SELECT %s FROM tb_cmed_conformidade WHERE nu_sanreg = $1 AND st_registro_ativo = 'ACTIVE' ORDER BY dt_referencia DESC LIMIT 1", cmedColumns)
		args = []interface{}{nuSanReg}
	}
	c, err := scanCMED(r.pool.QueryRow(ctx, query, args...))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query cmed by nu_sanreg: %w", err)
	}
	return c, nil
}

func (r *CMEDRepo) GetByEAN(ctx context.Context, ean string, dtReferencia string) (*entity.CMEDConformidade, error) {
	var query string
	var args []interface{}
	if dtReferencia != "" {
		query = fmt.Sprintf("SELECT %s FROM tb_cmed_conformidade WHERE (nu_ean1 = $1 OR nu_ean2 = $1 OR nu_ean3 = $1) AND dt_referencia = $2 AND st_registro_ativo = 'ACTIVE'", cmedColumns)
		args = []interface{}{ean, dtReferencia}
	} else {
		query = fmt.Sprintf("SELECT %s FROM tb_cmed_conformidade WHERE (nu_ean1 = $1 OR nu_ean2 = $1 OR nu_ean3 = $1) AND st_registro_ativo = 'ACTIVE' ORDER BY dt_referencia DESC LIMIT 1", cmedColumns)
		args = []interface{}{ean}
	}
	c, err := scanCMED(r.pool.QueryRow(ctx, query, args...))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query cmed by ean: %w", err)
	}
	return c, nil
}

func (r *CMEDRepo) List(ctx context.Context, filter repository.CMEDFilterParams) (*entity.CursorPage[entity.CMEDConformidade], error) {
	limit := filter.Limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	var where []string
	var args []interface{}
	argIdx := 1

	if filter.Ativo != nil {
		if *filter.Ativo {
			where = append(where, "st_registro_ativo = 'ACTIVE'")
		} else {
			where = append(where, "st_registro_ativo = 'INACTIVE'")
		}
	} else {
		where = append(where, "st_registro_ativo = 'ACTIVE'")
	}

	if filter.Nome != "" {
		where = append(where, fmt.Sprintf("(no_produto ILIKE '%%' || $%d || '%%' OR ds_substancia ILIKE '%%' || $%d || '%%')", argIdx, argIdx))
		args = append(args, filter.Nome)
		argIdx++
	}

	if filter.Registro != "" {
		regInt, err := strconv.ParseInt(filter.Registro, 10, 64)
		if err == nil {
			where = append(where, fmt.Sprintf("nu_sanreg = $%d", argIdx))
			args = append(args, regInt)
			argIdx++
		}
	}

	if filter.EAN != "" {
		where = append(where, fmt.Sprintf("(nu_ean1 = $%d OR nu_ean2 = $%d OR nu_ean3 = $%d)", argIdx, argIdx, argIdx))
		args = append(args, filter.EAN)
		argIdx++
	}

	if filter.Tarja != "" {
		where = append(where, fmt.Sprintf("ds_tarja ILIKE '%%' || $%d || '%%'", argIdx))
		args = append(args, filter.Tarja)
		argIdx++
	}

	if filter.TipoProduto != "" {
		where = append(where, fmt.Sprintf("tp_produto = $%d", argIdx))
		args = append(args, filter.TipoProduto)
		argIdx++
	}

	if filter.RegimePreco != "" {
		where = append(where, fmt.Sprintf("tp_regime_preco = $%d", argIdx))
		args = append(args, filter.RegimePreco)
		argIdx++
	}

	if filter.DTReferencia != "" {
		where = append(where, fmt.Sprintf("dt_referencia = $%d::date", argIdx))
		args = append(args, filter.DTReferencia)
		argIdx++
	}

	offset := decodeCursor(filter.Cursor)
	whereClause := strings.Join(where, " AND ")

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM tb_cmed_conformidade WHERE %s", whereClause)
	var total int64
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, fmt.Errorf("count cmed: %w", err)
	}

	dataQuery := fmt.Sprintf(
		"SELECT %s FROM tb_cmed_conformidade WHERE %s ORDER BY co_seq_id LIMIT $%d OFFSET $%d",
		cmedColumns, whereClause, argIdx, argIdx+1,
	)
	args = append(args, limit+1, offset)

	rows, err := r.pool.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("query cmed list: %w", err)
	}
	defer rows.Close()

	items, err := scanCMEDRows(rows)
	if err != nil {
		return nil, err
	}

	var nextCursor string
	if len(items) > limit {
		items = items[:limit]
		nextCursor = encodeCursor(offset + int64(limit))
	}

	return &entity.CursorPage[entity.CMEDConformidade]{
		Items:  items,
		Cursor: nextCursor,
		Limit:  limit,
		Total:  total,
	}, nil
}

func (r *CMEDRepo) GetHistorico(ctx context.Context, nuSanReg int64) ([]entity.CMEDConformidade, error) {
	query := fmt.Sprintf("SELECT %s FROM tb_cmed_conformidade WHERE nu_sanreg = $1 ORDER BY dt_referencia DESC", cmedColumns)
	rows, err := r.pool.Query(ctx, query, nuSanReg)
	if err != nil {
		return nil, fmt.Errorf("query cmed historico: %w", err)
	}
	defer rows.Close()
	return scanCMEDRows(rows)
}

func (r *CMEDRepo) UpsertBatch(ctx context.Context, records []entity.CMEDConformidade) (int64, error) {
	const batchSize = 500
	var count int64

	for i := 0; i < len(records); i += batchSize {
		end := i + batchSize
		if end > len(records) {
			end = len(records)
		}
		batch := records[i:end]

		for _, rec := range batch {
			jsPrecosPF := rec.JSPrecosPF
			if jsPrecosPF == nil {
				jsPrecosPF = json.RawMessage("null")
			}
			jsPrecosPMC := rec.JSPrecosPMC
			if jsPrecosPMC == nil {
				jsPrecosPMC = json.RawMessage("null")
			}

			const query = `INSERT INTO tb_cmed_conformidade (
				nu_sanreg, nu_ggrem, ds_substancia, nu_cnpj, no_laboratorio,
				no_produto, ds_apresentacao, ds_classe_terapeutica, tp_produto, tp_regime_preco,
				nu_ean1, nu_ean2, nu_ean3,
				vr_pf_sem_impostos, vr_pf_0, vr_pf_12, vr_pf_17, vr_pf_18, vr_pf_20,
				vr_pmc_sem_impostos, vr_pmc_0, vr_pmc_12, vr_pmc_17, vr_pmc_18, vr_pmc_20,
				js_precos_pf, js_precos_pmc,
				st_restricao_hospitalar, st_cap, st_confaz_87, st_icms_0, ds_analise_recural,
				ds_lista_pis_cofins, st_comercializacao, ds_tarja, ds_destinacao_comercial,
				dt_referencia, st_registro_ativo
			) VALUES (
				$1, $2, $3, $4, $5,
				$6, $7, $8, $9, $10,
				$11, $12, $13,
				$14, $15, $16, $17, $18, $19,
				$20, $21, $22, $23, $24, $25,
				$26, $27,
				$28, $29, $30, $31, $32,
				$33, $34, $35, $36,
				$37, $38
			) ON CONFLICT (nu_sanreg, dt_referencia) DO UPDATE SET
				nu_ggrem = EXCLUDED.nu_ggrem,
				ds_substancia = EXCLUDED.ds_substancia,
				nu_cnpj = EXCLUDED.nu_cnpj,
				no_laboratorio = EXCLUDED.no_laboratorio,
				no_produto = EXCLUDED.no_produto,
				ds_apresentacao = EXCLUDED.ds_apresentacao,
				ds_classe_terapeutica = EXCLUDED.ds_classe_terapeutica,
				tp_produto = EXCLUDED.tp_produto,
				tp_regime_preco = EXCLUDED.tp_regime_preco,
				nu_ean1 = EXCLUDED.nu_ean1,
				nu_ean2 = EXCLUDED.nu_ean2,
				nu_ean3 = EXCLUDED.nu_ean3,
				vr_pf_sem_impostos = EXCLUDED.vr_pf_sem_impostos,
				vr_pf_0 = EXCLUDED.vr_pf_0,
				vr_pf_12 = EXCLUDED.vr_pf_12,
				vr_pf_17 = EXCLUDED.vr_pf_17,
				vr_pf_18 = EXCLUDED.vr_pf_18,
				vr_pf_20 = EXCLUDED.vr_pf_20,
				vr_pmc_sem_impostos = EXCLUDED.vr_pmc_sem_impostos,
				vr_pmc_0 = EXCLUDED.vr_pmc_0,
				vr_pmc_12 = EXCLUDED.vr_pmc_12,
				vr_pmc_17 = EXCLUDED.vr_pmc_17,
				vr_pmc_18 = EXCLUDED.vr_pmc_18,
				vr_pmc_20 = EXCLUDED.vr_pmc_20,
				js_precos_pf = EXCLUDED.js_precos_pf,
				js_precos_pmc = EXCLUDED.js_precos_pmc,
				st_restricao_hospitalar = EXCLUDED.st_restricao_hospitalar,
				st_cap = EXCLUDED.st_cap,
				st_confaz_87 = EXCLUDED.st_confaz_87,
				st_icms_0 = EXCLUDED.st_icms_0,
				ds_analise_recural = EXCLUDED.ds_analise_recural,
				ds_lista_pis_cofins = EXCLUDED.ds_lista_pis_cofins,
				st_comercializacao = EXCLUDED.st_comercializacao,
				ds_tarja = EXCLUDED.ds_tarja,
				ds_destinacao_comercial = EXCLUDED.ds_destinacao_comercial,
				st_registro_ativo = EXCLUDED.st_registro_ativo`

			tag, err := r.pool.Exec(ctx, query,
				rec.NUSanReg, rec.NUGgrem, rec.DSSubstancia, rec.NUCnpj, rec.NOLaboratorio,
				rec.NOProduto, rec.DSApresentacao, rec.DSClasseTerapeutica, rec.TPProduto, rec.TPRegimePreco,
				rec.NUEAN1, rec.NUEAN2, rec.NUEAN3,
				rec.VRPFSemImpostos, rec.VRPF0, rec.VRPF12, rec.VRPF17, rec.VRPF18, rec.VRPF20,
				rec.VRPMCSemImpostos, rec.VRPMC0, rec.VRPMC12, rec.VRPMC17, rec.VRPMC18, rec.VRPMC20,
				jsPrecosPF, jsPrecosPMC,
				rec.STRestricaoHospitalar, rec.STCap, rec.STConfaz87, rec.STIcms0, rec.DSAnaliseRecural,
				rec.DSListaPisCofins, rec.STComercializacao, rec.DSTarja, rec.DSDestinacaoComercial,
				rec.DTReferencia, rec.STRegistroAtivo,
			)
			if err != nil {
				return count, fmt.Errorf("upsert cmed: %w", err)
			}
			count += tag.RowsAffected()
		}
	}

	return count, nil
}
