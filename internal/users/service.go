package users

import (
	"context"

	"go.uber.org/zap"
)

type Service struct {
	users *Repository

	logger *zap.Logger
}

func New(
	users *Repository,
	logger *zap.Logger,
) *Service {
	return &Service{
		users: users,

		logger: logger,
	}
}

func (s *Service) RegisterOrLogin(ctx context.Context, ident Identity) (*User, error) {
	if err := ident.Validate(); err != nil {
		return nil, err
	}

	return s.users.RegisterOrLogin(ctx, ident)
}

func (s *Service) Login(ctx context.Context, ident Identity) (*User, error) {
	if err := ident.Validate(); err != nil {
		return nil, err
	}

	return s.users.Login(ctx, ident)
}

func (s *Service) GetUser(ctx context.Context, userID string) (*User, error) {
	return s.users.GetUser(ctx, userID)
}

func (s *Service) GetIdentity(ctx context.Context, userID string, provider Provider) (*Identity, error) {
	return s.users.GetIdentity(ctx, userID, provider)
}

func (s *Service) ListActive(ctx context.Context, skip ...string) ([]User, error) {
	return s.users.ListActive(ctx, skip...)
}
