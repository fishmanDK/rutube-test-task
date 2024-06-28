package main

import (
	"context"
	"fmt"
	"github.com/fishmanDK/rutube-test-task/config"
	"github.com/fishmanDK/rutube-test-task/internal/handlers"
	"github.com/fishmanDK/rutube-test-task/internal/service"
	"github.com/fishmanDK/rutube-test-task/internal/storage"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	envLocal = "local"
	envDev   = "dev"
)

func main() {
	cfg := config.MustLoadConfig()
	cfgDB := config.MustLoadConfigDB()
	logger := setupLogger(cfg.Env)

	db, err := storage.MustStorage(cfgDB)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	srvc := service.MustService(db)
	handl := handlers.MustHandlers(srvc, logger)
	router := handl.InitRouts()

	server := &http.Server{
		Addr:         cfg.HttpServer.Addr,
		Handler:      router,
		ReadTimeout:  cfg.HttpServer.Timeout,
		WriteTimeout: cfg.HttpServer.Timeout,
		IdleTimeout:  cfg.HttpServer.IdleTimeout,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server error", slog.String("err", err.Error()))
			os.Exit(1)
		}
	}()

	logger.Info("Rutube-test-task Started")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	sig := <-stop
	fmt.Printf("Received signal: %v\n", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server shutdown error", slog.String("err", err.Error()))
		os.Exit(1)
	}

	logger.Info("Server gracefully stopped")
}

func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case envLocal:
		opts := &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}
		slogHandler := slog.NewTextHandler(os.Stdout, opts)
		logger = slog.New(slogHandler)
	case envDev:
		opts := &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}
		slogHandler := slog.NewJSONHandler(os.Stdout, opts)
		logger = slog.New(slogHandler)
	}

	return logger
}
