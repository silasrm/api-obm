package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/silasrm/api-obm/internal/interface/http/dto"
	"github.com/silasrm/api-obm/internal/usecase"
)

type DomainHandler struct {
	domainUseCase *usecase.DomainUsecase
}

func NewDomainHandler(uc *usecase.DomainUsecase) *DomainHandler {
	return &DomainHandler{domainUseCase: uc}
}

// ListDomains godoc
// @Summary Listar dominios por tipo
// @Description Lista domínios (tabelas td_*) por tipo com paginação e filtros
// @Tags Domain
// @Produce json
// @Security BearerAuth
// @Param domain path string true "Tipo de domínio (form, route, flavour, legal_category, licensing_authority, etc.)"
// @Param limit query int false "Limite por página" default(20)
// @Param cursor query string false "Cursor de paginação"
// @Param nome query string false "Filtro por nome (ILIKE)"
// @Param ativo query bool false "Filtro por status ativo"
// @Success 200 {object} dto.ListResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/domains/{domain} [get]
func (h *DomainHandler) List(c *gin.Context) {
	domainType := c.Param("domain")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	cursor := c.Query("cursor")
	filter := buildFilterParams(limit, cursor)
	filter.Nome = c.Query("nome")
	if ativo := c.Query("ativo"); ativo != "" {
		b := ativo == "true"
		filter.Ativo = &b
	}

	page, err := h.domainUseCase.List(c.Request.Context(), domainType, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error(), Code: http.StatusInternalServerError})
		return
	}

	c.JSON(http.StatusOK, dto.ListResponse{Items: page.Items, Cursor: page.Cursor, Limit: page.Limit, Total: page.Total})
}

// GetDomain godoc
// @Summary Obter domínio por tipo e ID
// @Description Retorna um domínio específico pelo tipo e ID
// @Tags Domain
// @Produce json
// @Security BearerAuth
// @Param domain path string true "Tipo de domínio"
// @Param id path int true "ID do domínio"
// @Success 200 {object} dto.DomainResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/domains/{domain}/{id} [get]
func (h *DomainHandler) GetByID(c *gin.Context) {
	domainType := c.Param("domain")
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id", Code: http.StatusBadRequest})
		return
	}

	domain, err := h.domainUseCase.GetByID(c.Request.Context(), domainType, id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "not found", Code: http.StatusNotFound})
		return
	}

	c.JSON(http.StatusOK, domain)
}
