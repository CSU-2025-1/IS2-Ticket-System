package rabbitmq

import (
	"context"
	"github.com/google/uuid"
	"ticket-service/internal/core"
	"ticket-service/pkg/rabbitmq"
)

type Ticket struct {
	ID              uuid.UUID `json:"id"`
	Type            string    `json:"type" `
	Title           string    `json:"title"`
	Priority        int       `json:"priority"`
	ResponsibleUUID uuid.UUID `json:"responsible_uuid"`
}

type RabbitMQ struct {
	writer *rabbitmq.Writer[Ticket]
}

func NewRabbitMQ(w *rabbitmq.Writer[Ticket]) (*RabbitMQ, error) {
	return &RabbitMQ{
		writer: w,
	}, nil
}

func (r *RabbitMQ) Save(ctx context.Context, ticket core.Ticket) error {
	return r.writer.Send(ctx, Ticket{
		ID:              ticket.UUID,
		Type:            ticket.RecipientType,
		Title:           ticket.Name,
		Priority:        ticket.Priority,
		ResponsibleUUID: ticket.ResponsibleUUID,
	})
}
