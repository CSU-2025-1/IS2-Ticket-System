package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"notification-service/config"
	"notification-service/internal/app"
	"notification-service/internal/usecase"
)

type Server struct {
	config            config.Http
	repositoryManager *app.RepositoryManager
	serviceManager    *app.ServiceManager
}

func New(
	config config.Http,
	repositoryManager *app.RepositoryManager,
	serviceManager *app.ServiceManager,
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

	router.POST("/api/mail/register", RegisterMailReceiver(registerMailReceiverUseCase))

	return router.Run(fmt.Sprintf(":%d", s.config.Port))
}
