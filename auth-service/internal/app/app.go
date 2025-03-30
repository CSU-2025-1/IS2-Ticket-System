package app

import (
	"auth-service/cmd/migrator"
	"auth-service/config"
	"auth-service/internal/consumer"
	"auth-service/internal/grpc"
	"auth-service/internal/http"
	"auth-service/internal/http/handler"
	"auth-service/internal/repository"
	"auth-service/internal/service"
	"auth-service/pkg/consul"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net"
)

type App struct {
	address string
	router  *gin.Engine

	consul *consul.Client
	uuid   string

	grpsAddress string
	controller  *handler.Controller

	userConsumer *consumer.Consumer

	grpcServer *grpc.Server
}

func Build(ctx context.Context, cfg config.Config) (*App, error) {
	if err := runMigration(ctx, cfg); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	consulClient := consul.New(*cfg.Consul)
	serviceUUID, err := consulClient.Register("auth", "auth", 8080)
	if err != nil {
		return nil, fmt.Errorf("failed to register auth service: %w", err)
	}

	repositories, err := repository.Init(ctx, cfg.Database, cfg.Hydra, cfg.Kafka)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize repositories: %w", err)
	}

	services := service.New(repositories, cfg.Hash)

	controller := &handler.Controller{
		Repository: repositories,
		Service:    services,
	}

	userConsumer := consumer.NewConsumer(services.Register, repositories.UserStream)

	grpcServer := grpc.NewServer(repositories.Hydra)

	router := http.SetupRouter(controller)

	return &App{
		controller: controller,

		consul: consulClient,
		uuid:   serviceUUID,

		address: cfg.Server.Address,
		router:  router,

		grpsAddress: cfg.Grpc.Address,
		grpcServer:  grpcServer,

		userConsumer: userConsumer,
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

func (a App) Run(ctx context.Context) error {
	errChan := make(chan error)

	go func(errChan chan error) {
		listener, err := net.Listen("tcp", a.grpsAddress)
		if err != nil {
			errChan <- fmt.Errorf("failed to listen: %w", err)
		}

		err = a.grpcServer.Run(listener)
		if err != nil {
			errChan <- fmt.Errorf("failed to start grpc server: %w", err)
		}
	}(errChan)

	go func() {
		if err := a.router.Run(a.address); err != nil {
			errChan <- fmt.Errorf("failed to start http server: %w", err)
		}
	}()

	go a.userConsumer.RunUserConsuming(ctx)

	return <-errChan
}

func (a App) Shutdown(ctx context.Context) {
	a.controller.Repository.Close()
	_ = a.consul.Deregister(a.uuid)
}
