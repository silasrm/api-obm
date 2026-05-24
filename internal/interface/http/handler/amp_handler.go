package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/silasrm/api-obm/internal/interface/http/dto"
	"github.com/silasrm/api-obm/internal/usecase"
)

type AMPHandler struct {
	ampUseCase *usecase.AMPUsecase
}

func NewAMPHandler(uc *usecase.AMPUsecase) *AMPHandler {
	return &AMPHandler{ampUseCase: uc}
}

// ListAMP godoc
// @Summary Listar AMPs
// @Description Lista Actual Medicinal Products com paginação e filtros
// @Tags AMP
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limite por página" default(20)
// @Param cursor query string false "Cursor de paginação"
// @Param nome query string false "Filtro por nome (ILIKE)"
// @Param codigo query string false "Filtro por código NU_APID"
// @Param fabricante query string false "Filtro por nome do fabricante (ILIKE)"
// @Param ativo query bool false "Filtro por status ativo"
// @Success 200 {object} dto.ListResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/amp [get]
func (h *AMPHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	cursor := c.Query("cursor")
	filter := buildFilterParams(limit, cursor)
	filter.Nome = c.Query("nome")
	filter.Codigo = c.Query("codigo")
	filter.Fabricante = c.Query("fabricante")
	if ativo := c.Query("ativo"); ativo != "" {
		b := ativo == "true"
		filter.Ativo = &b
	}

	page, err := h.ampUseCase.List(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error(), Code: http.StatusInternalServerError})
		return
	}

	c.JSON(http.StatusOK, dto.ListResponse{Items: page.Items, Cursor: page.Cursor, Limit: page.Limit, Total: page.Total})
}

// GetAMP godoc
// @Summary Obter AMP por ID
// @Description Retorna um Actual Medicinal Product pelo seu ID
// @Tags AMP
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID do AMP"
// @Success 200 {object} dto.AMPResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/amp/{id} [get]
func (h *AMPHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id", Code: http.StatusBadRequest})
		return
	}

	amp, err := h.ampUseCase.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "not found", Code: http.StatusNotFound})
		return
	}

	c.JSON(http.StatusOK, amp)
}

// GetAMPDetail godoc
// @Summary Obter detalhes do AMP
// @Description Retorna um AMP com detalhes completos (VMP, fornecedor, domínios, ingredientes, rotas, etc.)
// @Tags AMP
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID do AMP"
// @Success 200 {object} dto.AMPDetailResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/amp/{id}/detail [get]
func (h *AMPHandler) GetDetail(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id", Code: http.StatusBadRequest})
		return
	}

	detail, err := h.ampUseCase.GetDetailByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "not found", Code: http.StatusNotFound})
		return
	}

	c.JSON(http.StatusOK, detail)
}
