package settings

import (
	"fmt"

	"github.com/go-core-fx/logger"
	"github.com/zombie-check-bot/bot/internal/config"
	"go.uber.org/fx"
)

const MetadataGroup = "settings_metadata"

type Metadata struct {
	Key      string
	LoadFrom func(cfg config.Config) any
}

type Service struct {
	values map[string]any
}

type serviceParams struct {
	fx.In

	Config   config.Config
	Metadata []Metadata `group:"settings_metadata"`
}

func New(params serviceParams) (*Service, error) {
	values := make(map[string]any, len(params.Metadata))
	for _, meta := range params.Metadata {
		if meta.Key == "" {
			return nil, fmt.Errorf("settings metadata key is required")
		}
		if meta.LoadFrom == nil {
			return nil, fmt.Errorf("settings metadata loader is required for %s", meta.Key)
		}
		values[meta.Key] = meta.LoadFrom(params.Config)
	}
	return &Service{values: values}, nil
}

func (s *Service) Get(key string) (any, error) {
	v, ok := s.values[key]
	if !ok {
		return nil, fmt.Errorf("settings %q not found", key)
	}
	return v, nil
}

func Module() fx.Option {
	return fx.Module("settings", logger.WithNamedLogger("settings"), fx.Provide(New))
}
