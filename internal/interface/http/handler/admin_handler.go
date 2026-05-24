package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/meilisearch/meilisearch-go"
	"github.com/silasrm/api-obm/internal/interface/http/dto"
	"github.com/silasrm/api-obm/internal/usecase"
)

type AdminHandler struct {
	reindexUseCase *usecase.ReindexUsecase
	pool           *pgxpool.Pool
	meiliClient    meilisearch.ServiceManager
}

func NewAdminHandler(reindexUseCase *usecase.ReindexUsecase, pool *pgxpool.Pool, meiliClient meilisearch.ServiceManager) *AdminHandler {
	return &AdminHandler{
		reindexUseCase: reindexUseCase,
		pool:           pool,
		meiliClient:    meiliClient,
	}
}

// Reindex godoc
// @Summary Reindexar Meilisearch
// @Description Reindexa todos os dados (VMPs, AMPs, Fornecedores) no Meilisearch
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.ReindexResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/admin/reindex [post]
func (h *AdminHandler) Reindex(c *gin.Context) {
	indexed, err := h.reindexUseCase.Reindex(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error(), Code: http.StatusInternalServerError})
		return
	}

	c.JSON(http.StatusOK, dto.ReindexResponse{Status: "completed", Indexed: indexed})
}

// Health godoc
// @Summary Health check
// @Description Verifica o status do PostgreSQL e Meilisearch
// @Tags Health
// @Produce json
// @Success 200 {object} dto.HealthResponse
// @Router /health [get]
func (h *AdminHandler) Health(c *gin.Context) {
	pgStatus := "ok"
	if err := h.pool.Ping(context.Background()); err != nil {
		pgStatus = "error"
	}

	meiliStatus := "ok"
	if _, err := h.meiliClient.Health(); err != nil {
		meiliStatus = "error"
	}

	status := "ok"
	if pgStatus != "ok" || meiliStatus != "ok" {
		status = "degraded"
	}

	c.JSON(http.StatusOK, dto.HealthResponse{Status: status, Postgres: pgStatus, Meilisearch: meiliStatus})
}
