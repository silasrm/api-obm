package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/silasrm/api-obm/internal/interface/http/dto"
	"github.com/silasrm/api-obm/internal/usecase"
)

type AuthHandler struct {
	authUseCase *usecase.AuthUsecase
}

func NewAuthHandler(uc *usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{authUseCase: uc}
}

// Login godoc
// @Summary Autenticar usuário
// @Description Realiza login e retorna token JWT
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body dto.LoginRequest true "Credenciais"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request body", Code: http.StatusBadRequest})
		return
	}

	token, expiresIn, err := h.authUseCase.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "invalid credentials", Code: http.StatusUnauthorized})
		return
	}

	c.JSON(http.StatusOK, dto.LoginResponse{Token: token, ExpiresIn: expiresIn})
}
