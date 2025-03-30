package http

import "github.com/google/uuid"

type RegisterMailReceiverDto struct {
	UserUUID uuid.UUID `json:"user_uuid"`
	Mail     string    `json:"mail"`
}

type Error struct {
	Message string `json:"message"`
}
