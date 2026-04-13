package activity

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
)

type Service struct {
	repo   *Repository
	logger *zap.Logger
}

func New(repo *Repository, logger *zap.Logger) *Service { return &Service{repo: repo, logger: logger} }
func (s *Service) Ensure(ctx context.Context, userID string) (*State, error) {
	return s.repo.Ensure(ctx, userID, time.Now().UTC())
}
func (s *Service) Get(ctx context.Context, userID string) (*State, error) {
	return s.repo.Get(ctx, userID)
}

func (s *Service) UpdateSettings(ctx context.Context, userID string, settings Settings) (*State, error) {
	if err := settings.Validate(); err != nil {
		return nil, err
	}
	if err := s.repo.UpdateSettings(ctx, userID, settings); err != nil {
		return nil, err
	}
	return s.repo.Get(ctx, userID)
}

func (s *Service) MarkAlive(ctx context.Context, userID string) (*State, error) {
	now := time.Now().UTC()
	if err := s.repo.MarkAlive(ctx, userID, now); err != nil {
		return nil, err
	}
	state, err := s.repo.Get(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get state after mark alive: %w", err)
	}
	return state, nil
}
