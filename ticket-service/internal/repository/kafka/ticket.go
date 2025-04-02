package kafka

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"ticket-service/internal/core"
)

type Ticket struct {
	conn *kafka.Conn
}

func NewTicketSaver(conn *kafka.Conn) *Ticket {
	return &Ticket{
		conn: conn,
	}
}

type ticketMessage struct {
	ID              uuid.UUID
	Type            string
	Title           string
	Priority        int
	ResponsibleUUID uuid.UUID
}

func (t *Ticket) SendMessageTicketCreation(ticket core.Ticket) error {
	msg := ticketMessage{
		ID:              ticket.UUID,
		Type:            ticket.RecipientType,
		Title:           ticket.Name,
		Priority:        ticket.Priority,
		ResponsibleUUID: ticket.RecipientUUID,
	}

	ticketJson, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	_, err = t.conn.Write(ticketJson)
	if err != nil {
		return err
	}

	return nil
}
