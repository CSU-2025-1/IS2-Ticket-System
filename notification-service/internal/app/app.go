package app

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"notification-service/config"
	"notification-service/internal/repository"
	"notification-service/internal/service"
	"notification-service/internal/transport/http"
	consulPkg "notification-service/pkg/consul"
	"notification-service/pkg/grpc"
	loggerPkg "notification-service/pkg/logger"
	"sync"
)

type App struct {
	config config.Config
}

func New(config config.Config) *App {
	return &App{
		config: config,
	}
}

func (a *App) Run(ctx context.Context) error {
	appLogger := loggerPkg.NewPrettyStdout(loggerPkg.Debug)
	consulClient := consulPkg.New(consulPkg.Config(a.config.Consul))
	if err := consulClient.Configure(); err != nil {
		return err
	}

	pool, err := pgxpool.New(ctx, a.config.Postgres.ConnectionString)
	if err != nil {
		return err
	}

	if err := pool.Ping(ctx); err != nil {
		return err
	}

	userRepository := repository.NewUserRepository(
		grpc.New(
			appLogger,
			consulClient,
			"public-user-manager",
		),
	)
	receiverRepository := repository.NewReceiverRepository(pool)

	mailService := service.NewMailService(a.config.Mail)

	repositoryManager := &RepositoryManager{
		ReceiverRepository: receiverRepository,
		UserRepository:     userRepository,
	}

	serviceManager := &ServiceManager{
		MailService: mailService,
		Logger:      appLogger,
	}

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		wg.Done()
		if err := http.New(
			a.config.Http,
			repositoryManager,
			serviceManager,
		).Run(); err != nil {
			appLogger.Errorf("ошибка запускка http сервера")
		}
	}()

	go func() {

	}()

	wg.Wait()

	return nil
}
