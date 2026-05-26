package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/silasrm/api-obm/internal/domain/entity"
	"github.com/silasrm/api-obm/internal/domain/repository"
	redisrepo "github.com/silasrm/api-obm/internal/infrastructure/persistence/redis"
)

type SupplierCMEDResponse struct {
	Supplier *entity.Supplier         `json:"supplier"`
	CMED     []entity.CMEDConformidade `json:"cmed,omitempty"`
}

type SupplierCMEDUsecase struct {
	supplierRepo repository.SupplierRepository
	cmedRepo     repository.CMEDRepository
	cacheRepo    *redisrepo.CacheRepo
}

func NewSupplierCMEDUsecase(
	supplierRepo repository.SupplierRepository,
	cmedRepo repository.CMEDRepository,
	cacheRepo *redisrepo.CacheRepo,
) *SupplierCMEDUsecase {
	return &SupplierCMEDUsecase{
		supplierRepo: supplierRepo,
		cmedRepo:     cmedRepo,
		cacheRepo:    cacheRepo,
	}
}

func (u *SupplierCMEDUsecase) GetSupplierWithCMED(ctx context.Context, supplierID int64, dtReferencia string) (*SupplierCMEDResponse, error) {
	supplierIface, err := u.supplierRepo.GetByID(ctx, supplierID)
	if err != nil {
		return nil, fmt.Errorf("Supplier not found: %w", err)
	}
	if supplierIface == nil {
		return nil, fmt.Errorf("Supplier not found")
	}

	supplier, ok := supplierIface.(*entity.Supplier)
	if !ok {
		return nil, fmt.Errorf("Supplier not found")
	}

	resp := &SupplierCMEDResponse{Supplier: supplier}

	if supplier.NUCnpj.Valid && supplier.NUCnpj.String != "" {
		normalizedCNPJ := repository.NormalizeCNPJ(supplier.NUCnpj.String)
		cacheKey := fmt.Sprintf("supplier_cmed:%d:%s", supplierID, dtReferencia)

		if u.cacheRepo != nil {
			cached, err := u.cacheRepo.Get(ctx, cacheKey)
			if err == nil && cached != "" {
				var cmedList []entity.CMEDConformidade
				if json.Unmarshal([]byte(cached), &cmedList) == nil {
					resp.CMED = cmedList
					return resp, nil
				}
			}
		}

		cmedList, err := u.cmedRepo.GetByCNPJ(ctx, normalizedCNPJ, dtReferencia)
		if err != nil {
			log.Printf("Warning: failed to get CMED for supplier %d: %v", supplierID, err)
			resp.CMED = []entity.CMEDConformidade{}
			return resp, nil
		}

		resp.CMED = cmedList

		if u.cacheRepo != nil {
			if data, err := json.Marshal(cmedList); err == nil {
				if err := u.cacheRepo.Set(ctx, cacheKey, string(data), 0); err != nil {
					log.Printf("Warning: failed to cache CMED for supplier %d: %v", supplierID, err)
				}
			}
		}
	}

	return resp, nil
}
