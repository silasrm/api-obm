package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/silasrm/api-obm/internal/domain/repository"
	"github.com/silasrm/api-obm/internal/infrastructure/config"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
	userRepo repository.UserRepository
	jwtCfg   config.JWTConfig
}

func NewAuthUsecase(userRepo repository.UserRepository, jwtCfg config.JWTConfig) *AuthUsecase {
	return &AuthUsecase{userRepo: userRepo, jwtCfg: jwtCfg}
}

func (u *AuthUsecase) Login(ctx context.Context, username, password string) (string, int, error) {
	user, err := u.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return "", 0, fmt.Errorf("get user: %w", err)
	}
	if user == nil {
		return "", 0, errors.New("invalid credentials")
	}
	if !user.Active {
		return "", 0, errors.New("user inactive")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", 0, errors.New("invalid credentials")
	}

	expiresIn := u.jwtCfg.ExpirationHours * 3600
	exp := time.Now().Add(time.Duration(u.jwtCfg.ExpirationHours) * time.Hour)

	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      exp.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(u.jwtCfg.Secret))
	if err != nil {
		return "", 0, fmt.Errorf("sign token: %w", err)
	}

	return tokenStr, expiresIn, nil
}
