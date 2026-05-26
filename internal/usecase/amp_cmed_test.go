package usecase

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/silasrm/api-obm/internal/domain/entity"
	"github.com/silasrm/api-obm/internal/domain/repository"
	"github.com/silasrm/api-obm/internal/infrastructure/config"
	redisrepo "github.com/silasrm/api-obm/internal/infrastructure/persistence/redis"
	"github.com/stretchr/testify/assert"
)

type mockAMPRepoForAMPCMED struct {
	amp *entity.AMP
	err error
}

func (m *mockAMPRepoForAMPCMED) GetByID(ctx context.Context, id int64) (*entity.AMP, error) {
	return m.amp, m.err
}

func (m *mockAMPRepoForAMPCMED) GetDetailByID(ctx context.Context, id int64) (*entity.AMPDetail, error) {
	return nil, nil
}

func (m *mockAMPRepoForAMPCMED) List(ctx context.Context, filter repository.FilterParams) (*entity.CursorPage[entity.AMP], error) {
	return nil, nil
}

type mockVMPRepoForAMPCMED struct {
	vmp *entity.VMP
	err error
}

func (m *mockVMPRepoForAMPCMED) GetByID(ctx context.Context, id int64) (*entity.VMP, error) {
	return m.vmp, m.err
}

func (m *mockVMPRepoForAMPCMED) GetDetailByID(ctx context.Context, id int64) (*entity.VMPDetail, error) {
	return nil, nil
}

func (m *mockVMPRepoForAMPCMED) List(ctx context.Context, filter repository.FilterParams) (*entity.CursorPage[entity.VMP], error) {
	return nil, nil
}

type mockCMEDRepoForAMPCMED struct {
	cmed *entity.CMEDConformidade
	err  error
}

func (m *mockCMEDRepoForAMPCMED) GetByID(ctx context.Context, id int64) (*entity.CMEDConformidade, error) {
	return m.cmed, m.err
}

func (m *mockCMEDRepoForAMPCMED) GetByNuSanReg(ctx context.Context, nuSanReg int64, dtReferencia string) (*entity.CMEDConformidade, error) {
	return m.cmed, m.err
}

func (m *mockCMEDRepoForAMPCMED) GetByEAN(ctx context.Context, ean string, dtReferencia string) (*entity.CMEDConformidade, error) {
	return nil, nil
}

func (m *mockCMEDRepoForAMPCMED) GetByCNPJ(ctx context.Context, cnpj string, dtReferencia string) ([]entity.CMEDConformidade, error) {
	return nil, nil
}

func (m *mockCMEDRepoForAMPCMED) List(ctx context.Context, filter repository.CMEDFilterParams) (*entity.CursorPage[entity.CMEDConformidade], error) {
	return nil, nil
}

func (m *mockCMEDRepoForAMPCMED) GetHistorico(ctx context.Context, nuSanReg int64) ([]entity.CMEDConformidade, error) {
	return nil, nil
}

func (m *mockCMEDRepoForAMPCMED) UpsertBatch(ctx context.Context, records []entity.CMEDConformidade) (int64, error) {
	return 0, nil
}

func TestAMPCMEDUsecase_Success(t *testing.T) {
	nuNReg := sql.NullInt64{Int64: 12345, Valid: true}
	amp := &entity.AMP{COSeqID: 1, COVpID: 20, NUNReg: nuNReg}
	vmp := &entity.VMP{COSeqID: 20, NONm: "Test VMP"}
	cmed := &entity.CMEDConformidade{COSeqID: 100, NUSanReg: int64Ptr(12345)}

	ampRepo := &mockAMPRepoForAMPCMED{amp: amp}
	vmpRepo := &mockVMPRepoForAMPCMED{vmp: vmp}
	cmedRepo := &mockCMEDRepoForAMPCMED{cmed: cmed}

	uc := NewAMPCMEDUsecase(ampRepo, vmpRepo, cmedRepo, nil)

	result, err := uc.GetAMPWithCMED(context.Background(), 1, "2024-01-01")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, amp, result.AMP)
	assert.Equal(t, vmp, result.VMP)
	assert.Equal(t, cmed, result.CMED)
}

func TestAMPCMEDUsecase_NuNRegEmpty(t *testing.T) {
	amp := &entity.AMP{COSeqID: 1, COVpID: 20, NUNReg: sql.NullInt64{Valid: false}}
	vmp := &entity.VMP{COSeqID: 20, NONm: "Test VMP"}

	ampRepo := &mockAMPRepoForAMPCMED{amp: amp}
	vmpRepo := &mockVMPRepoForAMPCMED{vmp: vmp}
	cmedRepo := &mockCMEDRepoForAMPCMED{}

	uc := NewAMPCMEDUsecase(ampRepo, vmpRepo, cmedRepo, nil)

	result, err := uc.GetAMPWithCMED(context.Background(), 1, "2024-01-01")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, amp, result.AMP)
	assert.Nil(t, result.CMED)
}

func TestAMPCMEDUsecase_NuNRegZero(t *testing.T) {
	amp := &entity.AMP{COSeqID: 1, COVpID: 20, NUNReg: sql.NullInt64{Int64: 0, Valid: true}}
	vmp := &entity.VMP{COSeqID: 20, NONm: "Test VMP"}

	ampRepo := &mockAMPRepoForAMPCMED{amp: amp}
	vmpRepo := &mockVMPRepoForAMPCMED{vmp: vmp}
	cmedRepo := &mockCMEDRepoForAMPCMED{}

	uc := NewAMPCMEDUsecase(ampRepo, vmpRepo, cmedRepo, nil)

	result, err := uc.GetAMPWithCMED(context.Background(), 1, "2024-01-01")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Nil(t, result.CMED)
}

func TestAMPCMEDUsecase_AMPNotFound(t *testing.T) {
	ampRepo := &mockAMPRepoForAMPCMED{err: errors.New("not found")}
	vmpRepo := &mockVMPRepoForAMPCMED{}
	cmedRepo := &mockCMEDRepoForAMPCMED{}

	uc := NewAMPCMEDUsecase(ampRepo, vmpRepo, cmedRepo, nil)

	_, err := uc.GetAMPWithCMED(context.Background(), 999, "2024-01-01")
	assert.Error(t, err)
}

func TestAMPCMEDUsecase_AMPNil(t *testing.T) {
	ampRepo := &mockAMPRepoForAMPCMED{amp: nil}
	vmpRepo := &mockVMPRepoForAMPCMED{}
	cmedRepo := &mockCMEDRepoForAMPCMED{}

	uc := NewAMPCMEDUsecase(ampRepo, vmpRepo, cmedRepo, nil)

	_, err := uc.GetAMPWithCMED(context.Background(), 999, "2024-01-01")
	assert.Error(t, err)
}

func TestAMPCMEDUsecase_CMEDGracefulFailure(t *testing.T) {
	nuNReg := sql.NullInt64{Int64: 12345, Valid: true}
	amp := &entity.AMP{COSeqID: 1, COVpID: 20, NUNReg: nuNReg}
	vmp := &entity.VMP{COSeqID: 20, NONm: "Test VMP"}

	ampRepo := &mockAMPRepoForAMPCMED{amp: amp}
	vmpRepo := &mockVMPRepoForAMPCMED{vmp: vmp}
	cmedRepo := &mockCMEDRepoForAMPCMED{err: errors.New("cmed db error")}

	uc := NewAMPCMEDUsecase(ampRepo, vmpRepo, cmedRepo, nil)

	result, err := uc.GetAMPWithCMED(context.Background(), 1, "2024-01-01")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, amp, result.AMP)
	assert.Nil(t, result.CMED)
}

func TestAMPCMEDUsecase_CMEDNotFound(t *testing.T) {
	nuNReg := sql.NullInt64{Int64: 12345, Valid: true}
	amp := &entity.AMP{COSeqID: 1, COVpID: 20, NUNReg: nuNReg}
	vmp := &entity.VMP{COSeqID: 20, NONm: "Test VMP"}

	ampRepo := &mockAMPRepoForAMPCMED{amp: amp}
	vmpRepo := &mockVMPRepoForAMPCMED{vmp: vmp}
	cmedRepo := &mockCMEDRepoForAMPCMED{cmed: nil}

	uc := NewAMPCMEDUsecase(ampRepo, vmpRepo, cmedRepo, nil)

	result, err := uc.GetAMPWithCMED(context.Background(), 1, "2024-01-01")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Nil(t, result.CMED)
}

func TestAMPCMEDUsecase_CacheHit(t *testing.T) {
	nuNReg := sql.NullInt64{Int64: 12345, Valid: true}
	amp := &entity.AMP{COSeqID: 1, COVpID: 20, NUNReg: nuNReg}
	vmp := &entity.VMP{COSeqID: 20, NONm: "Test VMP"}
	cmed := &entity.CMEDConformidade{COSeqID: 100, NUSanReg: int64Ptr(12345)}

	ampRepo := &mockAMPRepoForAMPCMED{amp: amp}
	vmpRepo := &mockVMPRepoForAMPCMED{vmp: vmp}
	cmedRepo := &mockCMEDRepoForAMPCMED{cmed: cmed}

	uc := NewAMPCMEDUsecase(ampRepo, vmpRepo, cmedRepo, newNilClientCacheRepo())

	result, err := uc.GetAMPWithCMED(context.Background(), 1, "2024-01-01")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, cmed, result.CMED)
}

func TestAMPCMEDUsecase_NoCache(t *testing.T) {
	nuNReg := sql.NullInt64{Int64: 12345, Valid: true}
	amp := &entity.AMP{COSeqID: 1, COVpID: 20, NUNReg: nuNReg}
	vmp := &entity.VMP{COSeqID: 20, NONm: "Test VMP"}
	cmed := &entity.CMEDConformidade{COSeqID: 100, NUSanReg: int64Ptr(12345)}

	ampRepo := &mockAMPRepoForAMPCMED{amp: amp}
	vmpRepo := &mockVMPRepoForAMPCMED{vmp: vmp}
	cmedRepo := &mockCMEDRepoForAMPCMED{cmed: cmed}

	uc := NewAMPCMEDUsecase(ampRepo, vmpRepo, cmedRepo, nil)

	result, err := uc.GetAMPWithCMED(context.Background(), 1, "2024-01-01")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, cmed, result.CMED)
}

func newNilClientCacheRepoForAMPCMED() *redisrepo.CacheRepo {
	return redisrepo.NewCacheRepo(config.RedisConfig{})
}
