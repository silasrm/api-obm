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
	cmed        *entity.CMEDConformidade
	page        *entity.CursorPage[entity.CMEDConformidade]
	historico   []entity.CMEDConformidade
	cnpjResults []entity.CMEDConformidade
	err         error
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

func (m *mockCMEDRepoForHandler) GetByCNPJ(ctx context.Context, cnpj string, dtReferencia string) ([]entity.CMEDConformidade, error) {
	return m.cnpjResults, m.err
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

type mockSupplierRepoForHandler struct {
	supplier interface{}
	err      error
}

func (m *mockSupplierRepoForHandler) GetByID(ctx context.Context, id int64) (interface{}, error) {
	return m.supplier, m.err
}

func (m *mockSupplierRepoForHandler) List(ctx context.Context, filter repository.FilterParams) (interface{}, error) {
	return nil, nil
}

func newCMEDHandlerWithMocks(cmedRepo *mockCMEDRepoForHandler, amppRepo *mockAMPPRepoForHandler, ampRepo *mockAMPRepoForCMEDHandler, vmpRepo *mockVMPRepoForCMEDHandler, supplierRepo *mockSupplierRepoForHandler) *CMEDHandler {
	uc := usecase.NewCMEDUsecase(cmedRepo, nil)
	amppUC := usecase.NewAMPPCMEDUsecase(amppRepo, vmpRepo, ampRepo, cmedRepo, nil)
	ampUC := usecase.NewAMPCMEDUsecase(ampRepo, vmpRepo, cmedRepo, nil)
	supplierUC := usecase.NewSupplierCMEDUsecase(supplierRepo, cmedRepo, nil)
	return NewCMEDHandler(uc, amppUC, ampUC, supplierUC)
}

func defaultMockRepos(cmedRepo *mockCMEDRepoForHandler) (*mockAMPPRepoForHandler, *mockAMPRepoForCMEDHandler, *mockVMPRepoForCMEDHandler, *mockSupplierRepoForHandler) {
	return &mockAMPPRepoForHandler{}, &mockAMPRepoForCMEDHandler{}, &mockVMPRepoForCMEDHandler{}, &mockSupplierRepoForHandler{}
}

func TestCMEDHandler_List_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	sanReg := int64(12345)
	page := &entity.CursorPage[entity.CMEDConformidade]{
		Items: []entity.CMEDConformidade{{COSeqID: 1, NUSanReg: &sanReg, NOProduto: strPtrHandler("Test")}},
		Total: 1, Limit: 20,
	}
	repo := &mockCMEDRepoForHandler{page: page}
	ampp, amp, vmp, sup := defaultMockRepos(repo)
	handler := newCMEDHandlerWithMocks(repo, ampp, amp, vmp, sup)

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
	repo := &mockCMEDRepoForHandler{cmed: &entity.CMEDConformidade{COSeqID: 1, NUSanReg: &sanReg}}
	ampp, amp, vmp, sup := defaultMockRepos(repo)
	handler := newCMEDHandlerWithMocks(repo, ampp, amp, vmp, sup)

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
	ampp, amp, vmp, sup := defaultMockRepos(repo)
	handler := newCMEDHandlerWithMocks(repo, ampp, amp, vmp, sup)

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
	ampp, amp, vmp, sup := defaultMockRepos(repo)
	handler := newCMEDHandlerWithMocks(repo, ampp, amp, vmp, sup)

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
	repo := &mockCMEDRepoForHandler{cmed: &entity.CMEDConformidade{COSeqID: 1, NUSanReg: &sanReg}}
	ampp, amp, vmp, sup := defaultMockRepos(repo)
	handler := newCMEDHandlerWithMocks(repo, ampp, amp, vmp, sup)

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
	ampp, amp, vmp, sup := defaultMockRepos(repo)
	handler := newCMEDHandlerWithMocks(repo, ampp, amp, vmp, sup)

	r := gin.New()
	r.GET("/api/v1/cmed/registro/:registro", handler.GetByRegistro)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/cmed/registro/abc", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCMEDHandler_GetByEAN_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &mockCMEDRepoForHandler{cmed: &entity.CMEDConformidade{COSeqID: 1, NUEAN1: strPtrHandler("789123")}}
	ampp, amp, vmp, sup := defaultMockRepos(repo)
	handler := newCMEDHandlerWithMocks(repo, ampp, amp, vmp, sup)

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
	ampp, amp, vmp, sup := defaultMockRepos(repo)
	handler := newCMEDHandlerWithMocks(repo, ampp, amp, vmp, sup)

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
	historico := []entity.CMEDConformidade{
		{COSeqID: 1, NUSanReg: &sanReg, DTReferencia: "2024-01-01"},
		{COSeqID: 2, NUSanReg: &sanReg, DTReferencia: "2024-02-01"},
	}
	repo := &mockCMEDRepoForHandler{historico: historico}
	ampp, amp, vmp, sup := defaultMockRepos(repo)
	handler := newCMEDHandlerWithMocks(repo, ampp, amp, vmp, sup)

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
	ampp, amp, vmp, sup := defaultMockRepos(repo)
	handler := newCMEDHandlerWithMocks(repo, ampp, amp, vmp, sup)

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
	amppEnt := &entity.AMPP{COSeqID: 1, COApID: 10, NUSanReg: sanReg}
	ampEnt := &entity.AMP{COSeqID: 10, COVpID: 20}
	vmpEnt := &entity.VMP{COSeqID: 20, NONm: "Test VMP"}
	cmedSanReg := int64(12345)
	cmedRepo := &mockCMEDRepoForHandler{cmed: &entity.CMEDConformidade{COSeqID: 100, NUSanReg: &cmedSanReg}}
	amppRepo := &mockAMPPRepoForHandler{ampp: amppEnt}
	ampRepo := &mockAMPRepoForCMEDHandler{amp: ampEnt}
	vmpRepo := &mockVMPRepoForCMEDHandler{vmp: vmpEnt}
	supRepo := &mockSupplierRepoForHandler{}
	handler := newCMEDHandlerWithMocks(cmedRepo, amppRepo, ampRepo, vmpRepo, supRepo)

	r := gin.New()
	r.GET("/api/v1/cmed/ampp/:id", handler.GetAMPPWithCMED)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/cmed/ampp/1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCMEDHandler_GetAMPWithCMED_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	nuNReg := sql.NullInt64{Int64: 12345, Valid: true}
	ampEnt := &entity.AMP{COSeqID: 10, COVpID: 20, NUNReg: nuNReg}
	vmpEnt := &entity.VMP{COSeqID: 20, NONm: "Test VMP"}
	cmedSanReg := int64(12345)
	cmedRepo := &mockCMEDRepoForHandler{cmed: &entity.CMEDConformidade{COSeqID: 100, NUSanReg: &cmedSanReg}}
	amppRepo := &mockAMPPRepoForHandler{}
	ampRepo := &mockAMPRepoForCMEDHandler{amp: ampEnt}
	vmpRepo := &mockVMPRepoForCMEDHandler{vmp: vmpEnt}
	supRepo := &mockSupplierRepoForHandler{}
	handler := newCMEDHandlerWithMocks(cmedRepo, amppRepo, ampRepo, vmpRepo, supRepo)

	r := gin.New()
	r.GET("/api/v1/amp/:id/cmed", handler.GetAMPWithCMED)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/amp/10/cmed", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCMEDHandler_GetAMPWithCMED_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cmedRepo := &mockCMEDRepoForHandler{}
	ampRepo := &mockAMPRepoForCMEDHandler{err: fmt.Errorf("not found")}
	amppRepo := &mockAMPPRepoForHandler{}
	vmpRepo := &mockVMPRepoForCMEDHandler{}
	supRepo := &mockSupplierRepoForHandler{}
	handler := newCMEDHandlerWithMocks(cmedRepo, amppRepo, ampRepo, vmpRepo, supRepo)

	r := gin.New()
	r.GET("/api/v1/amp/:id/cmed", handler.GetAMPWithCMED)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/amp/999/cmed", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCMEDHandler_GetAMPWithCMED_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cmedRepo := &mockCMEDRepoForHandler{}
	ampp, amp, vmp, sup := defaultMockRepos(cmedRepo)
	handler := newCMEDHandlerWithMocks(cmedRepo, ampp, amp, vmp, sup)

	r := gin.New()
	r.GET("/api/v1/amp/:id/cmed", handler.GetAMPWithCMED)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/amp/abc/cmed", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCMEDHandler_GetSupplierWithCMED_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	supplier := &entity.Supplier{COSeqID: 5, NODescr: "Test Lab", NUCnpj: sql.NullString{String: "12345678000100", Valid: true}}
	cmedRepo := &mockCMEDRepoForHandler{cnpjResults: []entity.CMEDConformidade{{COSeqID: 100}}}
	amppRepo := &mockAMPPRepoForHandler{}
	ampRepo := &mockAMPRepoForCMEDHandler{}
	vmpRepo := &mockVMPRepoForCMEDHandler{}
	supRepo := &mockSupplierRepoForHandler{supplier: supplier}
	handler := newCMEDHandlerWithMocks(cmedRepo, amppRepo, ampRepo, vmpRepo, supRepo)

	r := gin.New()
	r.GET("/api/v1/suppliers/:id/cmed", handler.GetSupplierWithCMED)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/suppliers/5/cmed", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCMEDHandler_GetSupplierWithCMED_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cmedRepo := &mockCMEDRepoForHandler{}
	supRepo := &mockSupplierRepoForHandler{err: fmt.Errorf("not found")}
	amppRepo := &mockAMPPRepoForHandler{}
	ampRepo := &mockAMPRepoForCMEDHandler{}
	vmpRepo := &mockVMPRepoForCMEDHandler{}
	handler := newCMEDHandlerWithMocks(cmedRepo, amppRepo, ampRepo, vmpRepo, supRepo)

	r := gin.New()
	r.GET("/api/v1/suppliers/:id/cmed", handler.GetSupplierWithCMED)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/suppliers/999/cmed", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCMEDHandler_GetSupplierWithCMED_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cmedRepo := &mockCMEDRepoForHandler{}
	ampp, amp, vmp, sup := defaultMockRepos(cmedRepo)
	handler := newCMEDHandlerWithMocks(cmedRepo, ampp, amp, vmp, sup)

	r := gin.New()
	r.GET("/api/v1/suppliers/:id/cmed", handler.GetSupplierWithCMED)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/suppliers/abc/cmed", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func strPtrHandler(s string) *string { return &s }
