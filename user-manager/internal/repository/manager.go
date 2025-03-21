package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"user-mananger/internal/repository/postgres/group"
	"user-mananger/internal/repository/postgres/user"
)

type Manager struct {
	User  *user.Repository
	Group *group.Repository
}

func NewManager(db *pgxpool.Pool) *Manager {
	m := &Manager{
		User:  user.New(db),
		Group: group.New(db),
	}

	return m
}
