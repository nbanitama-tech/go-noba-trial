package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bytedance/go-noba-trial/internal/config"
	"github.com/bytedance/go-noba-trial/internal/datastore"
	httpapi "github.com/bytedance/go-noba-trial/internal/delivery/http"
	"github.com/bytedance/go-noba-trial/internal/usecase"
)

func main() {
	cfg := config.Load()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))
	slog.SetDefault(logger)

	dbCtx, dbCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer dbCancel()

	db, err := datastore.OpenPostgres(dbCtx, cfg.DatabaseURL)
	if err != nil {
		logger.Error("failed to connect to postgres", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := datastore.EnsureSchema(dbCtx, db); err != nil {
		logger.Error("failed to ensure postgres schema", "error", err)
		os.Exit(1)
	}

	healthUsecase := usecase.NewHealthUsecase(cfg.ServiceName)
	userRepository := datastore.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepository)
	router := httpapi.NewRouter(httpapi.RouterConfig{
		HealthUsecase: healthUsecase,
		UserUsecase:   userUsecase,
		Logger:        logger,
		BearerToken:   cfg.BearerToken,
	})

	server := &http.Server{
		Addr:              cfg.HTTPAddress,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		logger.Info("http server started", "address", cfg.HTTPAddress)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("http server failed", "error", err)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("http server shutdown failed", "error", err)
		os.Exit(1)
	}

	logger.Info("http server stopped")
}
