package entity

import "github.com/google/uuid"

type User struct {
	UUID     uuid.UUID
	Login    string
	Password string
}

type UserInfo struct {
	UUID  uuid.UUID
	Login string
}
