package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/silasrm/api-obm/internal/domain/entity"
	"github.com/silasrm/api-obm/internal/domain/repository"
	"github.com/stretchr/testify/assert"
)

type mockSupplierRepo struct {
	result interface{}
	err    error
}

func (m *mockSupplierRepo) GetByID(ctx context.Context, id int64) (interface{}, error) {
	return m.result, m.err
}

func (m *mockSupplierRepo) List(ctx context.Context, filter repository.FilterParams) (interface{}, error) {
	return m.result, m.err
}

func TestSupplierUsecase_GetByID_Success(t *testing.T) {
	supplier := &entity.Supplier{COSeqID: 1, NODescr: "Test Supplier"}
	repo := &mockSupplierRepo{result: supplier}
	uc := NewSupplierUsecase(repo)

	result, err := uc.GetByID(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, supplier, result)
}

func TestSupplierUsecase_GetByID_NotFound(t *testing.T) {
	repo := &mockSupplierRepo{result: nil}
	uc := NewSupplierUsecase(repo)

	result, err := uc.GetByID(context.Background(), 999)
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestSupplierUsecase_GetByID_RepoError(t *testing.T) {
	repo := &mockSupplierRepo{err: errors.New("db error")}
	uc := NewSupplierUsecase(repo)

	_, err := uc.GetByID(context.Background(), 1)
	assert.Error(t, err)
}

func TestSupplierUsecase_GetByID_TypeAssertionFail(t *testing.T) {
	repo := &mockSupplierRepo{result: "wrong type"}
	uc := NewSupplierUsecase(repo)

	result, err := uc.GetByID(context.Background(), 1)
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestSupplierUsecase_List_Success(t *testing.T) {
	page := &entity.CursorPage[entity.Supplier]{
		Items: []entity.Supplier{{COSeqID: 1}},
		Total: 1,
	}
	repo := &mockSupplierRepo{result: page}
	uc := NewSupplierUsecase(repo)

	result, err := uc.List(context.Background(), repository.FilterParams{})
	assert.NoError(t, err)
	assert.Equal(t, page, result)
}

func TestSupplierUsecase_List_TypeAssertionFail(t *testing.T) {
	repo := &mockSupplierRepo{result: "wrong type"}
	uc := NewSupplierUsecase(repo)

	result, err := uc.List(context.Background(), repository.FilterParams{})
	assert.NoError(t, err)
	assert.Nil(t, result)
}
