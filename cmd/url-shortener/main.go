package main

import (
	"fmt"
	"os"

	"log/slog"

	"github.com/YuryFilipovich/url-shortener/internal/config"
	"github.com/YuryFilipovich/url-shortener/internal/lib/logger/sl"
	"github.com/YuryFilipovich/url-shortener/internal/storage/sqlite"

	mwLogger "github.com/YuryFilipovich/url-shortener/internal/http-server/middleware/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg)

	log := setupLogger(cfg.Env)
	log.Info("starting url-shortener", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	log.Info(os.Getwd())

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	err = storage.DeleteURL("0")
	if err != nil {
		log.Error("%s", err)
		os.Exit(1)
	}
	log.Info("%s", storage)

	_ = storage

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.RealIP)
	router.Use(mwLogger.New(log))
	// TODO: run server

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
}
