package usecase

import (
	"context"

	"github.com/silasrm/api-obm/internal/domain/entity"
	"github.com/silasrm/api-obm/internal/domain/repository"
)

type SearchUsecase struct {
	searchRepo repository.SearchRepository
}

func NewSearchUsecase(searchRepo repository.SearchRepository) *SearchUsecase {
	return &SearchUsecase{searchRepo: searchRepo}
}

func (u *SearchUsecase) Search(ctx context.Context, query string, entities []string, filters map[string]string, limit int, cursor string) (*entity.CursorPage[entity.SearchHit], error) {
	hits, total, nextCursor, err := u.searchRepo.Search(ctx, query, entities, filters, limit, cursor)
	if err != nil {
		return nil, err
	}

	if limit <= 0 {
		limit = 20
	}

	return &entity.CursorPage[entity.SearchHit]{
		Items:  hits,
		Cursor: nextCursor,
		Limit:  limit,
		Total:  total,
	}, nil
}
