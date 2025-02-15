// Package v1 implements routing paths. Each services in own file.
package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/exp/slog"
	"net/http"

	"avito-shop/internal/usecase/auth"
	"avito-shop/internal/usecase/buy"
	"avito-shop/internal/usecase/info"
	"avito-shop/internal/usecase/send"
	"avito-shop/pkg/logger"
	// Swagger docs.
	_ "github.com/evrone/go-clean-template/docs"
)

// NewRouter -.
// Swagger spec:
// @title       Go Clean Template API
// @description Using a translation service as an example
// @version     1.0
// @host        localhost:8080
// @BasePath    /v1

func NewRouter(handler *gin.Engine,
	log *slog.Logger,
	authUC auth.Auth,
	buyUC buy.Buy,
	infoUC info.Info,
	sendUC send.Send,
) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// Swagger
	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	handler.GET("/swagger/*any", swaggerHandler)

	// K8s probe
	handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Prometheus metrics
	handler.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Routers
	v1 := handler.Group("/v1")
	{
		newSendRoutes(v1, sendUC, log)
		newAuthRoutes(v1, authUC, log)
		newInfoRoutes(v1, infoUC, log)
		newBuyRoutes(v1, buyUC, log)
	}
}
