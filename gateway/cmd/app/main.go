package main

import (
	"context"
	"gateway/config"
	"gateway/internal/application"
	"log"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if err := config.LoadDotEnv(); err != nil {
		log.Printf("continue working, but failed to load .env file: %s", err.Error())
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %s", err.Error())
	}

	if err = application.New(cfg).Run(ctx); err != nil {
		log.Fatalf("failed to run application: %s", err.Error())
	}
}
