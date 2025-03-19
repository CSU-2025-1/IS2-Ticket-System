package grpc

import (
	"context"
	auth_v1 "github.com/CSU-2025-1/IS2-Ticket-System-Proto/auth/gen/go/auth.v1"
	oryclient "github.com/ory/hydra-client-go"
	"google.golang.org/grpc"
	"net"
)

type TokenProvider interface {
	IntrospectOAuth2Token(ctx context.Context, token string) (*oryclient.OAuth2TokenIntrospection, error)
}

type Server struct {
	auth_v1.UnimplementedAuthServer
	token TokenProvider
}

func NewServer(token TokenProvider) *Server {
	return &Server{
		token: token,
	}
}

func (s *Server) Run(listener net.Listener) error {
	grpcServer := grpc.NewServer()
	auth_v1.RegisterAuthServer(grpcServer, s)
	return grpcServer.Serve(listener)
}

func (s *Server) IntrospectToken(ctx context.Context, in *auth_v1.IntrospectTokenRequest) (*auth_v1.IntrospectTokenResponse, error) {
	tokenInfo, err := s.token.IntrospectOAuth2Token(ctx, in.Token)
	if err != nil {
		return nil, err
	}

	return &auth_v1.IntrospectTokenResponse{
		UserUuid: *tokenInfo.Sub,
	}, nil
}
