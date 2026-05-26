package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/silasrm/api-obm/internal/domain/entity"
	"github.com/silasrm/api-obm/internal/domain/repository"
	"github.com/stretchr/testify/assert"
)

type mockCMEDRepo struct {
	cmed         *entity.CMEDConformidade
	cmedList     []entity.CMEDConformidade
	page         *entity.CursorPage[entity.CMEDConformidade]
	historico    []entity.CMEDConformidade
	err          error
	listErr      error
	historicoErr error
	lastFilter   repository.CMEDFilterParams
}

func (m *mockCMEDRepo) GetByID(ctx context.Context, id int64) (*entity.CMEDConformidade, error) {
	return m.cmed, m.err
}

func (m *mockCMEDRepo) GetByNuSanReg(ctx context.Context, nuSanReg int64, dtReferencia string) (*entity.CMEDConformidade, error) {
	return m.cmed, m.err
}

func (m *mockCMEDRepo) GetByEAN(ctx context.Context, ean string, dtReferencia string) (*entity.CMEDConformidade, error) {
	return m.cmed, m.err
}

func (m *mockCMEDRepo) List(ctx context.Context, filter repository.CMEDFilterParams) (*entity.CursorPage[entity.CMEDConformidade], error) {
	m.lastFilter = filter
	if m.listErr != nil {
		return nil, m.listErr
	}
	return m.page, nil
}

func (m *mockCMEDRepo) GetHistorico(ctx context.Context, nuSanReg int64) ([]entity.CMEDConformidade, error) {
	if m.historicoErr != nil {
		return nil, m.historicoErr
	}
	return m.historico, nil
}

func (m *mockCMEDRepo) UpsertBatch(ctx context.Context, records []entity.CMEDConformidade) (int64, error) {
	return 0, nil
}

func TestCMEDUsecase_List_Success(t *testing.T) {
	page := &entity.CursorPage[entity.CMEDConformidade]{
		Items: []entity.CMEDConformidade{{COSeqID: 1, NOProduto: strPtr("Test")}},
		Total: 1,
	}
	uc := NewCMEDUsecase(&mockCMEDRepo{page: page}, nil)

	result, err := uc.List(context.Background(), repository.CMEDFilterParams{})
	assert.NoError(t, err)
	assert.Equal(t, page, result)
	assert.Len(t, result.Items, 1)
}

func TestCMEDUsecase_List_DefaultLimit(t *testing.T) {
	repo := &mockCMEDRepo{page: &entity.CursorPage[entity.CMEDConformidade]{Items: []entity.CMEDConformidade{}, Total: 0}}
	uc := NewCMEDUsecase(repo, nil)

	_, err := uc.List(context.Background(), repository.CMEDFilterParams{Limit: 0})
	assert.NoError(t, err)
	assert.Equal(t, 20, repo.lastFilter.Limit)
}

func TestCMEDUsecase_List_NegativeLimit(t *testing.T) {
	repo := &mockCMEDRepo{page: &entity.CursorPage[entity.CMEDConformidade]{Items: []entity.CMEDConformidade{}, Total: 0}}
	uc := NewCMEDUsecase(repo, nil)

	_, err := uc.List(context.Background(), repository.CMEDFilterParams{Limit: -5})
	assert.NoError(t, err)
	assert.Equal(t, 20, repo.lastFilter.Limit)
}

func TestCMEDUsecase_GetByID_Success(t *testing.T) {
	cmed := &entity.CMEDConformidade{COSeqID: 1, NOProduto: strPtr("Test Product")}
	uc := NewCMEDUsecase(&mockCMEDRepo{cmed: cmed}, nil)

	result, err := uc.GetByID(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, cmed, result)
}

func TestCMEDUsecase_GetByID_NotFound(t *testing.T) {
	uc := NewCMEDUsecase(&mockCMEDRepo{cmed: nil}, nil)

	result, err := uc.GetByID(context.Background(), 999)
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestCMEDUsecase_GetByID_Error(t *testing.T) {
	uc := NewCMEDUsecase(&mockCMEDRepo{err: errors.New("db error")}, nil)

	_, err := uc.GetByID(context.Background(), 1)
	assert.Error(t, err)
}

func TestCMEDUsecase_GetByRegistro_Success(t *testing.T) {
	cmed := &entity.CMEDConformidade{COSeqID: 1, NUSanReg: int64Ptr(12345)}
	uc := NewCMEDUsecase(&mockCMEDRepo{cmed: cmed}, nil)

	result, err := uc.GetByRegistro(context.Background(), 12345, "")
	assert.NoError(t, err)
	assert.Equal(t, cmed, result)
}

func TestCMEDUsecase_GetByRegistro_WithDTReferencia(t *testing.T) {
	cmed := &entity.CMEDConformidade{COSeqID: 1, NUSanReg: int64Ptr(12345), DTReferencia: "2024-01-01"}
	uc := NewCMEDUsecase(&mockCMEDRepo{cmed: cmed}, nil)

	result, err := uc.GetByRegistro(context.Background(), 12345, "2024-01-01")
	assert.NoError(t, err)
	assert.Equal(t, cmed, result)
}

func TestCMEDUsecase_GetByEAN_Success(t *testing.T) {
	cmed := &entity.CMEDConformidade{COSeqID: 1, NUEAN1: strPtr("7891234567890")}
	uc := NewCMEDUsecase(&mockCMEDRepo{cmed: cmed}, nil)

	result, err := uc.GetByEAN(context.Background(), "7891234567890", "")
	assert.NoError(t, err)
	assert.Equal(t, cmed, result)
}

func TestCMEDUsecase_GetHistorico_Success(t *testing.T) {
	sanReg := int64(12345)
	cmed := &entity.CMEDConformidade{COSeqID: 1, NUSanReg: &sanReg}
	historico := []entity.CMEDConformidade{
		{COSeqID: 1, NUSanReg: &sanReg, DTReferencia: "2024-01-01"},
		{COSeqID: 2, NUSanReg: &sanReg, DTReferencia: "2024-02-01"},
	}
	uc := NewCMEDUsecase(&mockCMEDRepo{cmed: cmed, historico: historico}, nil)

	result, err := uc.GetHistorico(context.Background(), 1)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestCMEDUsecase_GetHistorico_NoSanReg(t *testing.T) {
	cmed := &entity.CMEDConformidade{COSeqID: 1, NUSanReg: nil}
	uc := NewCMEDUsecase(&mockCMEDRepo{cmed: cmed}, nil)

	result, err := uc.GetHistorico(context.Background(), 1)
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestCMEDUsecase_GetHistorico_NilCMED(t *testing.T) {
	uc := NewCMEDUsecase(&mockCMEDRepo{cmed: nil}, nil)

	result, err := uc.GetHistorico(context.Background(), 999)
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestCMEDUsecase_GetHistorico_Error(t *testing.T) {
	uc := NewCMEDUsecase(&mockCMEDRepo{err: errors.New("db error")}, nil)

	_, err := uc.GetHistorico(context.Background(), 1)
	assert.Error(t, err)
}

func strPtr(s string) *string {
	return &s
}

func int64Ptr(i int64) *int64 {
	return &i
}
