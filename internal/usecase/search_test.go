package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/silasrm/api-obm/internal/domain/entity"
	"github.com/stretchr/testify/assert"
)

type mockSearchRepo struct {
	hits      []entity.SearchHit
	total     int64
	cursor    string
	err       error
}

func (m *mockSearchRepo) Search(ctx context.Context, query string, entities []string, filters map[string]string, limit int, cursor string) ([]entity.SearchHit, int64, string, error) {
	if m.err != nil {
		return nil, 0, "", m.err
	}
	return m.hits, m.total, m.cursor, nil
}

func TestSearchUsecase_Search_Success(t *testing.T) {
	repo := &mockSearchRepo{
		hits:   []entity.SearchHit{{Entity: "vmp", ID: 1, Nome: "Test"}},
		total:  1,
		cursor: "next",
	}
	uc := NewSearchUsecase(repo)

	page, err := uc.Search(context.Background(), "test", []string{"vmp"}, nil, 10, "")
	assert.NoError(t, err)
	assert.Len(t, page.Items, 1)
	assert.Equal(t, int64(1), page.Total)
	assert.Equal(t, "next", page.Cursor)
	assert.Equal(t, 10, page.Limit)
}

func TestSearchUsecase_Search_DefaultLimit(t *testing.T) {
	repo := &mockSearchRepo{
		hits:  []entity.SearchHit{},
		total: 0,
	}
	uc := NewSearchUsecase(repo)

	page, err := uc.Search(context.Background(), "test", nil, nil, 0, "")
	assert.NoError(t, err)
	assert.Equal(t, 20, page.Limit)
}

func TestSearchUsecase_Search_NegativeLimit(t *testing.T) {
	repo := &mockSearchRepo{
		hits:  []entity.SearchHit{},
		total: 0,
	}
	uc := NewSearchUsecase(repo)

	page, err := uc.Search(context.Background(), "test", nil, nil, -5, "")
	assert.NoError(t, err)
	assert.Equal(t, 20, page.Limit)
}

func TestSearchUsecase_Search_RepoError(t *testing.T) {
	repo := &mockSearchRepo{err: errors.New("search error")}
	uc := NewSearchUsecase(repo)

	_, err := uc.Search(context.Background(), "test", nil, nil, 10, "")
	assert.Error(t, err)
}
