package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	PostgreSQL PostgresConfig
	Meilisearch MeilisearchConfig
	JWT JWTConfig
	Server ServerConfig
	Sync SyncConfig
	Redis RedisConfig
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
	CacheTTL int
}

type PostgresConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	SSLMode  string
}

type MeilisearchConfig struct {
	URL         string
	APIKey      string
	IndexPrefix string
}

type JWTConfig struct {
	Secret           string
	ExpirationHours int
}

type ServerConfig struct {
	Port    int
	GinMode string
}

type SyncConfig struct {
	OnStartup bool
}

func Load() *Config {
	godotenv.Load()

	return &Config{
		PostgreSQL: PostgresConfig{
			Host:     getEnv("PG_HOST", "localhost"),
			Port:     getEnvInt("PG_PORT", 5432),
			User:     getEnv("PG_USER", "obm"),
			Password: getEnv("PG_PASSWORD", "obm123"),
			Database: getEnv("PG_DATABASE", "dbportalobm"),
			SSLMode:  getEnv("PG_SSLMODE", "disable"),
		},
		Meilisearch: MeilisearchConfig{
			URL:         getEnv("MEILI_URL", "http://localhost:7700"),
			APIKey:      getEnv("MEILI_API_KEY", ""),
			IndexPrefix: getEnv("MEILI_INDEX_PREFIX", "obm_"),
		},
		JWT: JWTConfig{
			Secret:           getEnv("JWT_SECRET", "obm-secret-key-change-in-prod"),
			ExpirationHours: getEnvInt("JWT_EXPIRATION_HOURS", 24),
		},
		Server: ServerConfig{
			Port:    getEnvInt("SERVER_PORT", 8080),
			GinMode: getEnv("GIN_MODE", "release"),
		},
		Sync: SyncConfig{
			OnStartup: getEnvBool("SYNC_ON_STARTUP", true),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnvInt("REDIS_PORT", 6380),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
			CacheTTL: getEnvInt("REDIS_CACHE_TTL", 24),
		},
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	if v := os.Getenv(key); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}
	return fallback
}
