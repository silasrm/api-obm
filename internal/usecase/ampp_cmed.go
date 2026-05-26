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

type AMPPCMEDResponse struct {
	AMPP *entity.AMPP              `json:"ampp"`
	AMP  interface{}               `json:"amp,omitempty"`
	VMP  interface{}               `json:"vmp,omitempty"`
	CMED *entity.CMEDConformidade  `json:"cmed,omitempty"`
}

type AMPPCMEDUsecase struct {
	amppRepo  repository.AMPPRepository
	vmpRepo   repository.VMPRepository
	ampRepo   repository.AMPRepository
	cmedRepo  repository.CMEDRepository
	cacheRepo *redisrepo.CacheRepo
}

func NewAMPPCMEDUsecase(
	amppRepo repository.AMPPRepository,
	vmpRepo repository.VMPRepository,
	ampRepo repository.AMPRepository,
	cmedRepo repository.CMEDRepository,
	cacheRepo *redisrepo.CacheRepo,
) *AMPPCMEDUsecase {
	return &AMPPCMEDUsecase{
		amppRepo:  amppRepo,
		vmpRepo:   vmpRepo,
		ampRepo:   ampRepo,
		cmedRepo:  cmedRepo,
		cacheRepo: cacheRepo,
	}
}

func (u *AMPPCMEDUsecase) GetAMPPWithCMED(ctx context.Context, amppID int64, dtReferencia string) (*AMPPCMEDResponse, error) {
	amppIface, err := u.amppRepo.GetByID(ctx, amppID)
	if err != nil {
		return nil, fmt.Errorf("get ampp: %w", err)
	}
	if amppIface == nil {
		return nil, nil
	}

	ampp, ok := amppIface.(*entity.AMPP)
	if !ok {
		return nil, fmt.Errorf("invalid ampp type")
	}

	resp := &AMPPCMEDResponse{AMPP: ampp}

	if ampp.COApID > 0 {
		amp, err := u.ampRepo.GetByID(ctx, ampp.COApID)
		if err == nil && amp != nil {
			resp.AMP = amp
			if amp.COVpID > 0 {
				vmp, err := u.vmpRepo.GetByID(ctx, amp.COVpID)
				if err == nil && vmp != nil {
					resp.VMP = vmp
				}
			}
		}
	}

	if ampp.NUSanReg.Valid && ampp.NUSanReg.Int64 != 0 {
		cacheKey := fmt.Sprintf("ampp_cmed:%d:%s", amppID, dtReferencia)

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

		cmed, err := u.cmedRepo.GetByNuSanReg(ctx, ampp.NUSanReg.Int64, dtReferencia)
		if err != nil {
			log.Printf("Warning: failed to get CMED for ampp %d: %v", amppID, err)
		} else if cmed != nil {
			resp.CMED = cmed

			if u.cacheRepo != nil {
				if data, err := json.Marshal(cmed); err == nil {
					if err := u.cacheRepo.Set(ctx, cacheKey, string(data), 0); err != nil {
						log.Printf("Warning: failed to cache CMED for ampp %d: %v", amppID, err)
					}
				}
			}
		}
	}

	return resp, nil
}
