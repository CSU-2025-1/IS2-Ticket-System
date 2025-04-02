package consumer

import (
	"auth-service/internal/domain/constant"
	"auth-service/internal/domain/entity"
	"context"
	"log/slog"
)

type Registerer interface {
	RegisterUser(ctx context.Context, user *entity.User) error
}

type UserStream interface {
	GetNewUserStream(ctx context.Context) <-chan *entity.User
}

type Consumer struct {
	registerer Registerer
	stream     UserStream
	gcount     int
}

func NewConsumer(registerer Registerer, stream UserStream) *Consumer {
	return &Consumer{
		registerer: registerer,
		stream:     stream,
		gcount:     3,
	}
}

func (c *Consumer) RunUserConsuming(ctx context.Context) {
	users := c.stream.GetNewUserStream(ctx)

	for i := 0; i < c.gcount; i++ {
		go c.handler(ctx, users)
	}
}

func (c *Consumer) handler(ctx context.Context, users <-chan *entity.User) {
	for {
		select {
		case user := <-users:
			err := c.registerer.RegisterUser(ctx, user)
			if err != nil {
				slog.Error("failed to register user",
					slog.String(constant.ErrorField, err.Error()),
				)
			}
		case <-ctx.Done():
			return
		}
	}
}
