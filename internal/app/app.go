package app

import (
	"avito-shop/config"
	v1 "avito-shop/internal/controller/http/v1"
	repo "avito-shop/internal/repository"
	"avito-shop/internal/usecase/auth"
	"avito-shop/internal/usecase/buy"
	"avito-shop/internal/usecase/info"
	"avito-shop/internal/usecase/send"
	"avito-shop/pkg/httpserver"
	l "avito-shop/pkg/logger"
	_ "avito-shop/pkg/logger/handlers/slogpretty"
	"avito-shop/pkg/logger/sl"
	"avito-shop/pkg/postgres"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
	"os"
	"os/signal"
	"syscall"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	const op = "app.Run"

	//logger
	log := l.SetupLogger(cfg.Env)

	log.Info(
		"starting avito-shop",
		slog.String("env", cfg.Env),
		slog.String("version", cfg.Version),
	)
	log.Debug("debug messages are enabled")

	// Repository
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(-1)
	}
	defer pg.Close()

	// Use case
	authUseCase := auth.New(
		repo.NewUserRepo(pg),
		repo.NewBalanceRepo(pg),
	)

	buyUseCase := buy.New(
		repo.NewBalanceRepo(pg),
		repo.NewInventoryRepo(pg),
	)

	infoUseCase := info.New(
		repo.NewBalanceRepo(pg),
		repo.NewInventoryRepo(pg),
		repo.NewTransactionRepo(pg),
	)

	sendUseCase := send.New(
		repo.NewBalanceRepo(pg),
		repo.NewTransactionRepo(pg),
	)

	// HTTP Server
	handler := gin.New()
	v1.NewRouter(handler,
		log,
		authUseCase,
		buyUseCase,
		infoUseCase,
		sendUseCase,
	)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Info(op + s.String())
	case err = <-httpServer.Notify():
		log.Error(op, err)
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		log.Error("failed to stop server", fmt.Errorf("%s: %w", op, err))
	}

	log.Info("server stopped")

}
