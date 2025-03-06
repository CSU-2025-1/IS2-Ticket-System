package service

import (
	"auth-service/config"
	"auth-service/internal/repository"
	"auth-service/internal/service/auth"
	"auth-service/internal/service/hasher"
)

type Manager struct {
	Auth   *auth.Service
	Hasher *hasher.Service
}

func New(repository repository.Manager, hashConfig *config.HashConfig) Manager {
	var manager Manager

	manager.Auth = auth.New(repository.Hydra, repository.AuthData)
	manager.Hasher = hasher.New(hashConfig.Cost, hashConfig.Salt)

	return manager
}
