package usecase

import (
	"context"

	"github.com/silasrm/api-obm/internal/domain/entity"
	"github.com/silasrm/api-obm/internal/domain/repository"
)

type AMPUsecase struct {
	ampRepo repository.AMPRepository
}

func NewAMPUsecase(ampRepo repository.AMPRepository) *AMPUsecase {
	return &AMPUsecase{ampRepo: ampRepo}
}

func (u *AMPUsecase) GetByID(ctx context.Context, id int64) (*entity.AMP, error) {
	return u.ampRepo.GetByID(ctx, id)
}

func (u *AMPUsecase) GetDetailByID(ctx context.Context, id int64) (*entity.AMPDetail, error) {
	return u.ampRepo.GetDetailByID(ctx, id)
}

func (u *AMPUsecase) List(ctx context.Context, filters repository.FilterParams) (*entity.CursorPage[entity.AMP], error) {
	return u.ampRepo.List(ctx, filters)
}
