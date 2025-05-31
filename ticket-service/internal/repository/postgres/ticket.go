package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"ticket-service/internal/core"

	"github.com/google/uuid"
)

type TicketRepository struct {
	db *pgxpool.Pool
}

func NewTicketRepository(db *pgxpool.Pool) *TicketRepository {
	return &TicketRepository{
		db: db,
	}
}

func (r *TicketRepository) CreateTicket(ctx context.Context, ticket core.Ticket) error {
	query := `
		INSERT INTO ticket.tickets (
			uuid, name, description, status, created_on, updated_on, 
			created_by, recipient_type, recipient_uuid, responsible_uuid, priority
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
		)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		ticket.UUID,
		ticket.Name,
		ticket.Description,
		ticket.Status,
		ticket.CreatedOn,
		ticket.UpdatedOn,
		ticket.CreatedBy,
		ticket.RecipientType,
		ticket.RecipientUUID,
		ticket.ResponsibleUUID,
		ticket.Priority,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *TicketRepository) UpdateTicketStatus(ctx context.Context, ticketUUID uuid.UUID, status string) error {
	query := `
		UPDATE ticket.tickets 
		SET status = $1 
		WHERE uuid = $2
	`

	result, err := r.db.Exec(ctx, query, status, ticketUUID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("ticket not found")
	}

	return nil
}

func (r *TicketRepository) AssignResponsible(ctx context.Context, ticketUUID, responsibleUUID uuid.UUID) error {
	query := `
		UPDATE ticket.tickets 
		SET responsible_uuid = $1 
		WHERE uuid = $2
	`

	result, err := r.db.Exec(ctx, query, responsibleUUID, ticketUUID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("ticket not found")
	}

	return nil
}

func (r *TicketRepository) GetTickets(ctx context.Context, status string) ([]core.Ticket, error) {
	query := `
		SELECT 
			uuid, name, description, status, created_on, updated_on,
			created_by, recipient_type, recipient_uuid, responsible_uuid, priority
		FROM ticket.tickets
		ORDER BY created_on DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tickets []core.Ticket
	for rows.Next() {
		var ticket core.Ticket
		err := rows.Scan(
			&ticket.UUID,
			&ticket.Name,
			&ticket.Description,
			&ticket.Status,
			&ticket.CreatedOn,
			&ticket.UpdatedOn,
			&ticket.CreatedBy,
			&ticket.RecipientType,
			&ticket.RecipientUUID,
			&ticket.ResponsibleUUID,
			&ticket.Priority,
		)
		if err != nil {
			return nil, err
		}
		tickets = append(tickets, ticket)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tickets, nil
}
