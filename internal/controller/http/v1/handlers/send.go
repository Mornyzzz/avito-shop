package handlers

import (
	mw "avito-shop/internal/controller/middleware"
	"avito-shop/internal/usecase/send"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
	"net/http"
)

type SendRoute struct {
	sendUC send.Send
	log    *slog.Logger
}

func NewSendRoute(handler *gin.RouterGroup, sendUC send.Send, log *slog.Logger) {
	r := &SendRoute{sendUC, log}
	handler.POST("/sendCoin", mw.Auth(), r.Send)
}

type SendRequest struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}

func (r *SendRoute) Send(c *gin.Context) {
	var req SendRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		r.log.Error("Failed to parse request", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	username, exists := c.Get("username")
	if !exists {
		r.log.Error("Username not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	err := r.sendUC.SendCoin(c.Request.Context(), username.(string), req.ToUser, req.Amount)
	if err != nil {
		r.log.Error("Failed to send coins", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send coins"})
		return
	}

	c.JSON(http.StatusOK, "Coins send successfully")
}
