package model

import "github.com/google/uuid"

type Receiver struct {
	ID   uuid.UUID
	Mail string
}
