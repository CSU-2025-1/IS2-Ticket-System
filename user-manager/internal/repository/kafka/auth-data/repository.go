package auth_data

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"user-mananger/internal/domain/entity"
)

type Repository struct {
	k *kafka.Conn
}

func New(k *kafka.Conn) *Repository {
	return &Repository{
		k: k,
	}
}

func (r *Repository) SaveAuthData(ctx context.Context, user entity.User) error {
	bytes, err := json.Marshal(User{
		UserUUID: user.UUID.String(),
		Login:    user.Login,
		Password: user.Password,
	})
	if err != nil {
		return fmt.Errorf("marshal auth data: %w", err)
	}

	_, err = r.k.Write(bytes)
	if err != nil {
		return fmt.Errorf("write auth data: %w", err)
	}

	return nil
}
