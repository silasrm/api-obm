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

type mockAMPPRepoForCMED struct {
	ampp interface{}
	err  error
}

func (m *mockAMPPRepoForCMED) GetByID(ctx context.Context, id int64) (interface{}, error) {
	return m.ampp, m.err
}

func (m *mockAMPPRepoForCMED) List(ctx context.Context, filter repository.FilterParams) (interface{}, error) {
	return nil, nil
}

type mockVMPRepoForAMPPCMED struct {
	vmp *entity.VMP
	err error
}

func (m *mockVMPRepoForAMPPCMED) GetByID(ctx context.Context, id int64) (*entity.VMP, error) {
	return m.vmp, m.err
}

func (m *mockVMPRepoForAMPPCMED) GetDetailByID(ctx context.Context, id int64) (*entity.VMPDetail, error) {
	return nil, nil
}

func (m *mockVMPRepoForAMPPCMED) List(ctx context.Context, filter repository.FilterParams) (*entity.CursorPage[entity.VMP], error) {
	return nil, nil
}

type mockAMPRepoForAMPPCMED struct {
	amp *entity.AMP
	err error
}

func (m *mockAMPRepoForAMPPCMED) GetByID(ctx context.Context, id int64) (*entity.AMP, error) {
	return m.amp, m.err
}

func (m *mockAMPRepoForAMPPCMED) GetDetailByID(ctx context.Context, id int64) (*entity.AMPDetail, error) {
	return nil, nil
}

func (m *mockAMPRepoForAMPPCMED) List(ctx context.Context, filter repository.FilterParams) (*entity.CursorPage[entity.AMP], error) {
	return nil, nil
}

type mockCMEDRepoForAMPP struct {
	cmed *entity.CMEDConformidade
	err  error
}

func (m *mockCMEDRepoForAMPP) GetByID(ctx context.Context, id int64) (*entity.CMEDConformidade, error) {
	return m.cmed, m.err
}

func (m *mockCMEDRepoForAMPP) GetByNuSanReg(ctx context.Context, nuSanReg int64, dtReferencia string) (*entity.CMEDConformidade, error) {
	return m.cmed, m.err
}

func (m *mockCMEDRepoForAMPP) GetByEAN(ctx context.Context, ean string, dtReferencia string) (*entity.CMEDConformidade, error) {
	return nil, nil
}

func (m *mockCMEDRepoForAMPP) List(ctx context.Context, filter repository.CMEDFilterParams) (*entity.CursorPage[entity.CMEDConformidade], error) {
	return nil, nil
}

func (m *mockCMEDRepoForAMPP) GetHistorico(ctx context.Context, nuSanReg int64) ([]entity.CMEDConformidade, error) {
	return nil, nil
}

func (m *mockCMEDRepoForAMPP) UpsertBatch(ctx context.Context, records []entity.CMEDConformidade) (int64, error) {
	return 0, nil
}

func newNilClientCacheRepo() *redisrepo.CacheRepo {
	return redisrepo.NewCacheRepo(config.RedisConfig{})
}

func TestAMPPCMEDUsecase_Success(t *testing.T) {
	sanReg := sql.NullInt64{Int64: 12345, Valid: true}
	ampp := &entity.AMPP{COSeqID: 1, COApID: 10, NUSanReg: sanReg}
	amp := &entity.AMP{COSeqID: 10, COVpID: 20}
	vmp := &entity.VMP{COSeqID: 20, NONm: "Test VMP"}
	cmed := &entity.CMEDConformidade{COSeqID: 100, NUSanReg: int64Ptr(12345)}

	amppRepo := &mockAMPPRepoForCMED{ampp: ampp}
	ampRepo := &mockAMPRepoForAMPPCMED{amp: amp}
	vmpRepo := &mockVMPRepoForAMPPCMED{vmp: vmp}
	cmedRepo := &mockCMEDRepoForAMPP{cmed: cmed}

	uc := NewAMPPCMEDUsecase(amppRepo, vmpRepo, ampRepo, cmedRepo, nil)

	result, err := uc.GetAMPPWithCMED(context.Background(), 1, "")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, ampp, result.AMPP)
	assert.Equal(t, amp, result.AMP)
	assert.Equal(t, vmp, result.VMP)
	assert.Equal(t, cmed, result.CMED)
}

func TestAMPPCMEDUsecase_AMPPNoSanReg(t *testing.T) {
	ampp := &entity.AMPP{COSeqID: 1, COApID: 10, NUSanReg: sql.NullInt64{Valid: false}}
	amp := &entity.AMP{COSeqID: 10, COVpID: 20}
	vmp := &entity.VMP{COSeqID: 20, NONm: "Test VMP"}

	amppRepo := &mockAMPPRepoForCMED{ampp: ampp}
	ampRepo := &mockAMPRepoForAMPPCMED{amp: amp}
	vmpRepo := &mockVMPRepoForAMPPCMED{vmp: vmp}
	cmedRepo := &mockCMEDRepoForAMPP{}

	uc := NewAMPPCMEDUsecase(amppRepo, vmpRepo, ampRepo, cmedRepo, nil)

	result, err := uc.GetAMPPWithCMED(context.Background(), 1, "")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, ampp, result.AMPP)
	assert.Nil(t, result.CMED)
}

func TestAMPPCMEDUsecase_AMPPNoSanRegZero(t *testing.T) {
	ampp := &entity.AMPP{COSeqID: 1, COApID: 10, NUSanReg: sql.NullInt64{Int64: 0, Valid: true}}
	amp := &entity.AMP{COSeqID: 10, COVpID: 20}
	vmp := &entity.VMP{COSeqID: 20, NONm: "Test VMP"}

	amppRepo := &mockAMPPRepoForCMED{ampp: ampp}
	ampRepo := &mockAMPRepoForAMPPCMED{amp: amp}
	vmpRepo := &mockVMPRepoForAMPPCMED{vmp: vmp}
	cmedRepo := &mockCMEDRepoForAMPP{}

	uc := NewAMPPCMEDUsecase(amppRepo, vmpRepo, ampRepo, cmedRepo, nil)

	result, err := uc.GetAMPPWithCMED(context.Background(), 1, "")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Nil(t, result.CMED)
}

func TestAMPPCMEDUsecase_CMEDNotFound(t *testing.T) {
	sanReg := sql.NullInt64{Int64: 12345, Valid: true}
	ampp := &entity.AMPP{COSeqID: 1, COApID: 10, NUSanReg: sanReg}
	amp := &entity.AMP{COSeqID: 10, COVpID: 20}
	vmp := &entity.VMP{COSeqID: 20, NONm: "Test VMP"}

	amppRepo := &mockAMPPRepoForCMED{ampp: ampp}
	ampRepo := &mockAMPRepoForAMPPCMED{amp: amp}
	vmpRepo := &mockVMPRepoForAMPPCMED{vmp: vmp}
	cmedRepo := &mockCMEDRepoForAMPP{cmed: nil}

	uc := NewAMPPCMEDUsecase(amppRepo, vmpRepo, ampRepo, cmedRepo, nil)

	result, err := uc.GetAMPPWithCMED(context.Background(), 1, "")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Nil(t, result.CMED)
}

func TestAMPPCMEDUsecase_CacheHit(t *testing.T) {
	sanReg := sql.NullInt64{Int64: 12345, Valid: true}
	ampp := &entity.AMPP{COSeqID: 1, COApID: 10, NUSanReg: sanReg}
	amp := &entity.AMP{COSeqID: 10, COVpID: 20}
	vmp := &entity.VMP{COSeqID: 20, NONm: "Test VMP"}
	cmed := &entity.CMEDConformidade{COSeqID: 100, NUSanReg: int64Ptr(12345)}

	amppRepo := &mockAMPPRepoForCMED{ampp: ampp}
	ampRepo := &mockAMPRepoForAMPPCMED{amp: amp}
	vmpRepo := &mockVMPRepoForAMPPCMED{vmp: vmp}
	cmedRepo := &mockCMEDRepoForAMPP{cmed: cmed}

	uc := NewAMPPCMEDUsecase(amppRepo, vmpRepo, ampRepo, cmedRepo, newNilClientCacheRepo())

	result, err := uc.GetAMPPWithCMED(context.Background(), 1, "")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, cmed, result.CMED)
}

func TestAMPPCMEDUsecase_NoCache(t *testing.T) {
	sanReg := sql.NullInt64{Int64: 12345, Valid: true}
	ampp := &entity.AMPP{COSeqID: 1, COApID: 10, NUSanReg: sanReg}
	amp := &entity.AMP{COSeqID: 10, COVpID: 20}
	vmp := &entity.VMP{COSeqID: 20, NONm: "Test VMP"}
	cmed := &entity.CMEDConformidade{COSeqID: 100, NUSanReg: int64Ptr(12345)}

	amppRepo := &mockAMPPRepoForCMED{ampp: ampp}
	ampRepo := &mockAMPRepoForAMPPCMED{amp: amp}
	vmpRepo := &mockVMPRepoForAMPPCMED{vmp: vmp}
	cmedRepo := &mockCMEDRepoForAMPP{cmed: cmed}

	uc := NewAMPPCMEDUsecase(amppRepo, vmpRepo, ampRepo, cmedRepo, nil)

	result, err := uc.GetAMPPWithCMED(context.Background(), 1, "")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, cmed, result.CMED)
}

func TestAMPPCMEDUsecase_AMPPNil(t *testing.T) {
	amppRepo := &mockAMPPRepoForCMED{ampp: nil}
	ampRepo := &mockAMPRepoForAMPPCMED{}
	vmpRepo := &mockVMPRepoForAMPPCMED{}
	cmedRepo := &mockCMEDRepoForAMPP{}

	uc := NewAMPPCMEDUsecase(amppRepo, vmpRepo, ampRepo, cmedRepo, nil)

	result, err := uc.GetAMPPWithCMED(context.Background(), 999, "")
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestAMPPCMEDUsecase_AMPPGetError(t *testing.T) {
	amppRepo := &mockAMPPRepoForCMED{err: errors.New("db error")}
	ampRepo := &mockAMPRepoForAMPPCMED{}
	vmpRepo := &mockVMPRepoForAMPPCMED{}
	cmedRepo := &mockCMEDRepoForAMPP{}

	uc := NewAMPPCMEDUsecase(amppRepo, vmpRepo, ampRepo, cmedRepo, nil)

	_, err := uc.GetAMPPWithCMED(context.Background(), 1, "")
	assert.Error(t, err)
}
