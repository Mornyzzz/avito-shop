package app

import (
	"fmt"
	"os"
	"os/signal"
	_ "sync"
	"syscall"

	_ "github.com/avito-tech/go-transaction-manager/drivers/pgxv4/v2"
	_ "github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"

	"avito-shop/config"
	"avito-shop/internal/controller"
	"avito-shop/internal/controller/worker"
	_ "avito-shop/internal/repository"
	"avito-shop/pkg/httpserver"
	l "avito-shop/pkg/logger"
	_ "avito-shop/pkg/logger/handlers/slogpretty"
	"avito-shop/pkg/logger/sl"
	"avito-shop/pkg/postgres"
)

const (
	numWorkers = 18
	taskNum    = 18
)

func Run(cfg *config.Config) {
	const op = "app.Run"

	// logger
	log := l.SetupLogger(cfg.Env)

	log.Info(
		"starting avito-shop",
		slog.String("env", cfg.Env),
		slog.String("version", cfg.Version),
	)
	log.Debug("debug messages are enabled")

	// Repository
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PoolMax))
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(-1)
	}
	defer pg.Close()

	// Workers
	workerPool := worker.NewWorkerPool(numWorkers, taskNum)
	defer workerPool.Shutdown()

	// HTTP Server
	handler := gin.New()
	controller.NewRouter(handler,
		log,
		pg,
		workerPool,
	)

	// run server
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
	workerPool.Shutdown()

	err = httpServer.Shutdown()
	if err != nil {
		log.Error("failed to stop server", fmt.Errorf("%s: %w", op, err))
	}

	log.Info("server stopped")
}
