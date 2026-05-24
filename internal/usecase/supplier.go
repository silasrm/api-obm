package usecase

import (
	"context"

	"github.com/silasrm/api-obm/internal/domain/entity"
	"github.com/silasrm/api-obm/internal/domain/repository"
)

type SupplierUsecase struct {
	supplierRepo repository.SupplierRepository
}

func NewSupplierUsecase(supplierRepo repository.SupplierRepository) *SupplierUsecase {
	return &SupplierUsecase{supplierRepo: supplierRepo}
}

func (u *SupplierUsecase) GetByID(ctx context.Context, id int64) (*entity.Supplier, error) {
	result, err := u.supplierRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	s, ok := result.(*entity.Supplier)
	if !ok {
		return nil, nil
	}
	return s, nil
}

func (u *SupplierUsecase) List(ctx context.Context, filters repository.FilterParams) (*entity.CursorPage[entity.Supplier], error) {
	result, err := u.supplierRepo.List(ctx, filters)
	if err != nil {
		return nil, err
	}
	page, ok := result.(*entity.CursorPage[entity.Supplier])
	if !ok {
		return nil, nil
	}
	return page, nil
}
