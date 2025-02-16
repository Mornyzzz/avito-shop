package handlers

import (
	"avito-shop/internal/usecase/buy"
	mw "avito-shop/pkg/jwt"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
	"net/http"
)

type BuyRoute struct {
	buyUC buy.Buy
	log   *slog.Logger
}

func NewBuyRoute(handler *gin.RouterGroup, buyUC buy.Buy, log *slog.Logger) {
	r := &BuyRoute{buyUC, log}
	handler.GET("/buy/:item", mw.Auth(), r.Buy)
}

// BuyRequest представляет запрос на покупку товара
type BuyRequest struct {
	Item string `uri:"item" binding:"required"` // Название товара
}

func (r *BuyRoute) Buy(c *gin.Context) {
	var req BuyRequest

	if err := c.ShouldBindUri(&req); err != nil {
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

	err := r.buyUC.BuyItem(c.Request.Context(), username.(string), req.Item)
	if err != nil {
		r.log.Error("Failed to buy item", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to buy item"})
		return
	}

	c.JSON(http.StatusOK, "Item purchased successfully")
}
