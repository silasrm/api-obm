package meilisearch

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/meilisearch/meilisearch-go"
	"github.com/silasrm/api-obm/internal/domain/entity"
	"github.com/silasrm/api-obm/internal/infrastructure/config"
)

type MeilisearchRepo struct {
	client      meilisearch.ServiceManager
	indexPrefix string
}

func NewMeilisearchRepo(cfg config.MeilisearchConfig) *MeilisearchRepo {
	client := meilisearch.New(cfg.URL, meilisearch.WithAPIKey(cfg.APIKey))
	return &MeilisearchRepo{
		client:      client,
		indexPrefix: cfg.IndexPrefix,
	}
}

func (r *MeilisearchRepo) Search(ctx context.Context, query string, entities []string, filters map[string]string, limit int, cursor string) ([]entity.SearchHit, int64, string, error) {
	if len(entities) == 0 {
		entities = []string{"vmp", "amp", "supplier"}
	}

	offset := 0
	if cursor != "" {
		decoded, err := base64.StdEncoding.DecodeString(cursor)
		if err == nil {
			offset, _ = strconv.Atoi(string(decoded))
		}
	}

	if limit <= 0 {
		limit = 20
	}

	var allHits []entity.SearchHit
	var totalHits int64

	for _, ent := range entities {
		indexUID := r.indexPrefix + ent
		idx := r.client.Index(indexUID)

		var filterExprs []string
		for k, v := range filters {
			filterExprs = append(filterExprs, fmt.Sprintf("%s = %q", k, v))
		}

		req := &meilisearch.SearchRequest{
			Query:  query,
			Limit:  int64(limit),
			Offset: int64(offset),
		}
		if len(filterExprs) > 0 {
			req.Filter = &filterExprs
		}

		resp, err := idx.SearchWithContext(ctx, query, req)
		if err != nil {
			continue
		}

		totalHits += resp.TotalHits

		for _, hit := range resp.Hits {
			sh := entity.SearchHit{
				Entity: ent,
			}
			if idRaw, ok := hit["co_seq_id"]; ok {
				sh.ID = jsonInt(idRaw)
			}
			if nomeRaw, ok := hit["no_nm"]; ok {
				sh.Nome = jsonStr(nomeRaw)
			} else if nomeRaw, ok := hit["no_descr"]; ok {
				sh.Nome = jsonStr(nomeRaw)
			}
			if codRaw, ok := hit["nu_vpid"]; ok {
				sh.Codigo = jsonStr(codRaw)
			} else if codRaw, ok := hit["nu_apid"]; ok {
				sh.Codigo = jsonStr(codRaw)
			} else if codRaw, ok := hit["nu_cd"]; ok {
				sh.Codigo = jsonStr(codRaw)
			}
			if fabRaw, ok := hit["supplier_name"]; ok {
				sh.Fabricante = jsonStr(fabRaw)
			}
			if descRaw, ok := hit["ds_descr"]; ok {
				sh.Descricao = jsonStr(descRaw)
			}

			allHits = append(allHits, sh)
		}
	}

	newOffset := offset + limit
	nextCursor := base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(newOffset)))

	if int64(offset) >= totalHits {
		nextCursor = ""
	}

	return allHits, totalHits, nextCursor, nil
}

func jsonStr(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}
	var s string
	if raw[0] == '"' {
		s, _ = strconv.Unquote(string(raw))
	} else {
		s = string(raw)
	}
	return s
}

func jsonInt(raw json.RawMessage) int64 {
	if len(raw) == 0 {
		return 0
	}
	var n int64
	fmt.Sscanf(string(raw), "%d", &n)
	return n
}
