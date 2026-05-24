package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/silasrm/api-obm/internal/domain/entity"
	"github.com/silasrm/api-obm/internal/domain/repository"
	"github.com/stretchr/testify/assert"
)

type mockAMPRepo struct {
	amp    *entity.AMP
	detail *entity.AMPDetail
	page   *entity.CursorPage[entity.AMP]
	err    error
}

func (m *mockAMPRepo) GetByID(ctx context.Context, id int64) (*entity.AMP, error) {
	return m.amp, m.err
}

func (m *mockAMPRepo) GetDetailByID(ctx context.Context, id int64) (*entity.AMPDetail, error) {
	return m.detail, m.err
}

func (m *mockAMPRepo) List(ctx context.Context, filter repository.FilterParams) (*entity.CursorPage[entity.AMP], error) {
	return m.page, m.err
}

func TestAMPUsecase_GetByID(t *testing.T) {
	amp := &entity.AMP{COSeqID: 1, NONm: "Test AMP"}
	uc := NewAMPUsecase(&mockAMPRepo{amp: amp})

	result, err := uc.GetByID(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, amp, result)
}

func TestAMPUsecase_GetByID_Error(t *testing.T) {
	uc := NewAMPUsecase(&mockAMPRepo{err: errors.New("db error")})

	_, err := uc.GetByID(context.Background(), 1)
	assert.Error(t, err)
}

func TestAMPUsecase_GetDetailByID(t *testing.T) {
	detail := &entity.AMPDetail{AMP: entity.AMP{COSeqID: 1}}
	uc := NewAMPUsecase(&mockAMPRepo{detail: detail})

	result, err := uc.GetDetailByID(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, detail, result)
}

func TestAMPUsecase_List(t *testing.T) {
	page := &entity.CursorPage[entity.AMP]{Items: []entity.AMP{{COSeqID: 1}}, Total: 1}
	uc := NewAMPUsecase(&mockAMPRepo{page: page})

	result, err := uc.List(context.Background(), repository.FilterParams{})
	assert.NoError(t, err)
	assert.Equal(t, page, result)
}
