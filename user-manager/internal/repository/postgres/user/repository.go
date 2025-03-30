package user

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"user-mananger/internal/domain/entity"
	"user-mananger/internal/domain/errors/repository"
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreateUser(ctx context.Context, user entity.User) (uuid.UUID, error) {
	const query = `INSERT INTO users.users (login) VALUES ($1) RETURNING uuid`

	var id uuid.UUID
	err := r.db.QueryRow(ctx, query, user.Login).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create user: %w", err)
	}

	return id, nil
}

func (r *Repository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	const query = `DELETE FROM users.users WHERE uuid = $1`

	tag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return repository.ErrNotFound
	}

	return nil
}

func (r *Repository) GetUsers(ctx context.Context) ([]entity.User, error) {
	const query = `SELECT uuid, login FROM users.users`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	defer rows.Close()

	var users []entity.User
	for rows.Next() {
		var user entity.User
		if err := rows.Scan(&user.UUID, &user.Login); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return users, nil
}
