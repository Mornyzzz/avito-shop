// Package v1 implements routing paths. Each services in own file.
package controller

import (
	h "avito-shop/internal/controller/handlers"
	"avito-shop/internal/controller/worker"
	repo "avito-shop/internal/repository"
	"avito-shop/internal/usecase/auth"
	"avito-shop/internal/usecase/buy"
	"avito-shop/internal/usecase/info"
	"avito-shop/internal/usecase/send"
	"avito-shop/pkg/postgres"
	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/pgxv4/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
	// Swagger docs.
	_ "github.com/evrone/go-clean-template/docs"
)

func NewRouter(handler *gin.Engine,
	log *slog.Logger,
	pg *postgres.Postgres,
	wp *worker.Pool,
) {
	// options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// use cases
	authUseCase := auth.New(
		repo.NewUserRepo(pg),
		repo.NewBalanceRepo(pg),
		manager.Must(trmsqlx.NewDefaultFactory(pg.Pool)),
	)

	buyUseCase := buy.New(
		repo.NewBalanceRepo(pg),
		repo.NewInventoryRepo(pg),
		manager.Must(trmsqlx.NewDefaultFactory(pg.Pool)),
	)

	infoUseCase := info.New(
		repo.NewBalanceRepo(pg),
		repo.NewInventoryRepo(pg),
		repo.NewTransactionRepo(pg),
		manager.Must(trmsqlx.NewDefaultFactory(pg.Pool)),
	)

	sendUseCase := send.New(
		repo.NewBalanceRepo(pg),
		repo.NewTransactionRepo(pg),
		manager.Must(trmsqlx.NewDefaultFactory(pg.Pool)),
	)

	// router
	v1 := handler.Group("/api")
	{
		h.NewAuthRoute(v1, authUseCase, wp, log)
		h.NewBuyRoute(v1, buyUseCase, wp, log)
		h.NewInfoRoute(v1, infoUseCase, wp, log)
		h.NewSendRoute(v1, sendUseCase, wp, log)
	}

}
