package handlers

import (
	"avito-shop/internal/entity"
	"avito-shop/internal/usecase/auth"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
	"net/http"
)

type AuthRoute struct {
	authUC auth.Auth
	log    *slog.Logger
}

func NewAuthRoute(handler *gin.RouterGroup, authUC auth.Auth, log *slog.Logger) {
	r := &AuthRoute{authUC, log}
	handler.POST("/auth", r.Auth)
}

// AuthRequest представляет запрос на аутентификацию
type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r *AuthRoute) Auth(c *gin.Context) {
	var req AuthRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		r.log.Error("Failed to parse request", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	token, err := r.authUC.Login(c.Request.Context(), entity.User{Username: req.Username, Password: req.Password})
	if err != nil {
		r.log.Error("Authentication failed", slog.String("error", err.Error()))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{Token: token})
}

// AuthResponse представляет ответ на успешную аутентификацию
type AuthResponse struct {
	Token string `json:"token"`
}
