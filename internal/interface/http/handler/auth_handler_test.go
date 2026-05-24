package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/silasrm/api-obm/internal/domain/entity"
	"github.com/silasrm/api-obm/internal/infrastructure/config"
	"github.com/silasrm/api-obm/internal/usecase"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

type mockUserRepoForHandler struct {
	user *entity.User
}

func (m *mockUserRepoForHandler) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	return m.user, nil
}

func (m *mockUserRepoForHandler) Create(ctx context.Context, user *entity.User) error {
	return nil
}

func handlerHashPassword(password string) string {
	h, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(h)
}

func TestAuthHandler_Login_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := handlerHashPassword("pass123")
	repo := &mockUserRepoForHandler{
		user: &entity.User{
			ID: 1, Username: "admin", PasswordHash: h, Active: true,
		},
	}
	uc := usecase.NewAuthUsecase(repo, config.JWTConfig{Secret: "secret", ExpirationHours: 24})
	handler := NewAuthHandler(uc)

	r := gin.New()
	r.POST("/auth/login", handler.Login)

	body, _ := json.Marshal(map[string]string{"username": "admin", "password": "pass123"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NotEmpty(t, resp["token"])
}

func TestAuthHandler_Login_InvalidBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &mockUserRepoForHandler{}
	uc := usecase.NewAuthUsecase(repo, config.JWTConfig{Secret: "secret", ExpirationHours: 24})
	handler := NewAuthHandler(uc)

	r := gin.New()
	r.POST("/auth/login", handler.Login)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewReader([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuthHandler_Login_InvalidCredentials(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &mockUserRepoForHandler{user: nil}
	uc := usecase.NewAuthUsecase(repo, config.JWTConfig{Secret: "secret", ExpirationHours: 24})
	handler := NewAuthHandler(uc)

	r := gin.New()
	r.POST("/auth/login", handler.Login)

	body, _ := json.Marshal(map[string]string{"username": "admin", "password": "wrong"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
