package scheduler

import (
	"context"
	"fmt"
	"time"

	"github.com/zombie-check-bot/bot/internal/activity"
	"github.com/zombie-check-bot/bot/internal/notifications"
	"go.uber.org/zap"
)

type Service struct {
	config Config

	activitySvc      *activity.Service
	notificationsSvc *notifications.Service

	logger *zap.Logger
}

func NewService(
	config Config,
	activitySvc *activity.Service,
	notificationsSvc *notifications.Service,
	logger *zap.Logger,
) *Service {
	return &Service{config: config, activitySvc: activitySvc, notificationsSvc: notificationsSvc, logger: logger}
}

func (s *Service) Run(ctx context.Context) error {
	if s.config.CheckInterval <= 0 {
		return fmt.Errorf("%w: scheduler check interval must be positive: %s", ErrInvalidConfig, s.config.CheckInterval)
	}

	ticker := time.NewTicker(s.config.CheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			s.processScheduledChecks(ctx)
		}
	}
}

func (s *Service) processScheduledChecks(ctx context.Context) {
	if pending, err := s.activitySvc.ListPendingNotification(ctx); err != nil {
		s.logger.Error("list pending notification", zap.Error(err))
	} else {
		if len(pending) > 0 {
			s.logger.Info("found users pending notification", zap.Int("users", len(pending)))
		}

		for _, user := range pending {
			_, notifErr := s.notificationsSvc.SendAliveCheck(ctx, user.ID)
			if notifErr != nil {
				s.logger.Error(
					"send alive check",
					zap.String("user_id", user.ID),
					zap.Error(notifErr),
				)
			}
		}
	}

	if dead, err := s.activitySvc.ListDead(ctx); err != nil {
		s.logger.Error("list dead users", zap.Error(err))
	} else {
		if len(dead) > 0 {
			s.logger.Info("found dead users", zap.Int("users", len(dead)))
		}

		for _, user := range dead {
			_, notifErr := s.notificationsSvc.SendTrustedAlert(ctx, user.ID)
			if notifErr != nil {
				s.logger.Error(
					"send trusted alert",
					zap.String("user_id", user.ID),
					zap.Error(notifErr),
				)
			}
		}
	}
}
