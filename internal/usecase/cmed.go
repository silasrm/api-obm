package usecase

import (
	"context"

	"github.com/silasrm/api-obm/internal/domain/entity"
	"github.com/silasrm/api-obm/internal/domain/repository"
	redisrepo "github.com/silasrm/api-obm/internal/infrastructure/persistence/redis"
)

type CMEDUsecase struct {
	cmedRepo  repository.CMEDRepository
	cacheRepo *redisrepo.CacheRepo
}

func NewCMEDUsecase(cmedRepo repository.CMEDRepository, cacheRepo *redisrepo.CacheRepo) *CMEDUsecase {
	return &CMEDUsecase{cmedRepo: cmedRepo, cacheRepo: cacheRepo}
}

func (u *CMEDUsecase) List(ctx context.Context, filter repository.CMEDFilterParams) (*entity.CursorPage[entity.CMEDConformidade], error) {
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	return u.cmedRepo.List(ctx, filter)
}

func (u *CMEDUsecase) GetByID(ctx context.Context, id int64) (*entity.CMEDConformidade, error) {
	return u.cmedRepo.GetByID(ctx, id)
}

func (u *CMEDUsecase) GetByRegistro(ctx context.Context, nuSanReg int64, dtReferencia string) (*entity.CMEDConformidade, error) {
	return u.cmedRepo.GetByNuSanReg(ctx, nuSanReg, dtReferencia)
}

func (u *CMEDUsecase) GetByEAN(ctx context.Context, ean string, dtReferencia string) (*entity.CMEDConformidade, error) {
	return u.cmedRepo.GetByEAN(ctx, ean, dtReferencia)
}

func (u *CMEDUsecase) GetHistorico(ctx context.Context, id int64) ([]entity.CMEDConformidade, error) {
	cmed, err := u.cmedRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if cmed == nil || cmed.NUSanReg == nil {
		return nil, nil
	}
	return u.cmedRepo.GetHistorico(ctx, *cmed.NUSanReg)
}
