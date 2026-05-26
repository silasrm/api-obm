package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/silasrm/api-obm/internal/domain/repository"
	"github.com/silasrm/api-obm/internal/interface/http/dto"
	"github.com/silasrm/api-obm/internal/usecase"
)

type SearchHandler struct {
	searchUseCase *usecase.SearchUsecase
}

func NewSearchHandler(uc *usecase.SearchUsecase) *SearchHandler {
	return &SearchHandler{searchUseCase: uc}
}

// Search godoc
// @Summary Busca global
// @Description Busca medicamentos e fornecedores via Meilisearch com filtros
// @Tags Search
// @Produce json
// @Security BearerAuth
// @Param q query string true "Termo de busca"
// @Param entity query string false "Entidades para buscar (vmp,amp,supplier,cmed). Separadas por virgula"
// @Param limit query int false "Limite de resultados" default(20)
// @Param cursor query string false "Cursor de paginação"
// @Param filter[nome] query string false "Filtro por nome"
// @Param filter[codigo] query string false "Filtro por código"
// @Param filter[fabricante] query string false "Filtro por fabricante"
// @Param filter[descricao] query string false "Filtro por descrição"
// @Param filter[ativo] query string false "Filtro por status ativo"
// @Param filter[tarja] query string false "Filtro por tarja (CMED)"
// @Param filter[registro] query string false "Filtro por registro sanitário (CMED)"
// @Success 200 {object} dto.SearchResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/search [get]
func (h *SearchHandler) Search(c *gin.Context) {
	var req dto.SearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid query parameters", Code: http.StatusBadRequest})
		return
	}

	filters := make(map[string]string)
	if req.FilterNome != "" {
		filters["no_nm"] = req.FilterNome
	}
	if req.FilterCodigo != "" {
		filters["nu_vpid"] = req.FilterCodigo
	}
	if req.FilterFabricante != "" {
		filters["supplier_name"] = req.FilterFabricante
	}
	if req.FilterDescricao != "" {
		filters["ds_descr"] = req.FilterDescricao
	}
	if req.FilterAtivo != "" {
		filters["st_registro_ativo"] = req.FilterAtivo
	}
	if req.FilterTarja != "" {
		filters["ds_tarja"] = req.FilterTarja
	}
	if req.FilterRegistro != "" {
		filters["nu_sanreg"] = req.FilterRegistro
	}

	var entities []string
	if req.Entity != "" {
		for _, e := range strings.Split(req.Entity, ",") {
			e = strings.TrimSpace(e)
			if e != "" {
				entities = append(entities, e)
			}
		}
	}

	page, err := h.searchUseCase.Search(
		c.Request.Context(),
		req.Q,
		entities,
		filters,
		req.Limit,
		req.Cursor,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error(), Code: http.StatusInternalServerError})
		return
	}

	c.JSON(http.StatusOK, dto.SearchResponse{
		Query:    req.Q,
		Entities: entities,
		Hits:     page.Items,
		Cursor:   page.Cursor,
		Limit:    page.Limit,
		Total:    page.Total,
	})
}

func buildFilterParams(limit int, cursor string) repository.FilterParams {
	if limit <= 0 {
		limit = 20
	}
	return repository.FilterParams{
		Limit:  limit,
		Cursor: cursor,
	}
}
