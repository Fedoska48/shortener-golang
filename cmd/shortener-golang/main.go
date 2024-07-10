package main

import (
	"log/slog"
	"os"
	"shortener-golang/internal/config"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// TODO: init config: clearenv

	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("Runnning Shortener", slog.String("env", cfg.Env))
	log.Debug("Debug level")
	// TODO: init logger: slog: log/slog

	// TODO: init storage: sqlite

	// TODO: init router: chi, chi render

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
