package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/silasrm/api-obm/internal/domain/entity"
)

type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool}
}

func (r *UserRepo) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	const query = `SELECT id, username, password_hash, active, created_at FROM users WHERE username = $1`
	var u entity.User
	err := r.pool.QueryRow(ctx, query, username).Scan(&u.ID, &u.Username, &u.PasswordHash, &u.Active, &u.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query user: %w", err)
	}
	return &u, nil
}

func (r *UserRepo) Create(ctx context.Context, user *entity.User) error {
	const query = `INSERT INTO users (username, password_hash, active, created_at) VALUES ($1, $2, $3, $4) RETURNING id`
	if user.CreatedAt.IsZero() {
		user.CreatedAt = time.Now()
	}
	return r.pool.QueryRow(ctx, query, user.Username, user.PasswordHash, user.Active, user.CreatedAt).Scan(&user.ID)
}
