package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/silasrm/api-obm/internal/interface/http/dto"
	"github.com/silasrm/api-obm/internal/usecase"
)

type GenericHandler struct {
	genericUseCase *usecase.GenericUsecase
}

func NewGenericHandler(uc *usecase.GenericUsecase) *GenericHandler {
	return &GenericHandler{genericUseCase: uc}
}

// ListVTM godoc
// @Summary Listar VTMs
// @Description Lista Virtual Therapeutic Moieties com paginacao e filtros
// @Tags VTM
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limite por pagina" default(20)
// @Param cursor query string false "Cursor de paginacao"
// @Param nome query string false "Filtro por nome (ILIKE)"
// @Param ativo query bool false "Filtro por status ativo"
// @Success 200 {object} dto.ListResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/vtm [get]
func (h *GenericHandler) ListVTM(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	cursor := c.Query("cursor")
	filter := buildFilterParams(limit, cursor)
	filter.Nome = c.Query("nome")
	if ativo := c.Query("ativo"); ativo != "" {
		b := ativo == "true"
		filter.Ativo = &b
	}

	page, err := h.genericUseCase.ListVTMs(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error(), Code: http.StatusInternalServerError})
		return
	}

	c.JSON(http.StatusOK, dto.ListResponse{Items: page.Items, Cursor: page.Cursor, Limit: page.Limit, Total: page.Total})
}

// GetVTM godoc
// @Summary Obter VTM por ID
// @Description Retorna um Virtual Therapeutic Moiety pelo seu ID
// @Tags VTM
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID do VTM"
// @Success 200 {object} dto.VTMResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/vtm/{id} [get]
func (h *GenericHandler) GetVTM(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id", Code: http.StatusBadRequest})
		return
	}

	vtm, err := h.genericUseCase.GetVTM(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "not found", Code: http.StatusNotFound})
		return
	}

	c.JSON(http.StatusOK, vtm)
}

// ListVMPP godoc
// @Summary Listar VMPPs
// @Description Lista Virtual Medicinal Product Packs com paginacao e filtros
// @Tags VMPP
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limite por pagina" default(20)
// @Param cursor query string false "Cursor de paginacao"
// @Param nome query string false "Filtro por nome (ILIKE)"
// @Param ativo query bool false "Filtro por status ativo"
// @Success 200 {object} dto.ListResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/vmpp [get]
func (h *GenericHandler) ListVMPP(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	cursor := c.Query("cursor")
	filter := buildFilterParams(limit, cursor)
	filter.Nome = c.Query("nome")
	if ativo := c.Query("ativo"); ativo != "" {
		b := ativo == "true"
		filter.Ativo = &b
	}

	page, err := h.genericUseCase.ListVMPPs(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error(), Code: http.StatusInternalServerError})
		return
	}

	c.JSON(http.StatusOK, dto.ListResponse{Items: page.Items, Cursor: page.Cursor, Limit: page.Limit, Total: page.Total})
}

// GetVMPP godoc
// @Summary Obter VMPP por ID
// @Description Retorna um Virtual Medicinal Product Pack pelo seu ID
// @Tags VMPP
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID do VMPP"
// @Success 200 {object} dto.VMPPResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/vmpp/{id} [get]
func (h *GenericHandler) GetVMPP(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id", Code: http.StatusBadRequest})
		return
	}

	vmpp, err := h.genericUseCase.GetVMPP(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "not found", Code: http.StatusNotFound})
		return
	}

	c.JSON(http.StatusOK, vmpp)
}

// ListAMPP godoc
// @Summary Listar AMPPs
// @Description Lista Actual Medicinal Product Packs com paginacao e filtros
// @Tags AMPP
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limite por pagina" default(20)
// @Param cursor query string false "Cursor de paginacao"
// @Param nome query string false "Filtro por nome (ILIKE)"
// @Param ativo query bool false "Filtro por status ativo"
// @Success 200 {object} dto.ListResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/ampp [get]
func (h *GenericHandler) ListAMPP(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	cursor := c.Query("cursor")
	filter := buildFilterParams(limit, cursor)
	filter.Nome = c.Query("nome")
	if ativo := c.Query("ativo"); ativo != "" {
		b := ativo == "true"
		filter.Ativo = &b
	}

	page, err := h.genericUseCase.ListAMPPs(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error(), Code: http.StatusInternalServerError})
		return
	}

	c.JSON(http.StatusOK, dto.ListResponse{Items: page.Items, Cursor: page.Cursor, Limit: page.Limit, Total: page.Total})
}

// GetAMPP godoc
// @Summary Obter AMPP por ID
// @Description Retorna um Actual Medicinal Product Pack pelo seu ID
// @Tags AMPP
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID do AMPP"
// @Success 200 {object} dto.AMPPResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/ampp/{id} [get]
func (h *GenericHandler) GetAMPP(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id", Code: http.StatusBadRequest})
		return
	}

	ampp, err := h.genericUseCase.GetAMPP(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "not found", Code: http.StatusNotFound})
		return
	}

	c.JSON(http.StatusOK, ampp)
}

// ListDCB godoc
// @Summary Listar DCBs
// @Description Lista Denominacoes Comuns Brasileiras com paginacao e filtros
// @Tags DCB
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limite por pagina" default(20)
// @Param cursor query string false "Cursor de paginacao"
// @Param nome query string false "Filtro por nome (ILIKE)"
// @Param ativo query bool false "Filtro por status ativo"
// @Success 200 {object} dto.ListResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/dcb [get]
func (h *GenericHandler) ListDCB(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	cursor := c.Query("cursor")
	filter := buildFilterParams(limit, cursor)
	filter.Nome = c.Query("nome")
	if ativo := c.Query("ativo"); ativo != "" {
		b := ativo == "true"
		filter.Ativo = &b
	}

	page, err := h.genericUseCase.ListDCBs(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error(), Code: http.StatusInternalServerError})
		return
	}

	c.JSON(http.StatusOK, dto.ListResponse{Items: page.Items, Cursor: page.Cursor, Limit: page.Limit, Total: page.Total})
}

// GetDCB godoc
// @Summary Obter DCB por ID
// @Description Retorna uma Denominacao Comum Brasileira pelo seu ID
// @Tags DCB
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID do DCB"
// @Success 200 {object} dto.DCBResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/dcb/{id} [get]
func (h *GenericHandler) GetDCB(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id", Code: http.StatusBadRequest})
		return
	}

	dcb, err := h.genericUseCase.GetDCB(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "not found", Code: http.StatusNotFound})
		return
	}

	c.JSON(http.StatusOK, dcb)
}

// ListIngredients godoc
// @Summary Listar ingredientes
// @Description Lista Ingredient Substances com paginacao e filtros
// @Tags Ingredient
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limite por pagina" default(20)
// @Param cursor query string false "Cursor de paginacao"
// @Param nome query string false "Filtro por nome (ILIKE)"
// @Param ativo query bool false "Filtro por status ativo"
// @Success 200 {object} dto.ListResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/ingredients [get]
func (h *GenericHandler) ListIngredients(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	cursor := c.Query("cursor")
	filter := buildFilterParams(limit, cursor)
	filter.Nome = c.Query("nome")
	if ativo := c.Query("ativo"); ativo != "" {
		b := ativo == "true"
		filter.Ativo = &b
	}

	page, err := h.genericUseCase.ListIngredients(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error(), Code: http.StatusInternalServerError})
		return
	}

	c.JSON(http.StatusOK, dto.ListResponse{Items: page.Items, Cursor: page.Cursor, Limit: page.Limit, Total: page.Total})
}

// GetIngredient godoc
// @Summary Obter ingrediente por ID
// @Description Retorna um Ingredient Substance pelo seu ID
// @Tags Ingredient
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID do ingrediente"
// @Success 200 {object} dto.IngredientResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/ingredients/{id} [get]
func (h *GenericHandler) GetIngredient(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id", Code: http.StatusBadRequest})
		return
	}

	ing, err := h.genericUseCase.GetIngredient(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "not found", Code: http.StatusNotFound})
		return
	}

	c.JSON(http.StatusOK, ing)
}
