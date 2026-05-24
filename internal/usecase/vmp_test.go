package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/silasrm/api-obm/internal/domain/entity"
	"github.com/silasrm/api-obm/internal/domain/repository"
	"github.com/stretchr/testify/assert"
)

type mockVMPRepo struct {
	vmp    *entity.VMP
	detail *entity.VMPDetail
	page   *entity.CursorPage[entity.VMP]
	err    error
}

func (m *mockVMPRepo) GetByID(ctx context.Context, id int64) (*entity.VMP, error) {
	return m.vmp, m.err
}

func (m *mockVMPRepo) GetDetailByID(ctx context.Context, id int64) (*entity.VMPDetail, error) {
	return m.detail, m.err
}

func (m *mockVMPRepo) List(ctx context.Context, filter repository.FilterParams) (*entity.CursorPage[entity.VMP], error) {
	return m.page, m.err
}

func TestVMPUsecase_GetByID(t *testing.T) {
	vmp := &entity.VMP{COSeqID: 1, NONm: "Test VMP"}
	uc := NewVMPUsecase(&mockVMPRepo{vmp: vmp})

	result, err := uc.GetByID(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, vmp, result)
}

func TestVMPUsecase_GetByID_NotFound(t *testing.T) {
	uc := NewVMPUsecase(&mockVMPRepo{vmp: nil})

	result, err := uc.GetByID(context.Background(), 999)
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestVMPUsecase_GetByID_Error(t *testing.T) {
	uc := NewVMPUsecase(&mockVMPRepo{err: errors.New("db error")})

	_, err := uc.GetByID(context.Background(), 1)
	assert.Error(t, err)
}

func TestVMPUsecase_GetDetailByID(t *testing.T) {
	detail := &entity.VMPDetail{VMP: entity.VMP{COSeqID: 1}}
	uc := NewVMPUsecase(&mockVMPRepo{detail: detail})

	result, err := uc.GetDetailByID(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, detail, result)
}

func TestVMPUsecase_List(t *testing.T) {
	page := &entity.CursorPage[entity.VMP]{Items: []entity.VMP{{COSeqID: 1}}, Total: 1}
	uc := NewVMPUsecase(&mockVMPRepo{page: page})

	result, err := uc.List(context.Background(), repository.FilterParams{})
	assert.NoError(t, err)
	assert.Equal(t, page, result)
}
