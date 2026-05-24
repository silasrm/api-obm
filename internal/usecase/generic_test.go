package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/silasrm/api-obm/internal/domain/entity"
	"github.com/silasrm/api-obm/internal/domain/repository"
	"github.com/stretchr/testify/assert"
)

type mockGenericRepo struct {
	result interface{}
	err    error
}

func (m *mockGenericRepo) GetByID(ctx context.Context, id int64) (interface{}, error) {
	return m.result, m.err
}

func (m *mockGenericRepo) List(ctx context.Context, filter repository.FilterParams) (interface{}, error) {
	return m.result, m.err
}

func TestGenericUsecase_GetVTM_Success(t *testing.T) {
	vtm := &entity.VTM{COSeqID: 1, NONm: "Test VTM"}
	repo := &mockGenericRepo{result: vtm}
	uc := NewGenericUsecase(repo, &mockGenericRepo{}, &mockGenericRepo{}, &mockGenericRepo{}, &mockGenericRepo{})

	result, err := uc.GetVTM(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, vtm, result)
}

func TestGenericUsecase_GetVTM_TypeAssertionFail(t *testing.T) {
	repo := &mockGenericRepo{result: "wrong type"}
	uc := NewGenericUsecase(repo, &mockGenericRepo{}, &mockGenericRepo{}, &mockGenericRepo{}, &mockGenericRepo{})

	result, err := uc.GetVTM(context.Background(), 1)
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestGenericUsecase_GetVTM_NotFound(t *testing.T) {
	repo := &mockGenericRepo{result: nil}
	uc := NewGenericUsecase(repo, &mockGenericRepo{}, &mockGenericRepo{}, &mockGenericRepo{}, &mockGenericRepo{})

	result, err := uc.GetVTM(context.Background(), 999)
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestGenericUsecase_ListVTMs_Success(t *testing.T) {
	page := &entity.CursorPage[entity.VTM]{Items: []entity.VTM{{COSeqID: 1}}, Total: 1}
	repo := &mockGenericRepo{result: page}
	uc := NewGenericUsecase(repo, &mockGenericRepo{}, &mockGenericRepo{}, &mockGenericRepo{}, &mockGenericRepo{})

	result, err := uc.ListVTMs(context.Background(), repository.FilterParams{})
	assert.NoError(t, err)
	assert.Equal(t, page, result)
}

func TestGenericUsecase_GetVMPP_Success(t *testing.T) {
	vmpp := &entity.VMPP{COSeqID: 1, NONm: "Test VMPP"}
	repo := &mockGenericRepo{result: vmpp}
	uc := NewGenericUsecase(&mockGenericRepo{}, repo, &mockGenericRepo{}, &mockGenericRepo{}, &mockGenericRepo{})

	result, err := uc.GetVMPP(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, vmpp, result)
}

func TestGenericUsecase_GetAMPP_Success(t *testing.T) {
	ampp := &entity.AMPP{COSeqID: 1, NONm: "Test AMPP"}
	repo := &mockGenericRepo{result: ampp}
	uc := NewGenericUsecase(&mockGenericRepo{}, &mockGenericRepo{}, repo, &mockGenericRepo{}, &mockGenericRepo{})

	result, err := uc.GetAMPP(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, ampp, result)
}

func TestGenericUsecase_GetDCB_Success(t *testing.T) {
	dcb := &entity.DCB{COSeqID: 1, DSDcb: "Test DCB"}
	repo := &mockGenericRepo{result: dcb}
	uc := NewGenericUsecase(&mockGenericRepo{}, &mockGenericRepo{}, &mockGenericRepo{}, repo, &mockGenericRepo{})

	result, err := uc.GetDCB(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, dcb, result)
}

func TestGenericUsecase_GetIngredient_Success(t *testing.T) {
	ing := &entity.IngredientSubstance{COSeqID: 1, NONm: "Test Ingredient"}
	repo := &mockGenericRepo{result: ing}
	uc := NewGenericUsecase(&mockGenericRepo{}, &mockGenericRepo{}, &mockGenericRepo{}, &mockGenericRepo{}, repo)

	result, err := uc.GetIngredient(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, ing, result)
}

func TestGenericUsecase_GetByID_RepoError(t *testing.T) {
	repo := &mockGenericRepo{err: errors.New("db error")}
	uc := NewGenericUsecase(repo, &mockGenericRepo{}, &mockGenericRepo{}, &mockGenericRepo{}, &mockGenericRepo{})

	_, err := uc.GetVTM(context.Background(), 1)
	assert.Error(t, err)
}
