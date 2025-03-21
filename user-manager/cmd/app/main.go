package main

import (
	"context"
	"log"
	"user-mananger/config"
	"user-mananger/internal/repository/postgres"
)

func main() {
	cfg, err := config.Parse("config/config.yaml")
	if err != nil {
		log.Fatal("parse config file error", err.Error())
	}

	ctx := context.Background()

	db, err := postgres.Connect(ctx, cfg.Database)
	if err != nil {
		log.Fatal("connect database error", err.Error())
	}
	defer db.Close()
}
