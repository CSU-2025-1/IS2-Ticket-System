package rabbitmq

import (
	"context"
	"github.com/google/uuid"
	"notification-service/internal/repository"
	"notification-service/internal/service"
	"notification-service/internal/usecase"
	"notification-service/pkg/rabbitmq"
)

type Ticket struct {
	ID              uuid.UUID `json:"id"`
	Type            string    `json:"type" `
	Title           string    `json:"title"`
	Priority        int       `json:"priority"`
	ResponsibleUUID uuid.UUID `json:"responsible_uuid"`
}

type Server struct {
	repositoryManager *repository.RepositoryManager
	serviceManager    *service.ServiceManager
	rabbitConfig      rabbitmq.Config
}

func New(repositoryManager *repository.RepositoryManager, serviceManager *service.ServiceManager, rabbitConfig rabbitmq.Config) *Server {
	return &Server{
		repositoryManager: repositoryManager,
		serviceManager:    serviceManager,
		rabbitConfig:      rabbitConfig,
	}
}

func (s *Server) Run(ctx context.Context) error {
	reader, err := rabbitmq.CreateJsonReader[Ticket](s.rabbitConfig, s.handler, rabbitmq.DefaultConsumeOption)
	if err != nil {
		return err
	}

	err = reader.Run(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) handler(ticket Ticket) {
	notificationByTicketUseCase := usecase.NewNotificationByTicketUseCase(
		s.serviceManager.Logger,
		s.repositoryManager.ReceiverRepository,
		s.repositoryManager.UserRepository,
		s.serviceManager.MailService,
	)

	s.serviceManager.Logger.Infof("получен тикет")

	if err := notificationByTicketUseCase.Execute(context.Background(), usecase.Ticket(ticket)); err != nil {
		s.serviceManager.Logger.Errorf("ошибка выполнения юзкейса: %s", err.Error())
	}
}
