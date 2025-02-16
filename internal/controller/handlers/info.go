package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"

	"avito-shop/internal/controller/worker"
	"avito-shop/internal/entity"
	"avito-shop/internal/usecase/info"
	e "avito-shop/pkg/errors"
	mw "avito-shop/pkg/jwt"
)

type InfoRoute struct {
	infoUC info.Info
	log    *slog.Logger
	wp     worker.PoolI
}

func NewInfoRoute(handler *gin.RouterGroup, infoUC info.Info, wp worker.PoolI, log *slog.Logger) {
	r := &InfoRoute{infoUC, log, wp}
	handler.GET("/info", mw.AuthMW(), r.Info)
}

func (r *InfoRoute) Info(c *gin.Context) {
	resultChan := make(chan interface{}, 1)
	errorChan := make(chan error, 1)

	username, exists := c.Get("username")
	if !exists {
		r.log.Error("Username not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})

		return
	}

	r.wp.Submit(func() {
		information, err := r.infoUC.GetInfo(c.Request.Context(), username.(string))
		if err != nil {
			errorChan <- err
			return
		}

		resultChan <- information
	})

	select {
	case result := <-resultChan:
		c.JSON(http.StatusOK, result)
	case err := <-errorChan:
		r.log.Error("Failed to get info", slog.String("error", err.Error()))

		switch {
		case errors.Is(err, e.ErrInvalidCredentials), errors.Is(err, e.ErrNotFound):
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid credentials"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
	}
}

type InfoResponse struct {
	Coins       int                    `json:"coins"`
	Inventory   []entity.InventoryItem `json:"inventory"`
	CoinHistory entity.CoinHistory     `json:"coinHistory"`
}
