package usecase

import (
	"context"

	"github.com/silasrm/api-obm/internal/domain/entity"
	"github.com/silasrm/api-obm/internal/domain/repository"
)

type DomainUsecase struct {
	domainRepo repository.DomainRepository
}

func NewDomainUsecase(domainRepo repository.DomainRepository) *DomainUsecase {
	return &DomainUsecase{domainRepo: domainRepo}
}

func (u *DomainUsecase) List(ctx context.Context, domainType string, filters repository.FilterParams) (*entity.CursorPage[entity.Domain], error) {
	return u.domainRepo.List(ctx, domainType, filters)
}

func (u *DomainUsecase) GetByID(ctx context.Context, domainType string, id int64) (*entity.Domain, error) {
	return u.domainRepo.GetByID(ctx, domainType, id)
}
