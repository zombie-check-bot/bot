package activity

import (
	"context"
	"fmt"
	"time"

	"github.com/samber/lo"
	"github.com/zombie-check-bot/bot/internal/users"
	"go.uber.org/zap"
)

type Service struct {
	config Config

	activities *Repository

	usersSvc *users.Service

	logger *zap.Logger
}

func New(config Config, activities *Repository, usersSvc *users.Service, logger *zap.Logger) (*Service, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &Service{config: config, activities: activities, usersSvc: usersSvc, logger: logger}, nil
}

func (s *Service) MarkActive(ctx context.Context, userID string) error {
	return s.activities.Add(ctx, Activity{UserID: userID, CreatedAt: time.Now()})
}

func (s *Service) ListPendingNotification(ctx context.Context) ([]users.User, error) {
	return s.listNotActive(ctx, s.config.Pending)
}

func (s *Service) ListDead(ctx context.Context) ([]users.User, error) {
	return s.listNotActive(ctx, s.config.Deadline)
}

func (s *Service) listNotActive(ctx context.Context, duration time.Duration) ([]users.User, error) {
	active, err := s.activities.ListActiveSince(ctx, time.Now().Add(-duration))
	if err != nil {
		return nil, err
	}

	activeIDs := lo.Map(active, func(a Activity, _ int) string {
		return a.UserID
	})

	users, err := s.usersSvc.ListActive(ctx, activeIDs...)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}

	return users, nil
}
