package handler

import (
	"context"
	"database/sql"
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

type mockCMEDRepoForHandler struct {
	cmed      *entity.CMEDConformidade
	page      *entity.CursorPage[entity.CMEDConformidade]
	historico []entity.CMEDConformidade
	err       error
}

func (m *mockCMEDRepoForHandler) GetByID(ctx context.Context, id int64) (*entity.CMEDConformidade, error) {
	return m.cmed, m.err
}

func (m *mockCMEDRepoForHandler) GetByNuSanReg(ctx context.Context, nuSanReg int64, dtReferencia string) (*entity.CMEDConformidade, error) {
	return m.cmed, m.err
}

func (m *mockCMEDRepoForHandler) GetByEAN(ctx context.Context, ean string, dtReferencia string) (*entity.CMEDConformidade, error) {
	return m.cmed, m.err
}

func (m *mockCMEDRepoForHandler) List(ctx context.Context, filter repository.CMEDFilterParams) (*entity.CursorPage[entity.CMEDConformidade], error) {
	return m.page, m.err
}

func (m *mockCMEDRepoForHandler) GetHistorico(ctx context.Context, nuSanReg int64) ([]entity.CMEDConformidade, error) {
	return m.historico, m.err
}

func (m *mockCMEDRepoForHandler) UpsertBatch(ctx context.Context, records []entity.CMEDConformidade) (int64, error) {
	return 0, nil
}

type mockAMPPRepoForHandler struct {
	ampp interface{}
	err  error
}

func (m *mockAMPPRepoForHandler) GetByID(ctx context.Context, id int64) (interface{}, error) {
	return m.ampp, m.err
}

func (m *mockAMPPRepoForHandler) List(ctx context.Context, filter repository.FilterParams) (interface{}, error) {
	return nil, nil
}

type mockVMPRepoForCMEDHandler struct {
	vmp *entity.VMP
	err error
}

func (m *mockVMPRepoForCMEDHandler) GetByID(ctx context.Context, id int64) (*entity.VMP, error) {
	return m.vmp, m.err
}

func (m *mockVMPRepoForCMEDHandler) GetDetailByID(ctx context.Context, id int64) (*entity.VMPDetail, error) {
	return nil, nil
}

func (m *mockVMPRepoForCMEDHandler) List(ctx context.Context, filter repository.FilterParams) (*entity.CursorPage[entity.VMP], error) {
	return nil, nil
}

type mockAMPRepoForCMEDHandler struct {
	amp *entity.AMP
	err error
}

func (m *mockAMPRepoForCMEDHandler) GetByID(ctx context.Context, id int64) (*entity.AMP, error) {
	return m.amp, m.err
}

func (m *mockAMPRepoForCMEDHandler) GetDetailByID(ctx context.Context, id int64) (*entity.AMPDetail, error) {
	return nil, nil
}

func (m *mockAMPRepoForCMEDHandler) List(ctx context.Context, filter repository.FilterParams) (*entity.CursorPage[entity.AMP], error) {
	return nil, nil
}

func newCMEDHandler(repos ...interface{}) *usecase.CMEDUsecase {
	repo := &mockCMEDRepoForHandler{}
	if len(repos) > 0 {
		if r, ok := repos[0].(*mockCMEDRepoForHandler); ok {
			repo = r
		}
	}
	return usecase.NewCMEDUsecase(repo, nil)
}

func newAMPPCMEDHandler(amppRepo repository.AMPPRepository, vmpRepo repository.VMPRepository, ampRepo repository.AMPRepository, cmedRepo repository.CMEDRepository) *usecase.AMPPCMEDUsecase {
	return usecase.NewAMPPCMEDUsecase(amppRepo, vmpRepo, ampRepo, cmedRepo, nil)
}

func TestCMEDHandler_List_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	sanReg := int64(12345)
	page := &entity.CursorPage[entity.CMEDConformidade]{
		Items: []entity.CMEDConformidade{{COSeqID: 1, NUSanReg: &sanReg, NOProduto: strPtrHandler("Test")}},
		Total: 1,
		Limit: 20,
	}
	repo := &mockCMEDRepoForHandler{page: page}
	uc := usecase.NewCMEDUsecase(repo, nil)
	amppUC := newAMPPCMEDHandler(&mockAMPPRepoForHandler{}, &mockVMPRepoForCMEDHandler{}, &mockAMPRepoForCMEDHandler{}, repo)
	handler := NewCMEDHandler(uc, amppUC)

	r := gin.New()
	r.GET("/api/v1/cmed", handler.List)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/cmed?nome=test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCMEDHandler_GetByID_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	sanReg := int64(12345)
	cmed := &entity.CMEDConformidade{COSeqID: 1, NUSanReg: &sanReg}
	repo := &mockCMEDRepoForHandler{cmed: cmed}
	uc := usecase.NewCMEDUsecase(repo, nil)
	amppUC := newAMPPCMEDHandler(&mockAMPPRepoForHandler{}, &mockVMPRepoForCMEDHandler{}, &mockAMPRepoForCMEDHandler{}, repo)
	handler := NewCMEDHandler(uc, amppUC)

	r := gin.New()
	r.GET("/api/v1/cmed/:id", handler.GetByID)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/cmed/1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCMEDHandler_GetByID_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &mockCMEDRepoForHandler{cmed: nil, err: fmt.Errorf("not found")}
	uc := usecase.NewCMEDUsecase(repo, nil)
	amppUC := newAMPPCMEDHandler(&mockAMPPRepoForHandler{}, &mockVMPRepoForCMEDHandler{}, &mockAMPRepoForCMEDHandler{}, &mockCMEDRepoForHandler{})
	handler := NewCMEDHandler(uc, amppUC)

	r := gin.New()
	r.GET("/api/v1/cmed/:id", handler.GetByID)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/cmed/999", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCMEDHandler_GetByID_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &mockCMEDRepoForHandler{}
	uc := usecase.NewCMEDUsecase(repo, nil)
	amppUC := newAMPPCMEDHandler(&mockAMPPRepoForHandler{}, &mockVMPRepoForCMEDHandler{}, &mockAMPRepoForCMEDHandler{}, repo)
	handler := NewCMEDHandler(uc, amppUC)

	r := gin.New()
	r.GET("/api/v1/cmed/:id", handler.GetByID)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/cmed/abc", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCMEDHandler_GetByRegistro_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	sanReg := int64(12345)
	cmed := &entity.CMEDConformidade{COSeqID: 1, NUSanReg: &sanReg}
	repo := &mockCMEDRepoForHandler{cmed: cmed}
	uc := usecase.NewCMEDUsecase(repo, nil)
	amppUC := newAMPPCMEDHandler(&mockAMPPRepoForHandler{}, &mockVMPRepoForCMEDHandler{}, &mockAMPRepoForCMEDHandler{}, repo)
	handler := NewCMEDHandler(uc, amppUC)

	r := gin.New()
	r.GET("/api/v1/cmed/registro/:registro", handler.GetByRegistro)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/cmed/registro/123", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCMEDHandler_GetByRegistro_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &mockCMEDRepoForHandler{}
	uc := usecase.NewCMEDUsecase(repo, nil)
	amppUC := newAMPPCMEDHandler(&mockAMPPRepoForHandler{}, &mockVMPRepoForCMEDHandler{}, &mockAMPRepoForCMEDHandler{}, repo)
	handler := NewCMEDHandler(uc, amppUC)

	r := gin.New()
	r.GET("/api/v1/cmed/registro/:registro", handler.GetByRegistro)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/cmed/registro/abc", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCMEDHandler_GetByEAN_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cmed := &entity.CMEDConformidade{COSeqID: 1, NUEAN1: strPtrHandler("789123")}
	repo := &mockCMEDRepoForHandler{cmed: cmed}
	uc := usecase.NewCMEDUsecase(repo, nil)
	amppUC := newAMPPCMEDHandler(&mockAMPPRepoForHandler{}, &mockVMPRepoForCMEDHandler{}, &mockAMPRepoForCMEDHandler{}, repo)
	handler := NewCMEDHandler(uc, amppUC)

	r := gin.New()
	r.GET("/api/v1/cmed/ean/:ean", handler.GetByEAN)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/cmed/ean/789123", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCMEDHandler_GetByEAN_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &mockCMEDRepoForHandler{cmed: nil, err: nil}
	uc := usecase.NewCMEDUsecase(repo, nil)
	amppUC := newAMPPCMEDHandler(&mockAMPPRepoForHandler{}, &mockVMPRepoForCMEDHandler{}, &mockAMPRepoForCMEDHandler{}, repo)
	handler := NewCMEDHandler(uc, amppUC)

	r := gin.New()
	r.GET("/api/v1/cmed/ean/:ean", handler.GetByEAN)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/cmed/ean/000", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCMEDHandler_GetHistorico_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	sanReg := int64(12345)
	cmed := &entity.CMEDConformidade{COSeqID: 1, NUSanReg: &sanReg}
	historico := []entity.CMEDConformidade{
		{COSeqID: 1, NUSanReg: &sanReg, DTReferencia: "2024-01-01"},
		{COSeqID: 2, NUSanReg: &sanReg, DTReferencia: "2024-02-01"},
	}
	repo := &mockCMEDRepoForHandler{cmed: cmed, historico: historico}
	uc := usecase.NewCMEDUsecase(repo, nil)
	amppUC := newAMPPCMEDHandler(&mockAMPPRepoForHandler{}, &mockVMPRepoForCMEDHandler{}, &mockAMPRepoForCMEDHandler{}, repo)
	handler := NewCMEDHandler(uc, amppUC)

	r := gin.New()
	r.GET("/api/v1/cmed/:id/historico", handler.GetHistorico)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/cmed/1/historico", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCMEDHandler_GetHistorico_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &mockCMEDRepoForHandler{}
	uc := usecase.NewCMEDUsecase(repo, nil)
	amppUC := newAMPPCMEDHandler(&mockAMPPRepoForHandler{}, &mockVMPRepoForCMEDHandler{}, &mockAMPRepoForCMEDHandler{}, repo)
	handler := NewCMEDHandler(uc, amppUC)

	r := gin.New()
	r.GET("/api/v1/cmed/:id/historico", handler.GetHistorico)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/cmed/abc/historico", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCMEDHandler_GetAMPPWithCMED_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	sanReg := sql.NullInt64{Int64: 12345, Valid: true}
	ampp := &entity.AMPP{COSeqID: 1, COApID: 10, NUSanReg: sanReg}
	amp := &entity.AMP{COSeqID: 10, COVpID: 20}
	vmp := &entity.VMP{COSeqID: 20, NONm: "Test VMP"}
	cmedSanReg := int64(12345)
	cmed := &entity.CMEDConformidade{COSeqID: 100, NUSanReg: &cmedSanReg}

	cmedRepo := &mockCMEDRepoForHandler{cmed: cmed}
	amppRepo := &mockAMPPRepoForHandler{ampp: ampp}
	ampRepo := &mockAMPRepoForCMEDHandler{amp: amp}
	vmpRepo := &mockVMPRepoForCMEDHandler{vmp: vmp}

	uc := usecase.NewCMEDUsecase(cmedRepo, nil)
	amppUC := newAMPPCMEDHandler(amppRepo, vmpRepo, ampRepo, cmedRepo)
	handler := NewCMEDHandler(uc, amppUC)

	r := gin.New()
	r.GET("/api/v1/cmed/ampp/:id", handler.GetAMPPWithCMED)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/cmed/ampp/1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func strPtrHandler(s string) *string {
	return &s
}
