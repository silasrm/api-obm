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

type AMPCMEDResponse struct {
	AMP  *entity.AMP             `json:"amp"`
	VMP  *entity.VMP             `json:"vmp,omitempty"`
	CMED *entity.CMEDConformidade `json:"cmed,omitempty"`
}

type AMPCMEDUsecase struct {
	ampRepo   repository.AMPRepository
	vmpRepo   repository.VMPRepository
	cmedRepo  repository.CMEDRepository
	cacheRepo *redisrepo.CacheRepo
}

func NewAMPCMEDUsecase(
	ampRepo repository.AMPRepository,
	vmpRepo repository.VMPRepository,
	cmedRepo repository.CMEDRepository,
	cacheRepo *redisrepo.CacheRepo,
) *AMPCMEDUsecase {
	return &AMPCMEDUsecase{
		ampRepo:   ampRepo,
		vmpRepo:   vmpRepo,
		cmedRepo:  cmedRepo,
		cacheRepo: cacheRepo,
	}
}

func (u *AMPCMEDUsecase) GetAMPWithCMED(ctx context.Context, ampID int64, dtReferencia string) (*AMPCMEDResponse, error) {
	amp, err := u.ampRepo.GetByID(ctx, ampID)
	if err != nil {
		return nil, fmt.Errorf("AMP not found: %w", err)
	}
	if amp == nil {
		return nil, fmt.Errorf("AMP not found")
	}

	resp := &AMPCMEDResponse{AMP: amp}

	if amp.COVpID > 0 {
		vmp, err := u.vmpRepo.GetByID(ctx, amp.COVpID)
		if err == nil && vmp != nil {
			resp.VMP = vmp
		} else if err != nil {
			log.Printf("Warning: failed to get VMP for amp %d: %v", ampID, err)
		}
	}

	if amp.NUNReg.Valid && amp.NUNReg.Int64 != 0 {
		cacheKey := fmt.Sprintf("amp_cmed:%d:%s", ampID, dtReferencia)

		if u.cacheRepo != nil {
			cached, err := u.cacheRepo.Get(ctx, cacheKey)
			if err == nil && cached != "" {
				var cmed entity.CMEDConformidade
				if json.Unmarshal([]byte(cached), &cmed) == nil {
					resp.CMED = &cmed
					return resp, nil
				}
			}
		}

		cmed, err := u.cmedRepo.GetByNuSanReg(ctx, amp.NUNReg.Int64, dtReferencia)
		if err != nil {
			log.Printf("Warning: failed to get CMED for amp %d: %v", ampID, err)
		} else if cmed != nil {
			resp.CMED = cmed

			if u.cacheRepo != nil {
				if data, err := json.Marshal(cmed); err == nil {
					if err := u.cacheRepo.Set(ctx, cacheKey, string(data), 0); err != nil {
						log.Printf("Warning: failed to cache CMED for amp %d: %v", ampID, err)
					}
				}
			}
		}
	}

	return resp, nil
}
