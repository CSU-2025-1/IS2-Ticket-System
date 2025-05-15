package auth_data

import (
	"context"
	"user-mananger/internal/domain/entity"
	"user-mananger/pkg/rabbitmq"
)

type Repository struct {
	writer *rabbitmq.Writer[AuthData]
}

func NewRepository(writer *rabbitmq.Writer[AuthData]) *Repository {
	return &Repository{
		writer: writer,
	}
}

func (r *Repository) Save(ctx context.Context, authData entity.User) error {
	return r.writer.Send(ctx, AuthData(authData))
}
