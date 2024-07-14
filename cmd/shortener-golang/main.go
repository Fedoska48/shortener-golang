package main

import (
	"log/slog"
	"os"
	"shortener-golang/internal/config"
	"shortener-golang/internal/http-server/logger"
	"shortener-golang/internal/lib/logger/sl"
	"shortener-golang/internal/storage/sqlite"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// TODO: init config: clearenv

	cfg := config.MustLoad()

	// TODO: init logger: slog: log/slog
	log := setupLogger(cfg.Env)

	log.Info("Runnning Shortener", slog.String("env", cfg.Env))

	// TODO: init storage: sqlite
	storage, err := sqlite.New(cfg.StoragePath)

	_ = storage
	if err != nil {
		log.Error("fail to init storage", sl.Err(err))
		os.Exit(1)
		// or empty return
	}

	router := chi.NewRouter()
	// get request id
	router.Use(middleware.RequestID)
	// request logger
	router.Use(logger.New(log))
	// recover panic
	router.Use(middleware.Recoverer)
	// get url params
	router.Use(middleware.URLFormat)

	// TODO: init run server
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
