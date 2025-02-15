package app

import (
	"avito-shop/config"
	l "avito-shop/pkg/logger"
	"avito-shop/pkg/logger/handlers/slogpretty"
	"avito-shop/pkg/logger/sl"
	"avito-shop/pkg/postgres"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/exp/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	v1 "avito-shop/internal/controller/http/v1"
	repo "avito-shop/internal/repository"
	"avito-shop/internal/usecase/auth"
	"avito-shop/internal/usecase/buy"
	"avito-shop/internal/usecase/info"
	"avito-shop/internal/usecase/send"
	"github.com/evrone/go-clean-template/internal/usecase/webapi"
	"github.com/evrone/go-clean-template/pkg/httpserver"
	"github.com/evrone/go-clean-template/pkg/logger"
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
	v1.NewRouter(handler, log, translationUseCase)
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

func foo() {

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/url", func(r chi.Router) {
		r.Use(middleware.BasicAuth("avito-shop", map[string]string{
			cfg.HTTPServer.User: cfg.HTTPServer.Password,
		}))

		r.Post("/", save.New(log, storage))
		// TODO: add DELETE /url/{id}
	})

	router.Get("/{alias}", redirect.New(log, storage))

	log.Info("starting server", slog.String("address", cfg.Address))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("failed to start server")
		}
	}()

	log.Info("server started")

	<-done
	log.Info("stopping server")

	// TODO: move timeout to config
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to stop server", sl.Err(err))

		return
	}

	// TODO: close storage

	log.Info("server stopped")
}
