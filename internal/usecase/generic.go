package usecase

import (
	"context"

	"github.com/silasrm/api-obm/internal/domain/entity"
	"github.com/silasrm/api-obm/internal/domain/repository"
)

type GenericUsecase struct {
	vtmRepo        repository.VTMRepository
	vmppRepo       repository.VMPPRepository
	amppRepo       repository.AMPPRepository
	dcbRepo        repository.DCBRepository
	ingredientRepo repository.IngredientSubstanceRepository
}

func NewGenericUsecase(
	vtmRepo repository.VTMRepository,
	vmppRepo repository.VMPPRepository,
	amppRepo repository.AMPPRepository,
	dcbRepo repository.DCBRepository,
	ingredientRepo repository.IngredientSubstanceRepository,
) *GenericUsecase {
	return &GenericUsecase{
		vtmRepo:        vtmRepo,
		vmppRepo:       vmppRepo,
		amppRepo:       amppRepo,
		dcbRepo:        dcbRepo,
		ingredientRepo: ingredientRepo,
	}
}

func (u *GenericUsecase) GetVTM(ctx context.Context, id int64) (*entity.VTM, error) {
	result, err := u.vtmRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	v, ok := result.(*entity.VTM)
	if !ok {
		return nil, nil
	}
	return v, nil
}

func (u *GenericUsecase) ListVTMs(ctx context.Context, filter repository.FilterParams) (*entity.CursorPage[entity.VTM], error) {
	result, err := u.vtmRepo.List(ctx, filter)
	if err != nil {
		return nil, err
	}
	page, ok := result.(*entity.CursorPage[entity.VTM])
	if !ok {
		return nil, nil
	}
	return page, nil
}

func (u *GenericUsecase) GetVMPP(ctx context.Context, id int64) (*entity.VMPP, error) {
	result, err := u.vmppRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	v, ok := result.(*entity.VMPP)
	if !ok {
		return nil, nil
	}
	return v, nil
}

func (u *GenericUsecase) ListVMPPs(ctx context.Context, filter repository.FilterParams) (*entity.CursorPage[entity.VMPP], error) {
	result, err := u.vmppRepo.List(ctx, filter)
	if err != nil {
		return nil, err
	}
	page, ok := result.(*entity.CursorPage[entity.VMPP])
	if !ok {
		return nil, nil
	}
	return page, nil
}

func (u *GenericUsecase) GetAMPP(ctx context.Context, id int64) (*entity.AMPP, error) {
	result, err := u.amppRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	v, ok := result.(*entity.AMPP)
	if !ok {
		return nil, nil
	}
	return v, nil
}

func (u *GenericUsecase) ListAMPPs(ctx context.Context, filter repository.FilterParams) (*entity.CursorPage[entity.AMPP], error) {
	result, err := u.amppRepo.List(ctx, filter)
	if err != nil {
		return nil, err
	}
	page, ok := result.(*entity.CursorPage[entity.AMPP])
	if !ok {
		return nil, nil
	}
	return page, nil
}

func (u *GenericUsecase) GetDCB(ctx context.Context, id int64) (*entity.DCB, error) {
	result, err := u.dcbRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	v, ok := result.(*entity.DCB)
	if !ok {
		return nil, nil
	}
	return v, nil
}

func (u *GenericUsecase) ListDCBs(ctx context.Context, filter repository.FilterParams) (*entity.CursorPage[entity.DCB], error) {
	result, err := u.dcbRepo.List(ctx, filter)
	if err != nil {
		return nil, err
	}
	page, ok := result.(*entity.CursorPage[entity.DCB])
	if !ok {
		return nil, nil
	}
	return page, nil
}

func (u *GenericUsecase) GetIngredient(ctx context.Context, id int64) (*entity.IngredientSubstance, error) {
	result, err := u.ingredientRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	v, ok := result.(*entity.IngredientSubstance)
	if !ok {
		return nil, nil
	}
	return v, nil
}

func (u *GenericUsecase) ListIngredients(ctx context.Context, filter repository.FilterParams) (*entity.CursorPage[entity.IngredientSubstance], error) {
	result, err := u.ingredientRepo.List(ctx, filter)
	if err != nil {
		return nil, err
	}
	page, ok := result.(*entity.CursorPage[entity.IngredientSubstance])
	if !ok {
		return nil, nil
	}
	return page, nil
}
