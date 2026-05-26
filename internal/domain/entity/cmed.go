package entity

import "encoding/json"

type CMEDConformidade struct {
	COSeqID                int64           `json:"co_seq_id" db:"co_seq_id"`
	NUSanReg               *int64          `json:"nu_sanreg" db:"nu_sanreg"`
	NUGgrem                *string         `json:"nu_ggrem" db:"nu_ggrem"`
	DSSubstancia           *string         `json:"ds_substancia" db:"ds_substancia"`
	NUCnpj                 *string         `json:"nu_cnpj" db:"nu_cnpj"`
	NOLaboratorio          *string         `json:"no_laboratorio" db:"no_laboratorio"`
	NOProduto              *string         `json:"no_produto" db:"no_produto"`
	DSApresentacao         *string         `json:"ds_apresentacao" db:"ds_apresentacao"`
	DSClasseTerapeutica    *string         `json:"ds_classe_terapeutica" db:"ds_classe_terapeutica"`
	TPProduto              *string         `json:"tp_produto" db:"tp_produto"`
	TPRegimePreco          *string         `json:"tp_regime_preco" db:"tp_regime_preco"`
	NUEAN1                 *string         `json:"nu_ean1" db:"nu_ean1"`
	NUEAN2                 *string         `json:"nu_ean2" db:"nu_ean2"`
	NUEAN3                 *string         `json:"nu_ean3" db:"nu_ean3"`
	VRPFSemImpostos        *float64        `json:"vr_pf_sem_impostos" db:"vr_pf_sem_impostos"`
	VRPF0                  *float64        `json:"vr_pf_0" db:"vr_pf_0"`
	VRPF12                 *float64        `json:"vr_pf_12" db:"vr_pf_12"`
	VRPF17                 *float64        `json:"vr_pf_17" db:"vr_pf_17"`
	VRPF18                 *float64        `json:"vr_pf_18" db:"vr_pf_18"`
	VRPF20                 *float64        `json:"vr_pf_20" db:"vr_pf_20"`
	VRPMCSemImpostos       *float64        `json:"vr_pmc_sem_impostos" db:"vr_pmc_sem_impostos"`
	VRPMC0                 *float64        `json:"vr_pmc_0" db:"vr_pmc_0"`
	VRPMC12                *float64        `json:"vr_pmc_12" db:"vr_pmc_12"`
	VRPMC17                *float64        `json:"vr_pmc_17" db:"vr_pmc_17"`
	VRPMC18                *float64        `json:"vr_pmc_18" db:"vr_pmc_18"`
	VRPMC20                *float64        `json:"vr_pmc_20" db:"vr_pmc_20"`
	JSPrecosPF             json.RawMessage `json:"js_precos_pf" db:"js_precos_pf"`
	JSPrecosPMC            json.RawMessage `json:"js_precos_pmc" db:"js_precos_pmc"`
	STRestricaoHospitalar  *string         `json:"st_restricao_hospitalar" db:"st_restricao_hospitalar"`
	STCap                  *string         `json:"st_cap" db:"st_cap"`
	STConfaz87             *string         `json:"st_confaz_87" db:"st_confaz_87"`
	STIcms0                *string         `json:"st_icms_0" db:"st_icms_0"`
	DSAnaliseRecural       *string         `json:"ds_analise_recural" db:"ds_analise_recural"`
	DSListaPisCofins       *string         `json:"ds_lista_pis_cofins" db:"ds_lista_pis_cofins"`
	STComercializacao      *string         `json:"st_comercializacao" db:"st_comercializacao"`
	DSTarja                *string         `json:"ds_tarja" db:"ds_tarja"`
	DSDestinacaoComercial  *string         `json:"ds_destinacao_comercial" db:"ds_destinacao_comercial"`
	DTReferencia           string          `json:"dt_referencia" db:"dt_referencia"`
	STRegistroAtivo        string          `json:"st_registro_ativo" db:"st_registro_ativo"`
}
