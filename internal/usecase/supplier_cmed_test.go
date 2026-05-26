package usecase

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/silasrm/api-obm/internal/domain/entity"
	"github.com/silasrm/api-obm/internal/domain/repository"
	"github.com/stretchr/testify/assert"
)

type mockSupplierRepoForCMED struct {
	supplier interface{}
	err      error
}

func (m *mockSupplierRepoForCMED) GetByID(ctx context.Context, id int64) (interface{}, error) {
	return m.supplier, m.err
}

func (m *mockSupplierRepoForCMED) List(ctx context.Context, filter repository.FilterParams) (interface{}, error) {
	return nil, nil
}

type mockCMEDRepoForSupplier struct {
	cmedList []entity.CMEDConformidade
	err      error
}

func (m *mockCMEDRepoForSupplier) GetByID(ctx context.Context, id int64) (*entity.CMEDConformidade, error) {
	return nil, nil
}

func (m *mockCMEDRepoForSupplier) GetByNuSanReg(ctx context.Context, nuSanReg int64, dtReferencia string) (*entity.CMEDConformidade, error) {
	return nil, nil
}

func (m *mockCMEDRepoForSupplier) GetByEAN(ctx context.Context, ean string, dtReferencia string) (*entity.CMEDConformidade, error) {
	return nil, nil
}

func (m *mockCMEDRepoForSupplier) GetByCNPJ(ctx context.Context, cnpj string, dtReferencia string) ([]entity.CMEDConformidade, error) {
	return m.cmedList, m.err
}

func (m *mockCMEDRepoForSupplier) List(ctx context.Context, filter repository.CMEDFilterParams) (*entity.CursorPage[entity.CMEDConformidade], error) {
	return nil, nil
}

func (m *mockCMEDRepoForSupplier) GetHistorico(ctx context.Context, nuSanReg int64) ([]entity.CMEDConformidade, error) {
	return nil, nil
}

func (m *mockCMEDRepoForSupplier) UpsertBatch(ctx context.Context, records []entity.CMEDConformidade) (int64, error) {
	return 0, nil
}

func TestSupplierCMEDUsecase_Success(t *testing.T) {
	cnpj := sql.NullString{String: "00000000000000", Valid: true}
	supplier := &entity.Supplier{COSeqID: 1, NUCnpj: cnpj}
	cmedList := []entity.CMEDConformidade{
		{COSeqID: 100, NUCnpj: strPtr("00000000000000")},
		{COSeqID: 101, NUCnpj: strPtr("00000000000000")},
	}

	supplierRepo := &mockSupplierRepoForCMED{supplier: supplier}
	cmedRepo := &mockCMEDRepoForSupplier{cmedList: cmedList}

	uc := NewSupplierCMEDUsecase(supplierRepo, cmedRepo, nil)

	result, err := uc.GetSupplierWithCMED(context.Background(), 1, "2024-01")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, supplier, result.Supplier)
	assert.Equal(t, cmedList, result.CMED)
}

func TestSupplierCMEDUsecase_CNPJNormalization(t *testing.T) {
	cnpj := sql.NullString{String: "00.000.000/0000-00", Valid: true}
	supplier := &entity.Supplier{COSeqID: 1, NUCnpj: cnpj}
	cmedList := []entity.CMEDConformidade{
		{COSeqID: 100, NUCnpj: strPtr("00000000000000")},
	}

	supplierRepo := &mockSupplierRepoForCMED{supplier: supplier}
	cmedRepo := &mockCMEDRepoForSupplier{cmedList: cmedList}

	uc := NewSupplierCMEDUsecase(supplierRepo, cmedRepo, nil)

	result, err := uc.GetSupplierWithCMED(context.Background(), 1, "2024-01")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, cmedList, result.CMED)
}


func TestSupplierCMEDUsecase_CNPJNotValid(t *testing.T) {
	supplier := &entity.Supplier{COSeqID: 1, NUCnpj: sql.NullString{Valid: false}}

	supplierRepo := &mockSupplierRepoForCMED{supplier: supplier}
	cmedRepo := &mockCMEDRepoForSupplier{}

	uc := NewSupplierCMEDUsecase(supplierRepo, cmedRepo, nil)

	result, err := uc.GetSupplierWithCMED(context.Background(), 1, "2024-01")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, supplier, result.Supplier)
	assert.Nil(t, result.CMED)
}

func TestSupplierCMEDUsecase_SupplierNotFound(t *testing.T) {
	supplierRepo := &mockSupplierRepoForCMED{supplier: nil}
	cmedRepo := &mockCMEDRepoForSupplier{}

	uc := NewSupplierCMEDUsecase(supplierRepo, cmedRepo, nil)

	result, err := uc.GetSupplierWithCMED(context.Background(), 999, "2024-01")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "Supplier not found")
}

func TestSupplierCMEDUsecase_SupplierGetError(t *testing.T) {
	supplierRepo := &mockSupplierRepoForCMED{err: errors.New("db error")}
	cmedRepo := &mockCMEDRepoForSupplier{}

	uc := NewSupplierCMEDUsecase(supplierRepo, cmedRepo, nil)

	_, err := uc.GetSupplierWithCMED(context.Background(), 1, "2024-01")
	assert.Error(t, err)
}

func TestSupplierCMEDUsecase_CMEDLookupFailure(t *testing.T) {
	cnpj := sql.NullString{String: "00000000000000", Valid: true}
	supplier := &entity.Supplier{COSeqID: 1, NUCnpj: cnpj}

	supplierRepo := &mockSupplierRepoForCMED{supplier: supplier}
	cmedRepo := &mockCMEDRepoForSupplier{err: errors.New("cmed db error")}

	uc := NewSupplierCMEDUsecase(supplierRepo, cmedRepo, nil)

	result, err := uc.GetSupplierWithCMED(context.Background(), 1, "2024-01")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, supplier, result.Supplier)
	assert.Empty(t, result.CMED)
}

func TestSupplierCMEDUsecase_CacheHit(t *testing.T) {
	cnpj := sql.NullString{String: "00000000000000", Valid: true}
	supplier := &entity.Supplier{COSeqID: 1, NUCnpj: cnpj}
	cmedList := []entity.CMEDConformidade{
		{COSeqID: 100, NUCnpj: strPtr("00000000000000")},
	}

	supplierRepo := &mockSupplierRepoForCMED{supplier: supplier}
	cmedRepo := &mockCMEDRepoForSupplier{cmedList: cmedList}

	uc := NewSupplierCMEDUsecase(supplierRepo, cmedRepo, newNilClientCacheRepo())

	result, err := uc.GetSupplierWithCMED(context.Background(), 1, "2024-01")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, supplier, result.Supplier)
	assert.Equal(t, cmedList, result.CMED)
}

func TestSupplierCMEDUsecase_CacheMiss(t *testing.T) {
	cnpj := sql.NullString{String: "00000000000000", Valid: true}
	supplier := &entity.Supplier{COSeqID: 1, NUCnpj: cnpj}
	cmedList := []entity.CMEDConformidade{
		{COSeqID: 100, NUCnpj: strPtr("00000000000000")},
	}

	supplierRepo := &mockSupplierRepoForCMED{supplier: supplier}
	cmedRepo := &mockCMEDRepoForSupplier{cmedList: cmedList}

	uc := NewSupplierCMEDUsecase(supplierRepo, cmedRepo, nil)

	result, err := uc.GetSupplierWithCMED(context.Background(), 1, "2024-01")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, cmedList, result.CMED)
}


