package handlers

import (
	"avito-shop/internal/controller/worker"
	"avito-shop/internal/entity"
	"avito-shop/internal/usecase/auth"
	e "avito-shop/pkg/errors"
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
	"net/http"
)

type AuthRoute struct {
	authUC auth.Auth
	log    *slog.Logger
	wp     worker.PoolI
}

func NewAuthRoute(handler *gin.RouterGroup, authUC auth.Auth, wp worker.PoolI, log *slog.Logger) {
	r := &AuthRoute{authUC, log, wp}
	handler.POST("/auth", r.Auth)
}

type AuthRequest struct {
	Username string `json:"username"  binding:"required"`
	Password string `json:"password"  binding:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

func (r *AuthRoute) Auth(c *gin.Context) {
	resultChan := make(chan AuthResponse, 1)
	errorChan := make(chan error, 1)

	r.wp.Submit(func() {
		var req AuthRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			errorChan <- err
			return
		}

		token, err := r.authUC.Login(c.Request.Context(), entity.User{
			Username: req.Username,
			Password: req.Password,
		})

		if err != nil {
			errorChan <- err
			return
		}

		resultChan <- AuthResponse{Token: token}
	})

	select {
	case authResponse := <-resultChan:
		c.JSON(http.StatusOK, authResponse)
	case err := <-errorChan:
		r.log.Error("Authentication failed", slog.String("error", err.Error()))
		switch {
		case errors.Is(err, e.ErrInvalidCredentials), errors.Is(err, e.ErrNotFound):
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid credentials"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
	}
}
