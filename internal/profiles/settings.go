package profiles

import (
	"fmt"

	"github.com/zombie-check-bot/bot/internal/config"
	"github.com/zombie-check-bot/bot/internal/settings"
)

const settingsKey = "profiles"

func Metadata() settings.Metadata {
	return settings.Metadata{
		Key: settingsKey,
		LoadFrom: func(cfg config.Config) any {
			return Config{DefaultLocale: cfg.Profiles.DefaultLocale}
		},
	}
}

func LoadConfig(s *settings.Service) (Config, error) {
	value, err := s.Get(settingsKey)
	if err != nil {
		return Config{}, err
	}
	cfg, ok := value.(Config)
	if !ok {
		return Config{}, fmt.Errorf("invalid profiles config type")
	}
	return cfg, nil
}
