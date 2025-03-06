package entity

import "github.com/google/uuid"

type UserAuthData struct {
	UUID     uuid.UUID
	Login    string
	Password string
}
