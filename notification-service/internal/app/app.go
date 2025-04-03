package app

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"notification-service/config"
	"notification-service/internal/repository"
	"notification-service/internal/service"
	"notification-service/internal/transport/http"
	"notification-service/internal/transport/kafka"
	consulPkg "notification-service/pkg/consul"
	"notification-service/pkg/grpc"
	loggerPkg "notification-service/pkg/logger"
	"sync"
)

type App struct {
	config *config.Config
}

func New(config *config.Config) *App {
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
	registeredUUID, err := consulClient.Register(
		"public-notification-service",
		a.config.Http.Address,
		a.config.Http.Port,
	)
	if err != nil {
		return err
	}
	defer func() {
		if err := consulClient.Deregister(registeredUUID); err != nil {
			appLogger.Errorf("ошибка разрегистрации в Consul Registry")
		}
	}()

	appLogger.Infof("успешная регистрация в Consul Registry")

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
			"private-user-manager",
		),
	)
	receiverRepository := repository.NewReceiverRepository(pool)

	mailService := service.NewMailService(a.config.Mail)

	repositoryManager := &repository.RepositoryManager{
		ReceiverRepository: receiverRepository,
		UserRepository:     userRepository,
	}

	serviceManager := &service.ServiceManager{
		MailService: mailService,
		Logger:      appLogger,
	}

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := http.New(
			a.config.Http,
			repositoryManager,
			serviceManager,
		).Run(); err != nil {
			appLogger.Errorf("ошибка запуска http сервера")
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := kafka.New(
			a.config.Kafka,
			repositoryManager,
			serviceManager,
		).Run(ctx); err != nil {
			appLogger.Errorf("ошибка запуска kafka сервера")
		}

	}()

	appLogger.Infof("приложение запущено")
	wg.Wait()

	appLogger.Infof("приложение остановлено")
	return nil
}
