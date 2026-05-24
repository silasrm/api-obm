package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/silasrm/api-obm/internal/domain/entity"
	"github.com/silasrm/api-obm/internal/infrastructure/config"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

type mockUserRepo struct {
	user *entity.User
	err  error
}

func (m *mockUserRepo) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.user, nil
}

func (m *mockUserRepo) Create(ctx context.Context, user *entity.User) error {
	return nil
}

func hashPassword(password string) string {
	h, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(h)
}

func TestAuthUsecase_Login_Success(t *testing.T) {
	repo := &mockUserRepo{
		user: &entity.User{
			ID:           1,
			Username:     "admin",
			PasswordHash: hashPassword("password123"),
			Active:       true,
		},
	}
	uc := NewAuthUsecase(repo, config.JWTConfig{Secret: "secret", ExpirationHours: 24})

	token, expiresIn, err := uc.Login(context.Background(), "admin", "password123")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.Equal(t, 24*3600, expiresIn)
}

func TestAuthUsecase_Login_UserNotFound(t *testing.T) {
	repo := &mockUserRepo{user: nil}
	uc := NewAuthUsecase(repo, config.JWTConfig{Secret: "secret", ExpirationHours: 24})

	_, _, err := uc.Login(context.Background(), "admin", "password123")
	assert.Error(t, err)
	assert.Equal(t, "invalid credentials", err.Error())
}

func TestAuthUsecase_Login_WrongPassword(t *testing.T) {
	repo := &mockUserRepo{
		user: &entity.User{
			ID:           1,
			Username:     "admin",
			PasswordHash: hashPassword("password123"),
			Active:       true,
		},
	}
	uc := NewAuthUsecase(repo, config.JWTConfig{Secret: "secret", ExpirationHours: 24})

	_, _, err := uc.Login(context.Background(), "admin", "wrong")
	assert.Error(t, err)
	assert.Equal(t, "invalid credentials", err.Error())
}

func TestAuthUsecase_Login_InactiveUser(t *testing.T) {
	repo := &mockUserRepo{
		user: &entity.User{
			ID:           1,
			Username:     "admin",
			PasswordHash: hashPassword("password123"),
			Active:       false,
		},
	}
	uc := NewAuthUsecase(repo, config.JWTConfig{Secret: "secret", ExpirationHours: 24})

	_, _, err := uc.Login(context.Background(), "admin", "password123")
	assert.Error(t, err)
	assert.Equal(t, "user inactive", err.Error())
}

func TestAuthUsecase_Login_RepoError(t *testing.T) {
	repo := &mockUserRepo{err: errors.New("db error")}
	uc := NewAuthUsecase(repo, config.JWTConfig{Secret: "secret", ExpirationHours: 24})

	_, _, err := uc.Login(context.Background(), "admin", "password123")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "get user")
}
