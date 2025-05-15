package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"user-mananger/internal/repository/postgres/group"
	"user-mananger/internal/repository/postgres/user"
	auth_data "user-mananger/internal/repository/rabbitmq/auth-data"
	"user-mananger/pkg/rabbitmq"
)

type Manager struct {
	User     *user.Repository
	Group    *group.Repository
	AuthData *auth_data.Repository
}

func NewManager(db *pgxpool.Pool, cfg rabbitmq.Config) *Manager {
	writer, err := rabbitmq.CreateJsonWriter[auth_data.AuthData](cfg, rabbitmq.DefaultPublishOption)
	if err != nil {
		panic(err.Error())
	}

	m := &Manager{
		User:     user.New(db),
		Group:    group.New(db),
		AuthData: auth_data.NewRepository(writer),
	}

	return m
}
