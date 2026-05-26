package dto

import "github.com/silasrm/api-obm/internal/domain/entity"

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token     string `json:"token"`
	ExpiresIn int    `json:"expires_in"`
}

type SearchRequest struct {
	Q                string `form:"q"`
	Entity           string `form:"entity"`
	Limit            int    `form:"limit"`
	Cursor           string `form:"cursor"`
	FilterNome       string `form:"filter[nome]"`
	FilterCodigo     string `form:"filter[codigo]"`
	FilterFabricante string `form:"filter[fabricante]"`
	FilterDescricao  string `form:"filter[descricao]"`
	FilterAtivo     string `form:"filter[ativo]"`
	FilterTarja     string `form:"filter[tarja]"`
	FilterRegistro  string `form:"filter[registro]"`
}

type SearchResponse struct {
	Query    string              `json:"query"`
	Entities []string            `json:"entities"`
	Hits     []entity.SearchHit  `json:"hits"`
	Cursor   string              `json:"cursor"`
	Limit    int                 `json:"limit"`
	Total    int64               `json:"total"`
}

type ListResponse struct {
	Items  interface{} `json:"items"`
	Cursor string      `json:"cursor"`
	Limit  int         `json:"limit"`
	Total  int64       `json:"total"`
}

type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

type ReindexResponse struct {
	Status  string         `json:"status"`
	Indexed map[string]int64 `json:"indexed"`
}

type HealthResponse struct {
	Status       string `json:"status"`
	Postgres     string `json:"postgres"`
	Meilisearch  string `json:"meilisearch"`
}

type VMPResponse entity.VMP
type VMPDetailResponse entity.VMPDetail
type AMPResponse entity.AMP
type AMPDetailResponse entity.AMPDetail
type VTMResponse entity.VTM
type VMPPResponse entity.VMPP
type AMPPResponse entity.AMPP
type DCBResponse entity.DCB
type SupplierResponse entity.Supplier
type IngredientResponse entity.IngredientSubstance
type DomainResponse entity.Domain

type CMEDListRequest struct {
	Nome         string `form:"nome"`
	Registro     string `form:"registro"`
	EAN          string `form:"ean"`
	Tarja        string `form:"tarja"`
	TipoProduto  string `form:"tipo_produto"`
	RegimePreco  string `form:"regime_preco"`
	DTReferencia string `form:"dt_referencia"`
	Limit        int    `form:"limit"`
	Cursor       string `form:"cursor"`
}

type CMEDListResponse struct {
	Items interface{} `json:"items"`
	Cursor string     `json:"cursor,omitempty"`
	Limit  int        `json:"limit"`
	Total  int64      `json:"total"`
}
