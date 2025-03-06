package repository

import (
	"auth-service/internal/repository/postgres"
	"auth-service/internal/repository/postgres/user"
	"auth-service/pkg/hydra"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Manager struct {
	pgPool *pgxpool.Pool

	AuthData *user.Repository

	Hydra *hydra.Client
}

func Init(ctx context.Context, pgConfig *postgres.Config, hydraConfig *hydra.Config) (Manager, error) {
	var manager Manager
	var err error

	manager.pgPool, err = postgres.Connect(ctx, pgConfig)
	if err != nil {
		return Manager{}, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	manager.Hydra = hydra.New(ctx, hydraConfig)
	manager.AuthData = user.New(manager.pgPool)

	return manager, nil
}

func (manager Manager) Close() {
	manager.pgPool.Close()
}
