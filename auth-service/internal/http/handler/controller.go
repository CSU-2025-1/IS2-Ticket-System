package handler

import (
	"auth-service/internal/repository"
	"auth-service/internal/service"
)

type Controller struct {
	Repository repository.Manager
	Service    service.Manager
}
