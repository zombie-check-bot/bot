package contacts

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

type Service struct {
	repo   *Repository
	config Config

	logger *zap.Logger
}

func New(config Config, repo *Repository, logger *zap.Logger) *Service {
	return &Service{repo: repo, config: config, logger: logger}
}

func (s *Service) Count(ctx context.Context, userID string) (int, error) {
	return s.repo.Count(ctx, userID)
}

func (s *Service) Add(ctx context.Context, input ContactInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	count, err := s.repo.Count(ctx, input.UserID)
	if err != nil {
		return err
	}
	if s.config.MaxTrustedContacts > 0 && count >= s.config.MaxTrustedContacts {
		return fmt.Errorf("%w: %v", ErrLimitExceeded, s.config.MaxTrustedContacts)
	}

	return s.repo.Add(ctx, input)
}

func (s *Service) Remove(ctx context.Context, userID, contactID string) error {
	return s.repo.Delete(ctx, userID, contactID)
}

func (s *Service) List(ctx context.Context, userID string) ([]Contact, error) {
	return s.repo.List(ctx, userID)
}

func (s *Service) Activate(ctx context.Context, userID, contactID string) error {
	return s.repo.SetActive(ctx, userID, contactID, true)
}

func (s *Service) Deactivate(ctx context.Context, userID, contactID string) error {
	return s.repo.SetActive(ctx, userID, contactID, false)
}
