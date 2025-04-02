package user

import (
	"auth-service/internal/domain/entity"
	"github.com/google/uuid"
)

type AuthData struct {
	UUID     uuid.UUID `db:"uuid"`
	Login    string    `db:"login"`
	Password string    `db:"password"`
}

func (a AuthData) ToEntity() *entity.User {
	res := entity.User(a)
	return &res
}
