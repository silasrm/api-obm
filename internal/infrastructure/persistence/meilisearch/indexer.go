package meilisearch

import (
	"context"

	"github.com/meilisearch/meilisearch-go"
)

func (r *MeilisearchRepo) ConfigureIndexes(ctx context.Context) error {
	indexes := map[string]struct {
		Searchable []string
		Filterable []interface{}
		Sortable   []string
		Ranking    []string
	}{
		"vmp": {
			Searchable: []string{"no_nm", "no_abbrevnm", "no_nmprev", "nu_vpid", "nu_snomed"},
			Filterable: []interface{}{"co_vtmid", "co_pres_statcd", "co_non_availcd", "co_df_indcd", "co_brimunocd", "st_registro_ativo"},
			Sortable:   []string{"no_nm"},
			Ranking:    []string{"words", "typo", "proximity", "attribute", "sort", "exactness"},
		},
		"amp": {
			Searchable: []string{"no_nm", "ds_descr", "no_abbrevnm", "nu_apid", "supplier_name"},
			Filterable: []interface{}{"co_vpid", "co_suppcd", "co_flavourcd", "co_lic_authcd", "co_avail_restrictcd", "co_medclscd", "st_registro_ativo"},
			Sortable:   []string{"no_nm", "supplier_name"},
			Ranking:    []string{"words", "typo", "proximity", "attribute", "sort", "exactness"},
		},
		"supplier": {
			Searchable: []string{"no_descr", "nu_cd", "nu_cnpj"},
			Filterable: []interface{}{"co_countrycd", "st_registro_ativo"},
			Sortable:   []string{"no_descr"},
			Ranking:    []string{"words", "typo", "proximity", "attribute", "sort", "exactness"},
		},
	}

	for suffix, cfg := range indexes {
		uid := r.indexPrefix + suffix
		r.client.CreateIndex(&meilisearch.IndexConfig{Uid: uid})

		idx := r.client.Index(uid)

		idx.UpdateSearchableAttributes(&cfg.Searchable)
		idx.UpdateFilterableAttributes(&cfg.Filterable)
		idx.UpdateSortableAttributes(&cfg.Sortable)
		idx.UpdateRankingRules(&cfg.Ranking)
	}

	return nil
}

func (r *MeilisearchRepo) IndexVMPs(ctx context.Context, docs []map[string]interface{}) error {
	return r.batchIndex("vmp", docs)
}

func (r *MeilisearchRepo) IndexAMPs(ctx context.Context, docs []map[string]interface{}) error {
	return r.batchIndex("amp", docs)
}

func (r *MeilisearchRepo) IndexSuppliers(ctx context.Context, docs []map[string]interface{}) error {
	return r.batchIndex("supplier", docs)
}

func (r *MeilisearchRepo) batchIndex(suffix string, docs []map[string]interface{}) error {
	uid := r.indexPrefix + suffix
	idx := r.client.Index(uid)

	batchSize := 1000
	for i := 0; i < len(docs); i += batchSize {
		end := i + batchSize
		if end > len(docs) {
			end = len(docs)
		}
		batch := docs[i:end]
		_, err := idx.UpdateDocuments(batch, nil)
		if err != nil {
			return err
		}
	}

	return nil
}
