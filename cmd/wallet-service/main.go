package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"wallet-service/internal/config"
	"wallet-service/internal/http-server/handlers/wallet/save"
	"wallet-service/internal/http-server/middleware/logger"
	"wallet-service/internal/lib/logger/sl"
	"wallet-service/internal/services/wallet"
	"wallet-service/internal/storage/postgres"
	pgxdriver "wallet-service/pkg/pgx-driver"
	"wallet-service/pkg/pgx-driver/transaction"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg := config.MustLoad()

	log := sl.InitLogger(cfg.Env, os.Stdout)

	log.Debug("CONFIG", cfg)

	pCfg := cfg.Storage.Postgres

	storage, err := pgxdriver.New(
		cfg.Storage.Postgres.DSN,
		log,
		pgxdriver.MaxPoolSize(pCfg.MaxOpenConns),
		pgxdriver.MinConns(pCfg.MaxIdleConns),
		pgxdriver.MaxConnIdleTime(pCfg.MaxIdleTime))

	if err != nil {
		panic(err)
	}

	txManger, err := transaction.NewManager(storage, log)
	if err != nil {
		panic(err)
	}

	walletRepository := postgres.NewWalletRepository(log, storage)
	operationRepository := postgres.NewOperationRepository(log, storage)

	walletService := wallet.New(
		txManger,
		log,
		walletRepository,
		walletRepository,
		walletRepository,
		operationRepository)

	router := chi.NewRouter()

	// middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(logger.NewLoggerMiddleware(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/api/v1", func(r chi.Router) {
		r.Post("/wallet", save.New(log, walletService))
	})

	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		s := <-quit

		log.Info("caught signal", map[string]string{
			"signal": s.String(),
		})

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		shutdownError <- srv.Shutdown(ctx)
	}()

	log.Info("starting server wallet-service",
		slog.String("env", cfg.Env),
		slog.String("port", cfg.HTTPServer.Address),
	)

	err = srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}

	<-shutdownError

	err = srv.Close()
	if err != nil {
		log.Error("failed stop server", sl.Err(err))
		return
	}

	storage.Close()

	log.Info("stopped server", map[string]string{
		"addr": srv.Addr,
	})

}
