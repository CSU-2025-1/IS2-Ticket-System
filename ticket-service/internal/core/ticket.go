package core

import (
	"github.com/google/uuid"
	"time"
)

type Ticket struct {
	UUID            uuid.UUID
	Name            string
	Description     string
	Status          string
	CreatedOn       time.Time
	UpdatedOn       time.Time
	CreatedBy       uuid.UUID
	RecipientType   string
	RecipientUUID   uuid.UUID
	ResponsibleUUID uuid.UUID
	Priority        int
}
