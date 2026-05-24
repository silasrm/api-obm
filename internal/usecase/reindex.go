package usecase

import (
	"context"
	"fmt"
	"log"

	"github.com/silasrm/api-obm/internal/domain/repository"
	meilisearchrepo "github.com/silasrm/api-obm/internal/infrastructure/persistence/meilisearch"
)

type ReindexUsecase struct {
	syncRepo  repository.SyncRepository
	meiliRepo *meilisearchrepo.MeilisearchRepo
}

func NewReindexUsecase(syncRepo repository.SyncRepository, meiliRepo *meilisearchrepo.MeilisearchRepo) *ReindexUsecase {
	return &ReindexUsecase{
		syncRepo:  syncRepo,
		meiliRepo: meiliRepo,
	}
}

func (u *ReindexUsecase) Reindex(ctx context.Context) (map[string]int64, error) {
	if err := u.meiliRepo.ConfigureIndexes(ctx); err != nil {
		return nil, fmt.Errorf("configure indexes: %w", err)
	}

	indexed := make(map[string]int64)

	vmpDocs, err := u.syncRepo.GetAllVMPs(ctx)
	if err != nil {
		return nil, fmt.Errorf("get vmps: %w", err)
	}
	if err := u.meiliRepo.IndexVMPs(ctx, vmpDocs); err != nil {
		return nil, fmt.Errorf("index vmps: %w", err)
	}
	indexed["vmp"] = int64(len(vmpDocs))
	log.Printf("Indexed %d VMPs", len(vmpDocs))

	ampDocs, err := u.syncRepo.GetAllAMPs(ctx)
	if err != nil {
		return nil, fmt.Errorf("get amps: %w", err)
	}
	if err := u.meiliRepo.IndexAMPs(ctx, ampDocs); err != nil {
		return nil, fmt.Errorf("index amps: %w", err)
	}
	indexed["amp"] = int64(len(ampDocs))
	log.Printf("Indexed %d AMPs", len(ampDocs))

	suppDocs, err := u.syncRepo.GetAllSuppliers(ctx)
	if err != nil {
		return nil, fmt.Errorf("get suppliers: %w", err)
	}
	if err := u.meiliRepo.IndexSuppliers(ctx, suppDocs); err != nil {
		return nil, fmt.Errorf("index suppliers: %w", err)
	}
	indexed["supplier"] = int64(len(suppDocs))
	log.Printf("Indexed %d Suppliers", len(suppDocs))

	return indexed, nil
}
