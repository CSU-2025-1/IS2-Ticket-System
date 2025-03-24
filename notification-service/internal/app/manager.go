package app

import (
	"notification-service/internal/repository"
	"notification-service/internal/service"
	"notification-service/pkg/logger"
)

type RepositoryManager struct {
	ReceiverRepository *repository.ReceiverRepository
	UserRepository     *repository.UserRepository
}

type ServiceManager struct {
	MailService *service.MailService
	Logger      *logger.PrettyStdout
}
