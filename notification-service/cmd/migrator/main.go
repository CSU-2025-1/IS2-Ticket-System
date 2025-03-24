package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"notification-service/config"
	"os"
	"os/signal"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt)
	defer cancel()

	_ = config.LoadDotEnv()

	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	db, err := pgxpool.New(ctx, cfg.Postgres.ConnectionString)
	if err != nil {
		panic(err)
	}

	connection := stdlib.OpenDBFromPool(db)
	if connection == nil {
		panic(err)
	}

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	if err := goose.Up(connection, "././migrations"); err != nil {
		panic(err)
	}
}
