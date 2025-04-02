package service

import (
	"context"
	"github.com/google/uuid"
	"user-mananger/internal/domain/entity"
)

type UserSaver interface {
	SaveUser(ctx context.Context, login string) (uuid.UUID, error)
}

type RegistrationEventSaver interface {
	SaveRegistrationEvent(ctx context.Context, user entity.User) error
}

func RegisterUser(
	ctx context.Context,
	login, password string,
	saver UserSaver,
	event RegistrationEventSaver,
) (uuid.UUID, error) {

}
