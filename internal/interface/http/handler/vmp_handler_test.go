package handler

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/silasrm/api-obm/internal/domain/entity"
	"github.com/silasrm/api-obm/internal/domain/repository"
	"github.com/silasrm/api-obm/internal/usecase"
	"github.com/stretchr/testify/assert"
)

type mockVMPRepoForHandler struct {
	vmp    *entity.VMP
	detail *entity.VMPDetail
	page   *entity.CursorPage[entity.VMP]
	err    error
}

func (m *mockVMPRepoForHandler) GetByID(ctx context.Context, id int64) (*entity.VMP, error) {
	return m.vmp, m.err
}
func (m *mockVMPRepoForHandler) GetDetailByID(ctx context.Context, id int64) (*entity.VMPDetail, error) {
	return m.detail, m.err
}
func (m *mockVMPRepoForHandler) List(ctx context.Context, filter repository.FilterParams) (*entity.CursorPage[entity.VMP], error) {
	return m.page, m.err
}

func TestVMPHandler_GetByID_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &mockVMPRepoForHandler{}
	uc := usecase.NewVMPUsecase(repo)
	handler := NewVMPHandler(uc)

	r := gin.New()
	r.GET("/api/v1/vmp/:id", handler.GetByID)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/vmp/abc", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid id")
}

func TestVMPHandler_GetByID_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	vmp := &entity.VMP{COSeqID: 1, NONm: "Test VMP"}
	repo := &mockVMPRepoForHandler{vmp: vmp}
	uc := usecase.NewVMPUsecase(repo)
	handler := NewVMPHandler(uc)

	r := gin.New()
	r.GET("/api/v1/vmp/:id", handler.GetByID)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/vmp/1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestVMPHandler_GetByID_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &mockVMPRepoForHandler{vmp: nil, err: fmt.Errorf("not found")}
	uc := usecase.NewVMPUsecase(repo)
	handler := NewVMPHandler(uc)

	r := gin.New()
	r.GET("/api/v1/vmp/:id", handler.GetByID)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/vmp/999", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestVMPHandler_List(t *testing.T) {
	gin.SetMode(gin.TestMode)
	page := &entity.CursorPage[entity.VMP]{
		Items: []entity.VMP{{COSeqID: 1}},
		Total: 1,
		Limit: 20,
	}
	repo := &mockVMPRepoForHandler{page: page}
	uc := usecase.NewVMPUsecase(repo)
	handler := NewVMPHandler(uc)

	r := gin.New()
	r.GET("/api/v1/vmp", handler.List)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/vmp", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestVMPHandler_GetDetail_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &mockVMPRepoForHandler{}
	uc := usecase.NewVMPUsecase(repo)
	handler := NewVMPHandler(uc)

	r := gin.New()
	r.GET("/api/v1/vmp/:id/detail", handler.GetDetail)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/vmp/abc/detail", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestVMPHandler_List_WithFilters(t *testing.T) {
	gin.SetMode(gin.TestMode)
	page := &entity.CursorPage[entity.VMP]{
		Items: []entity.VMP{{COSeqID: 1}},
		Total: 1,
		Limit: 20,
	}
	repo := &mockVMPRepoForHandler{page: page}
	uc := usecase.NewVMPUsecase(repo)
	handler := NewVMPHandler(uc)

	r := gin.New()
	r.GET("/api/v1/vmp", handler.List)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/vmp?nome=test&ativo=true&limit=10", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
