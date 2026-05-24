package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/silasrm/api-obm/internal/interface/http/dto"
	"github.com/silasrm/api-obm/internal/usecase"
)

type VMPHandler struct {
	vmpUseCase *usecase.VMPUsecase
}

func NewVMPHandler(uc *usecase.VMPUsecase) *VMPHandler {
	return &VMPHandler{vmpUseCase: uc}
}

// ListVMP godoc
// @Summary Listar VMPs
// @Description Lista Virtual Medicinal Products com paginacao e filtros
// @Tags VMP
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limite por pagina" default(20)
// @Param cursor query string false "Cursor de paginacao"
// @Param nome query string false "Filtro por nome (ILIKE)"
// @Param codigo query string false "Filtro por codigo NU_VPID"
// @Param ativo query bool false "Filtro por status ativo"
// @Success 200 {object} dto.ListResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/vmp [get]
func (h *VMPHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	cursor := c.Query("cursor")
	filter := buildFilterParams(limit, cursor)
	filter.Nome = c.Query("nome")
	filter.Codigo = c.Query("codigo")
	if ativo := c.Query("ativo"); ativo != "" {
		b := ativo == "true"
		filter.Ativo = &b
	}

	page, err := h.vmpUseCase.List(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error(), Code: http.StatusInternalServerError})
		return
	}

	c.JSON(http.StatusOK, dto.ListResponse{Items: page.Items, Cursor: page.Cursor, Limit: page.Limit, Total: page.Total})
}

// GetVMP godoc
// @Summary Obter VMP por ID
// @Description Retorna um Virtual Medicinal Product pelo seu ID
// @Tags VMP
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID do VMP"
// @Success 200 {object} dto.VMPResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/vmp/{id} [get]
func (h *VMPHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id", Code: http.StatusBadRequest})
		return
	}

	vmp, err := h.vmpUseCase.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "not found", Code: http.StatusNotFound})
		return
	}

	c.JSON(http.StatusOK, vmp)
}

// GetVMPDetail godoc
// @Summary Obter detalhes do VMP
// @Description Retorna um VMP com detalhes completos (VTM, dominios, ingredientes, rotas, formas, etc.)
// @Tags VMP
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID do VMP"
// @Success 200 {object} dto.VMPDetailResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/vmp/{id}/detail [get]
func (h *VMPHandler) GetDetail(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id", Code: http.StatusBadRequest})
		return
	}

	detail, err := h.vmpUseCase.GetDetailByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "not found", Code: http.StatusNotFound})
		return
	}

	c.JSON(http.StatusOK, detail)
}
