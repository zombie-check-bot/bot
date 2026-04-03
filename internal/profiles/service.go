package profiles

import (
	"context"

	"go.uber.org/zap"
)

type Service struct {
	config Config

	profiles *Repository

	logger *zap.Logger
}

func NewService(config Config, profiles *Repository, logger *zap.Logger) *Service {
	return &Service{config: config, profiles: profiles, logger: logger}
}

func (s *Service) Upsert(ctx context.Context, userID string, profile Profile) (*Profile, error) {
	if profile.Locale == "" {
		profile.Locale = s.config.DefaultLocale
	}

	return s.profiles.Upsert(ctx, userID, profile)
}

func (s *Service) Get(ctx context.Context, userID string) (*Profile, error) {
	return s.profiles.Get(ctx, userID)
}
