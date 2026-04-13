package profiles

import (
	"github.com/go-core-fx/logger"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"profiles",
		logger.WithNamedLogger("profiles"),
		fx.Provide(fx.Annotate(Metadata, fx.ResultTags(`group:"settings_metadata"`))),
		fx.Provide(LoadConfig),
		fx.Provide(NewRepository, fx.Private),
		fx.Provide(NewService),
	)
}
