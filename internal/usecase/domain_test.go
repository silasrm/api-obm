package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/silasrm/api-obm/internal/domain/entity"
	"github.com/silasrm/api-obm/internal/domain/repository"
	"github.com/stretchr/testify/assert"
)

type mockDomainRepo struct {
	domain *entity.Domain
	page   *entity.CursorPage[entity.Domain]
	err    error
}

func (m *mockDomainRepo) GetByID(ctx context.Context, domainType string, id int64) (*entity.Domain, error) {
	return m.domain, m.err
}

func (m *mockDomainRepo) List(ctx context.Context, domainType string, filter repository.FilterParams) (*entity.CursorPage[entity.Domain], error) {
	return m.page, m.err
}

func TestDomainUsecase_GetByID(t *testing.T) {
	domain := &entity.Domain{COSeqID: 1, NODescr: "Test Domain"}
	uc := NewDomainUsecase(&mockDomainRepo{domain: domain})

	result, err := uc.GetByID(context.Background(), "form", 1)
	assert.NoError(t, err)
	assert.Equal(t, domain, result)
}

func TestDomainUsecase_GetByID_Error(t *testing.T) {
	uc := NewDomainUsecase(&mockDomainRepo{err: errors.New("db error")})

	_, err := uc.GetByID(context.Background(), "form", 1)
	assert.Error(t, err)
}

func TestDomainUsecase_List(t *testing.T) {
	page := &entity.CursorPage[entity.Domain]{Items: []entity.Domain{{COSeqID: 1}}, Total: 1}
	uc := NewDomainUsecase(&mockDomainRepo{page: page})

	result, err := uc.List(context.Background(), "form", repository.FilterParams{})
	assert.NoError(t, err)
	assert.Equal(t, page, result)
}
