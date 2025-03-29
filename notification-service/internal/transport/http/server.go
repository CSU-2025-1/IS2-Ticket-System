package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"notification-service/config"
	"notification-service/internal/repository"
	"notification-service/internal/service"
	"notification-service/internal/usecase"
)

type Server struct {
	config            config.Http
	repositoryManager *repository.RepositoryManager
	serviceManager    *service.ServiceManager
}

func New(
	config config.Http,
	repositoryManager *repository.RepositoryManager,
	serviceManager *service.ServiceManager,
) *Server {
	return &Server{
		repositoryManager: repositoryManager,
		serviceManager:    serviceManager,
		config:            config,
	}
}

func (s *Server) Run() error {
	router := gin.Default()

	registerMailReceiverUseCase := usecase.NewRegisterMailReceiverUseCase(
		s.serviceManager.Logger,
		s.repositoryManager.ReceiverRepository,
	)

	router.POST("/api/notification/mail/register", RegisterMailReceiver(registerMailReceiverUseCase))
	router.GET("/check", HealthCheck())

	return router.Run(fmt.Sprintf(":%d", s.config.Port))
}
