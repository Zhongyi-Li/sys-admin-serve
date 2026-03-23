package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"sys-admin-serve/internal/bootstrap"

	"go.uber.org/zap"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	configPath := os.Getenv("APP_CONFIG")
	if configPath == "" {
		configPath = "configs/config.yaml"
	}

	app, err := bootstrap.New(configPath)
	if err != nil {
		log.Fatalf("bootstrap application: %v", err)
	}

	defer func() {
		_ = app.Logger.Sync()
	}()

	if err := app.Run(ctx); err != nil {
		app.Logger.Fatal("application exited", zap.Error(err))
	}
}
