package service

import (
	"auth-service/config"
	"auth-service/internal/repository"
	"auth-service/internal/service/auth"
	"auth-service/internal/service/hasher"
	"auth-service/internal/service/register"
)

type Manager struct {
	Auth     *auth.Service
	Hasher   *hasher.Service
	Register *register.Service
}

func New(repository repository.Manager, hashConfig *config.HashConfig) Manager {
	var manager Manager

	manager.Hasher = hasher.New(hashConfig.Cost, hashConfig.Salt)
	manager.Auth = auth.New(repository.Hydra, repository.AuthData, manager.Hasher)
	manager.Register = register.New(repository.AuthData, manager.Hasher)

	return manager
}
