package repository

import (
	"context"

	"github.com/silasrm/api-obm/internal/domain/entity"
)

type FilterParams struct {
	Nome      string `json:"nome"`
	Descricao string `json:"descricao"`
	Codigo    string `json:"codigo"`
	Fabricante string `json:"fabricante"`
	Ativo     *bool  `json:"ativo"`
	Limit     int    `json:"limit"`
	Cursor    string `json:"cursor"`
}

type VTMRepository interface {
	GetByID(ctx context.Context, id int64) (interface{}, error)
	List(ctx context.Context, filter FilterParams) (interface{}, error)
}

type VMPRepository interface {
	GetByID(ctx context.Context, id int64) (*entity.VMP, error)
	GetDetailByID(ctx context.Context, id int64) (*entity.VMPDetail, error)
	List(ctx context.Context, filter FilterParams) (*entity.CursorPage[entity.VMP], error)
}

type VMPPRepository interface {
	GetByID(ctx context.Context, id int64) (interface{}, error)
	List(ctx context.Context, filter FilterParams) (interface{}, error)
}

type AMPRepository interface {
	GetByID(ctx context.Context, id int64) (*entity.AMP, error)
	GetDetailByID(ctx context.Context, id int64) (*entity.AMPDetail, error)
	List(ctx context.Context, filter FilterParams) (*entity.CursorPage[entity.AMP], error)
}

type AMPPRepository interface {
	GetByID(ctx context.Context, id int64) (interface{}, error)
	List(ctx context.Context, filter FilterParams) (interface{}, error)
}

type DCBRepository interface {
	GetByID(ctx context.Context, id int64) (interface{}, error)
	List(ctx context.Context, filter FilterParams) (interface{}, error)
}

type SupplierRepository interface {
	GetByID(ctx context.Context, id int64) (interface{}, error)
	List(ctx context.Context, filter FilterParams) (interface{}, error)
}

type IngredientSubstanceRepository interface {
	GetByID(ctx context.Context, id int64) (interface{}, error)
	List(ctx context.Context, filter FilterParams) (interface{}, error)
}

type DomainRepository interface {
	GetByID(ctx context.Context, domainType string, id int64) (*entity.Domain, error)
	List(ctx context.Context, domainType string, filter FilterParams) (*entity.CursorPage[entity.Domain], error)
}

type SearchRepository interface {
	Search(ctx context.Context, query string, entities []string, filters map[string]string, limit int, cursor string) ([]entity.SearchHit, int64, string, error)
}

type SyncRepository interface {
	GetAllVMPs(ctx context.Context) ([]map[string]interface{}, error)
	GetAllAMPs(ctx context.Context) ([]map[string]interface{}, error)
	GetAllSuppliers(ctx context.Context) ([]map[string]interface{}, error)
}

type UserRepository interface {
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
	Create(ctx context.Context, user *entity.User) error
}
