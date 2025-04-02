package user

import (
	"auth-service/internal/domain/entity"
	"auth-service/internal/domain/errors/repository"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r Repository) GetUserByLogin(ctx context.Context, login string) (*entity.User, error) {
	query := `select (uuid, login, password) from auth.users_auth_data where login = $1`

	var res AuthData
	err := r.db.QueryRow(ctx, query, login).Scan(&res)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}

		return nil, fmt.Errorf("user repo: get user by login: %w", err)
	}
	return res.ToEntity(), nil
}

func (r Repository) CreateUser(ctx context.Context, user *entity.User) error {
	query := `insert into auth.users_auth_data(uuid, login, password) values ($1, $2, $3)`

	_, err := r.db.Exec(ctx, query, user.UUID, user.Login, user.Password)
	if err != nil {
		return fmt.Errorf("user repo: create user: %w", err)
	}

	return nil
}
