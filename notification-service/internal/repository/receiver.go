package repository

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"notification-service/internal/model"
)

type ReceiverRepository struct {
	db *pgxpool.Pool
}

func NewReceiverRepository(db *pgxpool.Pool) *ReceiverRepository {
	return &ReceiverRepository{
		db: db,
	}
}

func (r *ReceiverRepository) CreateMailReceiver(ctx context.Context, userUUID uuid.UUID, mail string) error {
	sql := `INSERT INTO notification.mail_receiver VALUES($1, $2)`
	_, err := r.db.Exec(ctx, sql, userUUID, mail)
	return err
}

func (r *ReceiverRepository) GetAllMailReceiversByUUIDs(ctx context.Context, ids []uuid.UUID) ([]model.Receiver, error) {
	sql := `SELECT id, mail FROM notification.mail_receiver WHERE id = ANY($1)`
	rows, err := r.db.Query(ctx, sql, ids)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []model.Receiver{}, nil
		}

		return nil, err
	}

	return pgx.CollectRows(rows, pgx.RowToStructByNameLax[model.Receiver])
}
