package usecase

import (
	"context"

	"github.com/silasrm/api-obm/internal/domain/entity"
	"github.com/silasrm/api-obm/internal/domain/repository"
)

type VMPUsecase struct {
	vmpRepo repository.VMPRepository
}

func NewVMPUsecase(vmpRepo repository.VMPRepository) *VMPUsecase {
	return &VMPUsecase{vmpRepo: vmpRepo}
}

func (u *VMPUsecase) GetByID(ctx context.Context, id int64) (*entity.VMP, error) {
	return u.vmpRepo.GetByID(ctx, id)
}

func (u *VMPUsecase) GetDetailByID(ctx context.Context, id int64) (*entity.VMPDetail, error) {
	return u.vmpRepo.GetDetailByID(ctx, id)
}

func (u *VMPUsecase) List(ctx context.Context, filters repository.FilterParams) (*entity.CursorPage[entity.VMP], error) {
	return u.vmpRepo.List(ctx, filters)
}
