package handlers

import (
	"avito-shop/internal/controller/worker"
	"avito-shop/internal/usecase/buy"
	e "avito-shop/pkg/errors"
	mw "avito-shop/pkg/jwt"
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
	"net/http"
)

type BuyRoute struct {
	buyUC buy.Buy
	log   *slog.Logger
	wp    worker.PoolI
}

func NewBuyRoute(handler *gin.RouterGroup, buyUC buy.Buy, wp worker.PoolI, log *slog.Logger) {
	r := &BuyRoute{buyUC, log, wp}
	handler.GET("/buy/:item", mw.AuthMW(), r.Buy)
}

type BuyRequest struct {
	Item string `uri:"item" binding:"required"`
}

func (r *BuyRoute) Buy(c *gin.Context) {
	resultChan := make(chan string, 1)
	errorChan := make(chan error, 1)

	username, exists := c.Get("username")
	if !exists {
		r.log.Error("Username not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req BuyRequest
	if err := c.ShouldBindUri(&req); err != nil {
		r.log.Error("Failed to parse request", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	r.wp.Submit(func() {
		err := r.buyUC.BuyItem(c.Request.Context(), username.(string), req.Item)
		if err != nil {
			errorChan <- err
			return
		}

		resultChan <- "Item purchased successfully"
	})

	select {

	case result := <-resultChan:
		c.JSON(http.StatusOK, result) // Успешный ответ
	case err := <-errorChan:
		r.log.Error("Failed to buy item", slog.String("error", err.Error()))
		switch {
		case errors.Is(err, e.ErrInvalidCredentials), errors.Is(err, e.ErrNotFound):
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid credentials"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
	}

}
