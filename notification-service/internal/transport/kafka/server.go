package kafka

import (
	"context"
	"notification-service/config"
	"notification-service/internal/repository"
	"notification-service/internal/service"
	"notification-service/internal/usecase"
	kafkaPkg "notification-service/pkg/kafka"
)

type Server struct {
	config            config.Kafka
	repositoryManager *repository.RepositoryManager
	serviceManager    *service.ServiceManager
}

func New(
	config config.Kafka,
	repositoryManager *repository.RepositoryManager,
	serviceManager *service.ServiceManager) *Server {
	return &Server{
		config:            config,
		repositoryManager: repositoryManager,
		serviceManager:    serviceManager,
	}
}

func (s *Server) Run(ctx context.Context) error {
	notificationByTicketUseCase := usecase.NewNotificationByTicketUseCase(
		s.serviceManager.Logger,
		s.repositoryManager.ReceiverRepository,
		s.repositoryManager.UserRepository,
		s.serviceManager.MailService,
	)
	ticketChan := make(chan usecase.Ticket, 100)
	ticketCreationConsumer := kafkaPkg.New(
		kafkaPkg.ConsumerConfig{
			Address:         s.config.Address,
			GroupID:         s.config.GroupID,
			AutoOffsetReset: s.config.AutoOffsetReset,
			Topic:           s.config.TicketCreationTopic,
		},
		s.serviceManager.Logger,
	)
	if err := ticketCreationConsumer.Init(); err != nil {
		return err
	}

	go func() {
		if err := kafkaPkg.RunConsumer[usecase.Ticket](
			ticketCreationConsumer,
			s.serviceManager.Logger,
			ticketChan,
		); err != nil {
			s.serviceManager.Logger.Errorf("Kafka.Server.Run: %s", err.Error())
		}
	}()

	for {
		select {
		case ticket, ok := <-ticketChan:
			if !ok {
				continue
			}

			go func() {
				if err := notificationByTicketUseCase.Execute(ctx, ticket); err != nil {
					s.serviceManager.Logger.Errorf("ошибка выполнения юзкейса: %s", err.Error())
				}
			}()
		case <-ctx.Done():
			s.serviceManager.Logger.Infof("контект отменен")
			return nil
		}
	}
}
