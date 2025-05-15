package consumer

import (
	"auth-service/internal/domain/entity"
	"auth-service/pkg/rabbitmq"
	"context"
	"log/slog"
)

type Registerer interface {
	RegisterUser(ctx context.Context, user *entity.User) error
}

type Consumer struct {
	registerer Registerer
	reader     *rabbitmq.Reader[AuthData]
}

func NewConsumer(registerer Registerer, cfg rabbitmq.Config) (*Consumer, error) {
	c := &Consumer{
		registerer: registerer,
	}

	reader, err := rabbitmq.CreateJsonReader[AuthData](cfg, c.handler, rabbitmq.DefaultConsumeOption)
	if err != nil {
		return nil, err
	}

	c.reader = reader

	return c, nil
}

func (c *Consumer) RunUserConsuming(ctx context.Context) {
	err := c.reader.Run(ctx)
	if err != nil {
		panic(err.Error())
	}
}

func (c *Consumer) handler(data AuthData) {
	err := c.registerer.RegisterUser(context.Background(), &entity.User{
		UUID:     data.UUID,
		Login:    data.Login,
		Password: data.Password,
	})
	if err != nil {
		slog.Error(err.Error())
	}
}
