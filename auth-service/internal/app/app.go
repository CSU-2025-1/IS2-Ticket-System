package app

import (
	"auth-service/cmd/migrator"
	"auth-service/config"
	"auth-service/internal/http"
	"auth-service/internal/http/handler"
	"auth-service/internal/repository"
	"auth-service/internal/service"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
)

type App struct {
	address    string
	router     *gin.Engine
	controller *handler.Controller
}

func Build(ctx context.Context, cfg config.Config) (*App, error) {
	if err := runMigration(ctx, cfg); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	repositories, err := repository.Init(ctx, cfg.Database, cfg.Hydra)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize repositories: %w", err)
	}

	services := service.New(repositories, cfg.Hash)

	controller := &handler.Controller{
		Repository: repositories,
		Service:    services,
	}

	router := http.SetupRouter(controller)

	return &App{
		address:    cfg.Server.Address,
		router:     router,
		controller: controller,
	}, nil
}

func runMigration(ctx context.Context, cfg config.Config) error {
	err := migrator.MigratePostgres(ctx, *cfg.Database)
	if err != nil {
		return err
	}

	err = migrator.MigrateHydra(ctx, *cfg.Hydra)
	if err != nil {
		return err
	}

	return nil
}

func (a App) Run() error {
	if err := a.router.Run(a.address); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

func (a App) Shutdown(ctx context.Context) {
	a.controller.Repository.Close()
}
