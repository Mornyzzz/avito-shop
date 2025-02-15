package v1

import (
	"golang.org/x/exp/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"avito-shop/internal/entity"
	"avito-shop/internal/usecase/auth"
	"avito-shop/internal/usecase/buy"
	"avito-shop/internal/usecase/info"
	"avito-shop/internal/usecase/send"
	"github.com/evrone/go-clean-template/pkg/logger"
)

type authRoutes struct {
	authUC auth.Auth
	log    *slog.Logger
}

func newTranslationRoutes(handler *gin.RouterGroup, authUC auth.Auth, log *slog.Logger) {
	r := &authRoutes{authUC, log}

	h := handler.Group("/translation")
	{
		h.GET("/history", r.history)
		h.POST("/do-translate", r.doTranslate)
	}
}

type historyResponse struct {
	History []entity.Translation `json:"history"`
}

func (r *authRoutes) history(c *gin.Context) {
	translations, err := r.t.History(c.Request.Context())
	if err != nil {
		r.l.Error(err, "http - v1 - history")
		errorResponse(c, http.StatusInternalServerError, "database problems")

		return
	}

	c.JSON(http.StatusOK, historyResponse{translations})
}

type doTranslateRequest struct {
	Source      string `json:"source"       binding:"required"  example:"auto"`
	Destination string `json:"destination"  binding:"required"  example:"en"`
	Original    string `json:"original"     binding:"required"  example:"текст для перевода"`
}

func (r *authRoutes) doTranslate(c *gin.Context) {
	var request doTranslateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		r.l.Error(err, "http - v1 - doTranslate")
		errorResponse(c, http.StatusBadRequest, "invalid request body")

		return
	}

	translation, err := r.t.Translate(
		c.Request.Context(),
		entity.Translation{
			Source:      request.Source,
			Destination: request.Destination,
			Original:    request.Original,
		},
	)
	if err != nil {
		r.l.Error(err, "http - v1 - doTranslate")
		errorResponse(c, http.StatusInternalServerError, "translation service problems")

		return
	}

	c.JSON(http.StatusOK, translation)
}
