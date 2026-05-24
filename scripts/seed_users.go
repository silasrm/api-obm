package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	godotenv.Load()

	host := getEnv("PG_HOST", "localhost")
	port := getEnvInt("PG_PORT", 5433)
	user := getEnv("PG_USER", "obm")
	password := getEnv("PG_PASSWORD", "obm123")
	database := getEnv("PG_DATABASE", "dbportalobm")
	sslmode := getEnv("PG_SSLMODE", "disable")

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s database=%s sslmode=%s",
		host, port, user, password, database, sslmode)

	ctx := context.Background()
	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse config: %v\n", err)
		os.Exit(1)
	}
	poolCfg.MaxConns = 5

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "create pool: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "ping: %v\n", err)
		os.Exit(1)
	}

	users := []struct {
		username string
		password string
	}{
		{username: "admin", password: "admin123"},
		{username: "viewer", password: "viewer123"},
	}

	for _, u := range users {
		hash, err := bcrypt.GenerateFromPassword([]byte(u.password), bcrypt.DefaultCost)
		if err != nil {
			fmt.Fprintf(os.Stderr, "hash password for %s: %v\n", u.username, err)
			os.Exit(1)
		}

		var id int64
		err = pool.QueryRow(ctx,
			`INSERT INTO users (username, password_hash, active, created_at) VALUES ($1, $2, TRUE, NOW()) ON CONFLICT (username) DO UPDATE SET password_hash = EXCLUDED.password_hash RETURNING id`,
			u.username, string(hash),
		).Scan(&id)
		if err != nil {
			fmt.Fprintf(os.Stderr, "insert user %s: %v\n", u.username, err)
			os.Exit(1)
		}
		fmt.Printf("Seeded user: %s (id=%d)\n", u.username, id)
	}

	fmt.Println("Seed complete.")
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
