package group

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

func (r *Repository) CreateGroup(ctx context.Context, name string) (uuid.UUID, error) {
	const query = `INSERT INTO users.groups (name) VALUES ($1) RETURNING uuid`

	var id uuid.UUID
	err := r.db.QueryRow(ctx, query, name).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create group: %w", err)
	}

	return id, nil
}

func (r *Repository) DeleteGroup(ctx context.Context, id uuid.UUID) error {
	const query = `DELETE FROM users.groups WHERE uuid = $1`

	tag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete group: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return repository.ErrNotFound
	}

	return nil
}

func (r *Repository) GetGroups(ctx context.Context) ([]entity.Group, error) {
	const query = `SELECT uuid, name FROM users.groups`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get groups: %w", err)
	}
	defer rows.Close()

	var groups []entity.Group
	for rows.Next() {
		var group entity.Group
		if err := rows.Scan(&group.UUID, &group.Name); err != nil {
			return nil, fmt.Errorf("failed to scan group: %w", err)
		}
		groups = append(groups, group)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return groups, nil
}

func (r *Repository) AddUsersToGroup(ctx context.Context, groupID uuid.UUID, userIDs []uuid.UUID) error {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	const query = `INSERT INTO users.group_users (group_uuid, user_uuid) VALUES ($1, $2) ON CONFLICT DO NOTHING`

	for _, userID := range userIDs {
		if _, err := tx.Exec(ctx, query, groupID, userID); err != nil {
			return fmt.Errorf("failed to add user to group: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *Repository) RemoveUsersFromGroup(ctx context.Context, groupID uuid.UUID, userIDs []uuid.UUID) error {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	const query = `DELETE FROM users.group_users WHERE group_uuid = $1 AND user_uuid = $2`

	for _, userID := range userIDs {
		if _, err := tx.Exec(ctx, query, groupID, userID); err != nil {
			return fmt.Errorf("failed to remove user from group: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *Repository) GetGroupUsers(ctx context.Context, groupID uuid.UUID) ([]uuid.UUID, error) {
	const query = `SELECT user_uuid FROM users.group_users WHERE group_uuid = $1`

	rows, err := r.db.Query(ctx, query, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to get group users: %w", err)
	}
	defer rows.Close()

	var userIDs []uuid.UUID
	for rows.Next() {
		var userID uuid.UUID
		if err := rows.Scan(&userID); err != nil {
			return nil, fmt.Errorf("failed to scan user uuid: %w", err)
		}
		userIDs = append(userIDs, userID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return userIDs, nil
}
