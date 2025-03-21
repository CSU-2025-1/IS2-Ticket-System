package grpc

import (
	"context"
	"fmt"
	users_v1 "github.com/CSU-2025-1/IS2-Ticket-System-Proto/user-manager/gen/go/users.v1"
	"github.com/google/uuid"
)

type GroupUsersGetter interface {
	GetGroupUsers(ctx context.Context, groupUUID uuid.UUID) ([]uuid.UUID, error)
}

type Server struct {
	users_v1.UnimplementedUsersServer

	group GroupUsersGetter
}

func NewServer(group GroupUsersGetter) *Server {
	return &Server{
		group: group,
	}
}

func (s *Server) GetUsersByGroupID(ctx context.Context, in *users_v1.GetUsersByGroupIDRequest) (*users_v1.GetUsersByGroupIDResponse, error) {
	userUUID, err := uuid.Parse(in.GroupId)
	if err != nil {
		return nil, fmt.Errorf("invalid group id: %w", err)
	}

	users, err := s.group.GetGroupUsers(ctx, userUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	resp := &users_v1.GetUsersByGroupIDResponse{
		UserUuid: make([]string, len(users)),
	}

	for i := range users {
		resp.UserUuid[i] = users[i].String()
	}

	return resp, nil
}
