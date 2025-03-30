package entity

import "github.com/google/uuid"

type Group struct {
	UUID uuid.UUID
	Name string
}
