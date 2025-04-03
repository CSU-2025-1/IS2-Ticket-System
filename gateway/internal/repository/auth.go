package repository

import (
	"context"
	"fmt"
	"gateway/pkg/grpc"
	auth_v1 "github.com/CSU-2025-1/IS2-Ticket-System-Proto/auth/gen/go/auth.v1"
	"time"
)

type Auth struct {
	pool *grpc.Pool
}

func NewAuth(pool *grpc.Pool) *Auth {
	return &Auth{
		pool: pool,
	}
}

func (a *Auth) Auth(ctx context.Context, token string) (map[string]string, error) {
	client := auth_v1.NewAuthClient(a.pool.Get())

	fmt.Println(time.Now())
	response, err := client.IntrospectToken(ctx, &auth_v1.IntrospectTokenRequest{Token: token})
	if err != nil {
		return nil, err
	}
	fmt.Println(time.Now())
	return map[string]string{
		"X-User-UUID": response.UserUuid,
	}, nil
}
