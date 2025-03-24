package repository

import (
	"context"
	users_v1 "github.com/CSU-2025-1/IS2-Ticket-System-Proto/user-manager/gen/go/users.v1"
	"github.com/google/uuid"
	"notification-service/pkg/grpc"
)

type UserRepository struct {
	pool *grpc.Pool
}

func NewUserRepository(pool *grpc.Pool) *UserRepository {
	return &UserRepository{
		pool: pool,
	}
}

func (u *UserRepository) GetAllUserIDsByGroupID(ctx context.Context, groupID uuid.UUID) ([]uuid.UUID, error) {
	client := users_v1.NewUsersClient(u.pool.Get())
	response, err := client.GetUsersByGroupID(
		ctx,
		&users_v1.GetUsersByGroupIDRequest{
			GroupId: groupID.String(),
		},
	)
	if err != nil {
		return nil, err
	}

	uuids := make([]uuid.UUID, len(response.UserUuid))
	for i := range uuids {
		uuids[i], err = uuid.Parse(response.UserUuid[i])
		if err != nil {
			return nil, err
		}
	}

	return uuids, err
}
