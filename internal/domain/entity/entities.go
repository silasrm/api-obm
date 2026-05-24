package entity

import (
	"database/sql"
	"time"
)

type VTM struct {
	COSeqID              int64     `json:"co_seq_id" db:"co_seq_id"`
	NUVtmID              int64     `json:"nu_vtmid" db:"nu_vtmid"`
	NONm                 string    `json:"no_nm" db:"no_nm"`
	NOAbbrevNm           string    `json:"no_abbrevnm" db:"no_abbrevnm"`
	AUAuditCreateDatetime time.Time `json:"au_audit_create_datetime" db:"au_audit_create_datetime"`
	AUAuditChangeDatetime time.Time `json:"au_audit_change_datetime" db:"au_audit_change_datetime"`
	AUAuditCreateUsername string    `json:"au_audit_create_username" db:"au_audit_create_username"`
	AUAuditChangeUsername string    `json:"au_audit_change_username" db:"au_audit_change_username"`
	STRegistroAtivo      string    `json:"st_registro_ativo" db:"st_registro_ativo"`
	STChangeApprove      string    `json:"st_change_approve" db:"st_change_approve"`
}

type VMP struct {
	COSeqID         int64          `json:"co_seq_id" db:"co_seq_id"`
	NUVpID          string         `json:"nu_vpid" db:"nu_vpid"`
	DTVpIDDt        int64          `json:"dt_vpiddt" db:"dt_vpiddt"`
	COVpIDPrev      sql.NullInt64  `json:"co_vpidprev" db:"co_vpidprev"`
	STInvalid       sql.NullBool   `json:"st_invalid" db:"st_invalid"`
	COVtmID         int64          `json:"co_vtmid" db:"co_vtmid"`
	STCombProdCd    sql.NullBool   `json:"st_combprodcdd" db:"st_combprodcdd"`
	NONm            string         `json:"no_nm" db:"no_nm"`
	NUSnomed        sql.NullString `json:"nu_snomed" db:"nu_snomed"`
	COBasisCd       sql.NullInt64  `json:"co_basiscd" db:"co_basiscd"`
	DTNmDt          sql.NullInt64  `json:"dt_nmdt" db:"dt_nmdt"`
	NONmPrev        sql.NullString `json:"no_nmprev" db:"no_nmprev"`
	COBasisPrevCd   sql.NullInt64  `json:"co_basis_prevcd" db:"co_basis_prevcd"`
	CONmChangeCd    sql.NullInt64  `json:"co_nmchangecd" db:"co_nmchangecd"`
	NOAbbrevNm      sql.NullString `json:"no_abbrevnm" db:"no_abbrevnm"`
	STSugF          sql.NullBool   `json:"st_sug_f" db:"st_sug_f"`
	STGluF          sql.NullBool   `json:"st_glu_f" db:"st_glu_f"`
	STPresF         sql.NullBool   `json:"st_pres_f" db:"st_pres_f"`
	STCfcFree       sql.NullBool   `json:"st_cfc_free" db:"st_cfc_free"`
	COPresStatCd    sql.NullInt64  `json:"co_pres_statcd" db:"co_pres_statcd"`
	CONonAvailCd    sql.NullInt64  `json:"co_non_availcd" db:"co_non_availcd"`
	DTNonAvailDt    sql.NullInt64  `json:"dt_non_availdt" db:"dt_non_availdt"`
	CODFIndCd       sql.NullInt64  `json:"co_df_indcd" db:"co_df_indcd"`
	QTUdfs          sql.NullFloat64 `json:"qt_udfs" db:"qt_udfs"`
	COUdfsUomCd     sql.NullInt64  `json:"co_udfs_uomcd" db:"co_udfs_uomcd"`
	COUnitDoseUomCd sql.NullInt64  `json:"co_unit_dose_uomcd" db:"co_unit_dose_uomcd"`
	STFtn           sql.NullBool   `json:"st_ftn" db:"st_ftn"`
	COBrimunoCd     sql.NullInt64  `json:"co_brimunocd" db:"co_brimunocd"`
	STHorus         sql.NullBool   `json:"st_horus" db:"st_horus"`
	DSNotification  sql.NullString `json:"ds_notification" db:"ds_notification"`
	COTargetCd      sql.NullInt64  `json:"co_targetcd" db:"co_targetcd"`
	COAnvsClsCd     sql.NullInt64  `json:"co_anvsclscd" db:"co_anvsclscd"`
	NUVersion       sql.NullInt64  `json:"nu_version" db:"nu_version"`
	STRegistroAtivo string         `json:"st_registro_ativo" db:"st_registro_ativo"`
	STChangeApprove string         `json:"st_change_approve" db:"st_change_approve"`
}

type VMPP struct {
	COSeqID                  int64          `json:"co_seq_id" db:"co_seq_id"`
	NUVppID                  string         `json:"nu_vppid" db:"nu_vppid"`
	COVpID                   int64          `json:"co_vpid" db:"co_vpid"`
	NONm                     string         `json:"no_nm" db:"no_nm"`
	STCombPackCd             sql.NullBool   `json:"st_combpackcd" db:"st_combpackcd"`
	QTQtyvalFornecimento     sql.NullFloat64 `json:"qt_qtyval_fornecimento" db:"qt_qtyval_fornecimento"`
	COQtyUomCdFornecimento   sql.NullInt64  `json:"co_qty_uomcd_fornecimento" db:"co_qty_uomcd_fornecimento"`
	QTQtyval                 sql.NullFloat64 `json:"qt_qtyval" db:"qt_qtyval"`
	COQtyUomCd               sql.NullInt64  `json:"co_qty_uomcd" db:"co_qty_uomcd"`
	STRegistroAtivo          string         `json:"st_registro_ativo" db:"st_registro_ativo"`
	STChangeApprove          string         `json:"st_change_approve" db:"st_change_approve"`
}

type AMP struct {
	COSeqID                 int64          `json:"co_seq_id" db:"co_seq_id"`
	NUApID                  string         `json:"nu_apid" db:"nu_apid"`
	COVpID                  int64          `json:"co_vpid" db:"co_vpid"`
	STCombProdCd            sql.NullBool   `json:"st_combprodcdd" db:"st_combprodcdd"`
	NONm                    string         `json:"no_nm" db:"no_nm"`
	DSDescr                 sql.NullString `json:"ds_descr" db:"ds_descr"`
	NOAbbrevNm              sql.NullString `json:"no_abbrevnm" db:"no_abbrevnm"`
	COSuppCd                int64          `json:"co_suppcd" db:"co_suppcd"`
	COFlavourCd             sql.NullInt64  `json:"co_flavourcd" db:"co_flavourcd"`
	COLicAuthChangeCd       sql.NullInt64  `json:"co_lic_authchangecd" db:"co_lic_authchangecd"`
	STParallelImport        sql.NullBool   `json:"st_parallel_import" db:"st_parallel_import"`
	COLicAuthCd             sql.NullInt64  `json:"co_lic_authcd" db:"co_lic_authcd"`
	COAvailRestrictCd       sql.NullInt64  `json:"co_avail_restrictcd" db:"co_avail_restrictcd"`
	COMedClsCd              sql.NullInt64  `json:"co_medclscd" db:"co_medclscd"`
	COMonitoringReasonCd    sql.NullInt64  `json:"co_monitoringreasoncd" db:"co_monitoringreasoncd"`
	COEnteralAdminStatusCd  sql.NullInt64  `json:"co_enteraladminstatuscd" db:"co_enteraladminstatuscd"`
	DSEnteralTubesAdminObs  sql.NullString `json:"ds_enteraltubesadminobs" db:"ds_enteraltubesadminobs"`
	NUNReg                  sql.NullInt64  `json:"nu_nreg" db:"nu_nreg"`
	NUNProc                 sql.NullString `json:"nu_nproc" db:"nu_nproc"`
	NUVencReg               sql.NullString `json:"nu_vencreg" db:"nu_vencreg"`
	NUValidity              sql.NullInt64  `json:"nu_validity" db:"nu_validity"`
	COValidityUnit          sql.NullString `json:"co_validityunit" db:"co_validityunit"`
	STRegistroAtivo         string         `json:"st_registro_ativo" db:"st_registro_ativo"`
	STChangeApprove         string         `json:"st_change_approve" db:"st_change_approve"`
}

type AMPP struct {
	COSeqID         int64          `json:"co_seq_id" db:"co_seq_id"`
	NUAppID         string         `json:"nu_appid" db:"nu_appid"`
	COVppID         int64          `json:"co_vppid" db:"co_vppid"`
	NONm            string         `json:"no_nm" db:"no_nm"`
	COApID          int64          `json:"co_apid" db:"co_apid"`
	DSSubP          sql.NullString `json:"ds_subp" db:"ds_subp"`
	STCombPackCd    sql.NullBool   `json:"st_combpackcd" db:"st_combpackcd"`
	COLegalCatCd    sql.NullInt64  `json:"co_legal_catcd" db:"co_legal_catcd"`
	CODiscCd        sql.NullInt64  `json:"co_disccd" db:"co_disccd"`
	DTDiscDt        sql.NullInt64  `json:"dt_discdt" db:"dt_discdt"`
	NUSanReg        sql.NullInt64  `json:"nu_sanreg" db:"nu_sanreg"`
	DTRegPublic     sql.NullInt64  `json:"dt_regpublic" db:"dt_regpublic"`
	STHosp          sql.NullBool   `json:"st_hosp" db:"st_hosp"`
	STRegistroAtivo string         `json:"st_registro_ativo" db:"st_registro_ativo"`
	COPrimaryPck    sql.NullInt64  `json:"co_primarypck" db:"co_primarypck"`
	COSecondaryPck  sql.NullInt64  `json:"co_secondarypck" db:"co_secondarypck"`
	COIndFarmPopCd  sql.NullInt64  `json:"co_indfarmpopcd" db:"co_indfarmpopcd"`
	NUEAN13a        sql.NullString `json:"nu_ean13a" db:"nu_ean13a"`
	NUEAN13b        sql.NullString `json:"nu_ean13b" db:"nu_ean13b"`
	NUEAN13c        sql.NullString `json:"nu_ean13c" db:"nu_ean13c"`
	STChangeApprove string         `json:"st_change_approve" db:"st_change_approve"`
}

type DCB struct {
	COSeqID         int64          `json:"co_seq_id" db:"co_seq_id"`
	NUDcb           int64          `json:"nu_dcb" db:"nu_dcb"`
	DSDcb           string         `json:"ds_dcb" db:"ds_dcb"`
	COCas           sql.NullString `json:"co_cas" db:"co_cas"`
	DSClassif       sql.NullString `json:"ds_classif" db:"ds_classif"`
	STRegistroAtivo string         `json:"st_registro_ativo" db:"st_registro_ativo"`
	STChangeApprove string         `json:"st_change_approve" db:"st_change_approve"`
}

type Supplier struct {
	COSeqID         int64          `json:"co_seq_id" db:"co_seq_id"`
	NUCd            string         `json:"nu_cd" db:"nu_cd"`
	DTCdDt          sql.NullInt64  `json:"dt_cddt" db:"dt_cddt"`
	NODescr         string         `json:"no_descr" db:"no_descr"`
	NUCnpj          sql.NullString `json:"nu_cnpj" db:"nu_cnpj"`
	NUNAuth         sql.NullInt64  `json:"nu_nauth" db:"nu_nauth"`
	COCountryCd     sql.NullInt64  `json:"co_countrycd" db:"co_countrycd"`
	STRegistroAtivo string         `json:"st_registro_ativo" db:"st_registro_ativo"`
	STChangeApprove string         `json:"st_change_approve" db:"st_change_approve"`
}

type IngredientSubstance struct {
	COSeqID         int64          `json:"co_seq_id" db:"co_seq_id"`
	NUIsID          string         `json:"nu_isid" db:"nu_isid"`
	DTIsIDDt        sql.NullInt64  `json:"dt_isiddt" db:"dt_isiddt"`
	NONm            string         `json:"no_nm" db:"no_nm"`
	NONmPtBr        sql.NullString `json:"no_nm_pt_br" db:"no_nm_pt_br"`
	COIsrCd         sql.NullInt64  `json:"co_isrcid" db:"co_isrcid"`
	CODcbCd         sql.NullInt64  `json:"co_dcbcd" db:"co_dcbcd"`
	CONCas          sql.NullString `json:"co_ncas" db:"co_ncas"`
	STRegistroAtivo string         `json:"st_registro_ativo" db:"st_registro_ativo"`
	STChangeApprove string         `json:"st_change_approve" db:"st_change_approve"`
}

type Domain struct {
	COSeqID           int64          `json:"co_seq_id" db:"co_seq_id"`
	NUCd              string         `json:"nu_cd" db:"nu_cd"`
	NODescr           string         `json:"no_descr" db:"no_descr"`
	NODescrPtBr       sql.NullString `json:"no_descr_pt_br" db:"no_descr_pt_br"`
	STRegistroAtivo   string         `json:"st_registro_ativo" db:"st_registro_ativo"`
	STChangeApprove   string         `json:"st_change_approve" db:"st_change_approve"`
	COParent          sql.NullInt64  `json:"co_parent" db:"co_parent"`
	STHierarchyLevel  sql.NullInt64  `json:"st_hierarchy_level" db:"st_hierarchy_level"`
	NOAbbrevNm        sql.NullString `json:"no_abbrevnm" db:"no_abbrevnm"`
	DSScope           sql.NullString `json:"ds_scope" db:"ds_scope"`
	DSLimitPresc      sql.NullString `json:"ds_limitpresc" db:"ds_limitpresc"`
	DSQtypPresc       sql.NullString `json:"ds_qtypresc" db:"ds_qtypresc"`
	DSValPresc        sql.NullString `json:"ds_valpresc" db:"ds_valpresc"`
	COAlpha2          sql.NullString `json:"co_alpha_2" db:"co_alpha_2"`
	COAlpha3          sql.NullString `json:"co_alpha_3" db:"co_alpha_3"`
	COLatitude        sql.NullFloat64 `json:"co_latitude" db:"co_latitude"`
	COLongitude       sql.NullFloat64 `json:"co_longitude" db:"co_longitude"`
	NOSigla           sql.NullString `json:"no_sigla" db:"no_sigla"`
	PNIDisponivel     sql.NullBool   `json:"pni_disponivel" db:"pni_disponivel"`
	NUCatmat          sql.NullInt64  `json:"nu_catmat" db:"nu_catmat"`
	DSCatmat          sql.NullString `json:"ds_catmat" db:"ds_catmat"`
	NOHorus           sql.NullString `json:"no_horus" db:"no_horus"`
	DTCdDt            sql.NullInt64  `json:"dt_cddt" db:"dt_cddt"`
	COFFacCd          sql.NullInt64  `json:"co_ffacd" db:"co_ffacd"`
	NUVersion         sql.NullInt64  `json:"nu_version" db:"nu_version"`
}

type User struct {
	ID           int64     `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	PasswordHash string    `json:"-" db:"password_hash"`
	Active       bool      `json:"active" db:"active"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

type SearchHit struct {
	Entity         string  `json:"entity"`
	ID             int64   `json:"id"`
	Nome           string  `json:"nome"`
	Codigo         string  `json:"codigo"`
	Fabricante     string  `json:"fabricante,omitempty"`
	Descricao      string  `json:"descricao,omitempty"`
	RelevanceScore float64 `json:"relevance_score"`
}

type VMPDetail struct {
	VMP              VMP              `json:"vmp"`
	VTM              *VTM             `json:"vtm,omitempty"`
	BasisOfName      *Domain          `json:"basis_of_name,omitempty"`
	PresStatus       *Domain          `json:"pres_status,omitempty"`
	DFIndicator      *Domain          `json:"df_indicator,omitempty"`
	BrImunologico    *Domain          `json:"br_imunologico,omitempty"`
	AnvsClass        *Domain          `json:"anvs_class,omitempty"`
	Forms            []Domain         `json:"forms,omitempty"`
	Routes           []Domain         `json:"routes,omitempty"`
	ATCClasses       []Domain         `json:"atc_classes,omitempty"`
	Catmats          []Domain         `json:"catmats,omitempty"`
	Ingredients      []VMPIngredient  `json:"ingredients,omitempty"`
	ControlDrugInfos []Domain         `json:"control_drug_infos,omitempty"`
	RenamesBR        []Domain         `json:"renames_br,omitempty"`
	LocalAplicacao   []Domain         `json:"local_aplicacao,omitempty"`
}

type VMPIngredient struct {
	COSeqID             int64           `json:"co_seq_id" db:"co_seq_id"`
	COVmpID             int64           `json:"co_vmpid" db:"co_vmpid"`
	COIscd              int64           `json:"co_iscd" db:"co_iscd"`
	COBsSubID           sql.NullInt64   `json:"co_bs_subid" db:"co_bs_subid"`
	COBasisStrntCd      sql.NullInt64   `json:"co_basis_strntcd" db:"co_basis_strntcd"`
	QTStrntNmrtrVal     sql.NullFloat64 `json:"qt_strnt_nmrtr_val" db:"qt_strnt_nmrtr_val"`
	COStrntNmrtrUomCd   sql.NullInt64   `json:"co_strnt_nmrtr_uomcd" db:"co_strnt_nmrtr_uomcd"`
	QTStrntDnmtrVal     sql.NullFloat64 `json:"qt_strnt_dnmtr_val" db:"qt_strnt_dnmtr_val"`
	COStrntDnmtrUomCd   sql.NullInt64   `json:"co_strnt_dnmtr_uomcd" db:"co_strnt_dnmtr_uomcd"`
	IngredientSubstance *IngredientSubstance `json:"ingredient_substance,omitempty"`
}

type AMPDetail struct {
	AMP             AMP              `json:"amp"`
	VMP             *VMP             `json:"vmp,omitempty"`
	Supplier        *Supplier        `json:"supplier,omitempty"`
	Flavour         *Domain          `json:"flavour,omitempty"`
	LicAuth         *Domain          `json:"lic_auth,omitempty"`
	MedClass        *Domain          `json:"med_class,omitempty"`
	AvailRestriction *Domain         `json:"avail_restriction,omitempty"`
	Routes          []Domain         `json:"routes,omitempty"`
	PreservConds    []Domain         `json:"preserv_conds,omitempty"`
	Ingredients     []AMPIngredient  `json:"ingredients,omitempty"`
}

type AMPIngredient struct {
	COSeqID       int64           `json:"co_seq_id" db:"co_seq_id"`
	COIsID        int64           `json:"co_isid" db:"co_isid"`
	COApID        int64           `json:"co_apid" db:"co_apid"`
	QTStrnth      sql.NullFloat64 `json:"qt_strnth" db:"qt_strnth"`
	COUomCd       sql.NullInt64   `json:"co_uomcd" db:"co_uomcd"`
	IngredientSubstance *IngredientSubstance `json:"ingredient_substance,omitempty"`
}

type CursorPage[T any] struct {
	Items  []T    `json:"items"`
	Cursor string `json:"cursor,omitempty"`
	Limit  int    `json:"limit"`
	Total  int64  `json:"total"`
}
