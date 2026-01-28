package main

import (
	"log/slog"
	"os"

	"hitalent_tt/internal/api"
	"hitalent_tt/internal/config"
	"hitalent_tt/internal/storage/psql"
)

const logLevelInfo = "info"
const logLevelDebug = "debug"

func main() {

	cfg := config.NewConfig()

	log := setLogger(cfg.LogLevel)

	log.Info("starting hitalent_tt service...")

	storage, err := psql.New(cfg.StoragePath)
	if err != nil {
		log.Error("could not connect to storage", err)
		os.Exit(1)
	}

	srv := api.NewAPI(storage, cfg, log)
	srv.Start()
}

func setLogger(logLevel string) *slog.Logger {

	var log *slog.Logger

	switch logLevel {
	case logLevelInfo:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	case logLevelDebug:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}

	return log
}
