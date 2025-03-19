package register

import (
	"auth-service/internal/domain/entity"
	"context"
	"fmt"
)

type UserCreator interface {
	CreateUser(ctx context.Context, user *entity.User) error
}

type Hasher interface {
	Hash(in string) string
}

type Service struct {
	saver  UserCreator
	hasher Hasher
}

func New(saver UserCreator, hasher Hasher) *Service {
	return &Service{
		saver:  saver,
		hasher: hasher,
	}
}

func (s *Service) RegisterUser(ctx context.Context, user *entity.User) error {
	user.Password = s.hasher.Hash(user.Password)
	err := s.saver.CreateUser(ctx, user)
	if err != nil {
		return fmt.Errorf("could not create user: %w", err)
	}

	return nil
}
