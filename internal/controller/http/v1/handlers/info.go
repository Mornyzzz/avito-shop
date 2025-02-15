package handlers

import (
	mw "avito-shop/internal/controller/middleware"
	"avito-shop/internal/entity"
	"avito-shop/internal/usecase/info"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
	"net/http"
)

type InfoRoute struct {
	infoUC info.Info
	log    *slog.Logger
}

func NewInfoRoute(handler *gin.RouterGroup, infoUC info.Info, log *slog.Logger) {
	r := &InfoRoute{infoUC, log}
	handler.GET("/info", mw.Auth(), r.Info)
}

type InfoRequest struct {
	Item string `uri:"item" binding:"required"`
}

func (r *InfoRoute) Info(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		r.log.Error("Username not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	information, err := r.infoUC.GetInfo(c.Request.Context(), username.(string))
	if err != nil {
		r.log.Error("Failed to get info", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get info"})
		return
	}

	c.JSON(http.StatusOK, information)
}

type InfoResponse struct {
	Coins       int                    `json:"coins"`
	Inventory   []entity.InventoryItem `json:"inventory"`
	CoinHistory entity.CoinHistory     `json:"coinHistory"`
}
