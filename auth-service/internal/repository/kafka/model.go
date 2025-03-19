package kafka

import (
	"auth-service/internal/domain/entity"
	"github.com/google/uuid"
)

type User struct {
	UserUUID string `json:"user_uuid"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (u *User) ToEntity() (*entity.User, error) {
	userUUID, err := uuid.Parse(u.UserUUID)
	if err != nil {
		return nil, err
	}
	return &entity.User{
		UUID:     userUUID,
		Login:    u.Login,
		Password: u.Password,
	}, nil
}
