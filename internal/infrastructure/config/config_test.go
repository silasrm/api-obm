package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnv(t *testing.T) {
	os.Clearenv()
	assert.Equal(t, "fallback", getEnv("TEST_KEY", "fallback"))

	t.Setenv("TEST_KEY", "value")
	assert.Equal(t, "value", getEnv("TEST_KEY", "fallback"))

	t.Setenv("TEST_EMPTY", "")
	assert.Equal(t, "fallback", getEnv("TEST_EMPTY", "fallback"))
}

func TestGetEnvInt(t *testing.T) {
	os.Clearenv()
	assert.Equal(t, 42, getEnvInt("TEST_INT", 42))

	t.Setenv("TEST_INT", "8094")
	assert.Equal(t, 8094, getEnvInt("TEST_INT", 42))

	t.Setenv("TEST_BAD_INT", "abc")
	assert.Equal(t, 0, getEnvInt("TEST_BAD_INT", 0))
}

func TestGetEnvBool(t *testing.T) {
	os.Clearenv()
	assert.True(t, getEnvBool("TEST_BOOL", true))
	assert.False(t, getEnvBool("TEST_BOOL_F", false))

	t.Setenv("TEST_BOOL", "true")
	assert.True(t, getEnvBool("TEST_BOOL", false))

	t.Setenv("TEST_BOOL_0", "0")
	assert.False(t, getEnvBool("TEST_BOOL_0", true))

	t.Setenv("TEST_BOOL_BAD", "maybe")
	assert.True(t, getEnvBool("TEST_BOOL_BAD", true))
}

func TestLoad_Defaults(t *testing.T) {
	os.Clearenv()
	cfg := Load()
	assert.Equal(t, "localhost", cfg.PostgreSQL.Host)
	assert.Equal(t, 5432, cfg.PostgreSQL.Port)
	assert.Equal(t, "obm", cfg.PostgreSQL.User)
	assert.Equal(t, "disable", cfg.PostgreSQL.SSLMode)
	assert.Equal(t, "http://localhost:7700", cfg.Meilisearch.URL)
	assert.Equal(t, "obm_", cfg.Meilisearch.IndexPrefix)
	assert.Equal(t, "obm-secret-key-change-in-prod", cfg.JWT.Secret)
	assert.Equal(t, 24, cfg.JWT.ExpirationHours)
	assert.Equal(t, 8080, cfg.Server.Port)
	assert.Equal(t, "release", cfg.Server.GinMode)
	assert.True(t, cfg.Sync.OnStartup)
}

func TestLoad_Overrides(t *testing.T) {
	os.Clearenv()
	t.Setenv("PG_HOST", "db.example.com")
	t.Setenv("PG_PORT", "5433")
	t.Setenv("SERVER_PORT", "8094")
	t.Setenv("JWT_SECRET", "my-secret")
	t.Setenv("SYNC_ON_STARTUP", "false")

	cfg := Load()
	assert.Equal(t, "db.example.com", cfg.PostgreSQL.Host)
	assert.Equal(t, 5433, cfg.PostgreSQL.Port)
	assert.Equal(t, 8094, cfg.Server.Port)
	assert.Equal(t, "my-secret", cfg.JWT.Secret)
	assert.False(t, cfg.Sync.OnStartup)
}
