package main

import (
	"context"
	"log"
	"user-mananger/config"
	"user-mananger/internal/http"
	"user-mananger/internal/repository"
	"user-mananger/internal/repository/kafka"
	"user-mananger/internal/repository/postgres"
	"user-mananger/pkg/consul"
)

func main() {
	cfg, err := config.Parse("config/config.yaml")
	if err != nil {
		log.Fatal("parse config file error", err.Error())
	}

	consulClient := consul.New(*cfg.Consul)
	if err := consulClient.Configure(); err != nil {
		log.Fatalf("failed to configure consul client: %s", err.Error())
	}

	serviceUUID, err := consulClient.Register("public-user-manager", "user-manager", 8080)
	if err != nil {
		log.Fatal("register service error", err.Error())
	}
	defer consulClient.Deregister(serviceUUID)

	ctx := context.Background()

	db, err := postgres.Connect(ctx, cfg.Database)
	if err != nil {
		log.Fatal("connect database error", err.Error())
	}
	defer db.Close()

	broker, err := kafka.Connect(ctx, *cfg.Kafka)
	if err != nil {
		log.Fatal("connect kafka error", err.Error())
	}

	repositories := repository.NewManager(db, broker)

	r := http.SetupRouter(*repositories)

	if err = r.Run(cfg.Server.Address); err != nil {
		log.Fatal("start http server error", err.Error())
	}
}
