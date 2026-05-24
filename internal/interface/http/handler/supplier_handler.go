package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/silasrm/api-obm/internal/interface/http/dto"
	"github.com/silasrm/api-obm/internal/usecase"
)

type SupplierHandler struct {
	supplierUseCase *usecase.SupplierUsecase
}

func NewSupplierHandler(uc *usecase.SupplierUsecase) *SupplierHandler {
	return &SupplierHandler{supplierUseCase: uc}
}

// ListSuppliers godoc
// @Summary Listar fornecedores
// @Description Lista fornecedores com paginação e filtros
// @Tags Supplier
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limite por página" default(20)
// @Param cursor query string false "Cursor de paginação"
// @Param nome query string false "Filtro por nome (ILIKE)"
// @Param codigo query string false "Filtro por código NU_CD"
// @Param ativo query bool false "Filtro por status ativo"
// @Success 200 {object} dto.ListResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/suppliers [get]
func (h *SupplierHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	cursor := c.Query("cursor")
	filter := buildFilterParams(limit, cursor)
	filter.Nome = c.Query("nome")
	filter.Codigo = c.Query("codigo")
	if ativo := c.Query("ativo"); ativo != "" {
		b := ativo == "true"
		filter.Ativo = &b
	}

	page, err := h.supplierUseCase.List(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error(), Code: http.StatusInternalServerError})
		return
	}

	c.JSON(http.StatusOK, dto.ListResponse{Items: page.Items, Cursor: page.Cursor, Limit: page.Limit, Total: page.Total})
}

// GetSupplier godoc
// @Summary Obter fornecedor por ID
// @Description Retorna um fornecedor pelo seu ID
// @Tags Supplier
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID do fornecedor"
// @Success 200 {object} dto.SupplierResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/suppliers/{id} [get]
func (h *SupplierHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id", Code: http.StatusBadRequest})
		return
	}

	supplier, err := h.supplierUseCase.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "not found", Code: http.StatusNotFound})
		return
	}

	c.JSON(http.StatusOK, supplier)
}
