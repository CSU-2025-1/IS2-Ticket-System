package auth

import (
	"auth-service/internal/domain/entity"
	"auth-service/internal/domain/errors/repository"
	"auth-service/internal/domain/errors/service"
	"context"
	"errors"
	"fmt"
	oryclient "github.com/ory/hydra-client-go"
	"log/slog"
)

type Hasher interface {
	Hash(in string) string
}

type UserGetter interface {
	GetUserByLogin(ctx context.Context, login string) (*entity.User, error)
}

type Authenticator interface {
	AcceptLoginRequest(ctx context.Context, challenge, sub string) (*oryclient.CompletedRequest, error)
	AcceptConsentRequest(ctx context.Context, challenge string, scopes []string) (*oryclient.CompletedRequest, error)
}

type Service struct {
	authenticator Authenticator
	userGetter    UserGetter
	hasher        Hasher
}

func New(hydra Authenticator, userRepo UserGetter, hasher Hasher) *Service {
	return &Service{
		authenticator: hydra,
		userGetter:    userRepo,
		hasher:        hasher,
	}
}

func (s *Service) Authenticate(ctx context.Context, challenge, login, password string) (string, error) {
	user, err := s.userGetter.GetUserByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			slog.Warn("Login not found",
				slog.String("login", login),
			)
			return "", service.ErrInvalidCredentials
		}

		return "", fmt.Errorf("auth service: get user by login: %w", err)
	}

	if user.Password != s.hasher.Hash(password) {
		return "", service.ErrInvalidCredentials
	}

	oryAcceptLoginResp, err := s.authenticator.AcceptLoginRequest(ctx, challenge, user.UUID.String())
	if err != nil {
		return "", fmt.Errorf("auth service: accept login request: %w", err)
	}

	return oryAcceptLoginResp.GetRedirectTo(), nil
}

func (s *Service) Consent(ctx context.Context, challenge string, scopes []string) (string, error) {
	oryAcceptConsentResp, err := s.authenticator.AcceptConsentRequest(ctx, challenge, scopes)
	if err != nil {
		return "", fmt.Errorf("auth service: accept consent request: %w", err)
	}

	return oryAcceptConsentResp.GetRedirectTo(), nil
}
