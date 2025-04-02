package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/segmentio/kafka-go"
	auth_data "user-mananger/internal/repository/kafka/auth-data"
	"user-mananger/internal/repository/postgres/group"
	"user-mananger/internal/repository/postgres/user"
)

type Manager struct {
	User     *user.Repository
	Group    *group.Repository
	AuthData *auth_data.Repository
}

func NewManager(db *pgxpool.Pool, kafka *kafka.Conn) *Manager {
	m := &Manager{
		User:     user.New(db),
		Group:    group.New(db),
		AuthData: auth_data.New(kafka),
	}

	return m
}
