package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/silasrm/api-obm/internal/domain/repository"
	"github.com/silasrm/api-obm/internal/interface/http/dto"
	"github.com/silasrm/api-obm/internal/usecase"
)

type CMEDHandler struct {
	cmedUC     *usecase.CMEDUsecase
	amppCmedUC *usecase.AMPPCMEDUsecase
}

func NewCMEDHandler(cmedUC *usecase.CMEDUsecase, amppCmedUC *usecase.AMPPCMEDUsecase) *CMEDHandler {
	return &CMEDHandler{cmedUC: cmedUC, amppCmedUC: amppCmedUC}
}

func (h *CMEDHandler) List(c *gin.Context) {
	var req dto.CMEDListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid query parameters", Code: http.StatusBadRequest})
		return
	}

	limit := req.Limit
	if limit <= 0 {
		limit = 20
	}

	filter := repository.CMEDFilterParams{
		Nome:         req.Nome,
		Registro:     req.Registro,
		EAN:          req.EAN,
		Tarja:        req.Tarja,
		TipoProduto:  req.TipoProduto,
		RegimePreco:  req.RegimePreco,
		DTReferencia: req.DTReferencia,
		Limit:        limit,
		Cursor:       req.Cursor,
	}

	page, err := h.cmedUC.List(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error(), Code: http.StatusInternalServerError})
		return
	}

	c.JSON(http.StatusOK, dto.CMEDListResponse{Items: page.Items, Cursor: page.Cursor, Limit: page.Limit, Total: page.Total})
}

func (h *CMEDHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id", Code: http.StatusBadRequest})
		return
	}

	cmed, err := h.cmedUC.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "not found", Code: http.StatusNotFound})
		return
	}

	c.JSON(http.StatusOK, cmed)
}

func (h *CMEDHandler) GetByRegistro(c *gin.Context) {
	registro, err := strconv.ParseInt(c.Param("registro"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid registro", Code: http.StatusBadRequest})
		return
	}

	dtReferencia := c.Query("dt_referencia")

	cmed, err := h.cmedUC.GetByRegistro(c.Request.Context(), registro, dtReferencia)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error(), Code: http.StatusInternalServerError})
		return
	}
	if cmed == nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "not found", Code: http.StatusNotFound})
		return
	}

	c.JSON(http.StatusOK, cmed)
}

func (h *CMEDHandler) GetByEAN(c *gin.Context) {
	ean := c.Param("ean")
	if ean == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid ean", Code: http.StatusBadRequest})
		return
	}

	dtReferencia := c.Query("dt_referencia")

	cmed, err := h.cmedUC.GetByEAN(c.Request.Context(), ean, dtReferencia)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error(), Code: http.StatusInternalServerError})
		return
	}
	if cmed == nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "not found", Code: http.StatusNotFound})
		return
	}

	c.JSON(http.StatusOK, cmed)
}

func (h *CMEDHandler) GetHistorico(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id", Code: http.StatusBadRequest})
		return
	}

	historico, err := h.cmedUC.GetHistorico(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error(), Code: http.StatusInternalServerError})
		return
	}

	c.JSON(http.StatusOK, historico)
}

func (h *CMEDHandler) GetAMPPWithCMED(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id", Code: http.StatusBadRequest})
		return
	}

	dtReferencia := c.Query("dt_referencia")

	result, err := h.amppCmedUC.GetAMPPWithCMED(c.Request.Context(), id, dtReferencia)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error(), Code: http.StatusInternalServerError})
		return
	}
	if result == nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "not found", Code: http.StatusNotFound})
		return
	}

	c.JSON(http.StatusOK, result)
}
