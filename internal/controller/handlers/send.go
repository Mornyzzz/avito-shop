package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"

	"avito-shop/internal/controller/worker"
	"avito-shop/internal/usecase/send"
	e "avito-shop/pkg/errors"
	mw "avito-shop/pkg/jwt"
)

type SendRoute struct {
	sendUC send.Send
	log    *slog.Logger
	wp     worker.PoolI
}

func NewSendRoute(handler *gin.RouterGroup, sendUC send.Send, wp worker.PoolI, log *slog.Logger) {
	r := &SendRoute{sendUC, log, wp}
	handler.POST("/sendCoin", mw.AuthMW(), r.Send)
}

type SendRequest struct {
	ToUser string `json:"toUser" binding:"required"`
	Amount int    `json:"amount" binding:"required,gt=0,lte=1000000"`
}

func (r *SendRoute) Send(c *gin.Context) {
	resultChan := make(chan string, 1)
	errorChan := make(chan error, 1)

	username, exists := c.Get("username")
	if !exists {
		r.log.Error("Username not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})

		return
	}

	var req SendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		r.log.Error("Failed to parse request", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})

		return
	}

	r.wp.Submit(func() {
		err := r.sendUC.SendCoin(c.Request.Context(), username.(string), req.ToUser, req.Amount)
		if err != nil {
			errorChan <- err

			return
		}

		resultChan <- "Coins sent successfully"
	})

	select {
	case result := <-resultChan:
		c.JSON(http.StatusOK, result)
	case err := <-errorChan:
		r.log.Error("Failed to send coins", slog.String("error", err.Error()))

		switch {
		case errors.Is(err, e.ErrInvalidCredentials), errors.Is(err, e.ErrNotFound):
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid credentials"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
	}
}
