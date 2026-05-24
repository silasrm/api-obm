//go:build integration

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/silasrm/api-obm/internal/domain/entity"
	"github.com/silasrm/api-obm/internal/infrastructure/config"
	postgresrepo "github.com/silasrm/api-obm/internal/infrastructure/persistence/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func startPostgres(ctx context.Context, t *testing.T) (*pgxpool.Pool, func()) {
	t.Helper()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:16-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "obm",
			"POSTGRES_PASSWORD": "obm123",
			"POSTGRES_DB":       "dbportalobm",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").
			WithOccurrence(2).
			WithStartupTimeout(30 * time.Second),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	host, err := container.Host(ctx)
	require.NoError(t, err)
	natPort, err := container.MappedPort(ctx, "5432")
	require.NoError(t, err)

	cfg := config.PostgresConfig{
		Host:     host,
		Port:     int(natPort.Num()),
		User:     "obm",
		Password: "obm123",
		Database: "dbportalobm",
		SSLMode:  "disable",
	}

	pool, err := postgresrepo.NewPool(ctx, cfg)
	require.NoError(t, err)

	cleanup := func() {
		pool.Close()
		container.Terminate(ctx)
	}

	return pool, cleanup
}

func setupSchema(ctx context.Context, t *testing.T, pool *pgxpool.Pool) {
	t.Helper()
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username TEXT NOT NULL UNIQUE,
			password_hash TEXT NOT NULL,
			active BOOLEAN NOT NULL DEFAULT true,
			created_at TIMESTAMP NOT NULL DEFAULT NOW()
		)`,
	}
	for _, q := range queries {
		_, err := pool.Exec(ctx, q)
		require.NoError(t, err)
	}
}

func TestIntegration_Postgres_UserRepo(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	ctx := context.Background()
	pool, cleanup := startPostgres(ctx, t)
	defer cleanup()

	setupSchema(ctx, t, pool)

	repo := postgresrepo.NewUserRepo(pool)

	t.Run("GetByUsername not found", func(t *testing.T) {
		user, err := repo.GetByUsername(ctx, "nonexistent")
		assert.NoError(t, err)
		assert.Nil(t, user)
	})

	t.Run("Create and GetByUsername", func(t *testing.T) {
		user := &entity.User{
			Username:     "testuser",
			PasswordHash: "$2a$10$hashed",
			Active:       true,
		}
		err := repo.Create(ctx, user)
		require.NoError(t, err)
		assert.Greater(t, user.ID, int64(0))

		found, err := repo.GetByUsername(ctx, "testuser")
		require.NoError(t, err)
		require.NotNil(t, found)
		assert.Equal(t, "testuser", found.Username)
		assert.Equal(t, "$2a$10$hashed", found.PasswordHash)
		assert.True(t, found.Active)
	})
}

func TestIntegration_Postgres_NewPool(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	ctx := context.Background()
	pool, cleanup := startPostgres(ctx, t)
	defer cleanup()

	err := pool.Ping(ctx)
	assert.NoError(t, err)

	stats := pool.Stat()
	assert.Equal(t, int32(25), stats.MaxConns())
}
