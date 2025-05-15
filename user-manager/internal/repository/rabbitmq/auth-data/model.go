package auth_data

import "github.com/google/uuid"

type AuthData struct {
	UUID     uuid.UUID `json:"uuid"`
	Login    string    `json:"login"`
	Password string    `json:"password"`
}
