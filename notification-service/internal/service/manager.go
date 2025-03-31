package service

import "notification-service/pkg/logger"

type ServiceManager struct {
	MailService *MailService
	Logger      *logger.PrettyStdout
}
