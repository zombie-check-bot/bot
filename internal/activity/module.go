package activity

import (
	"github.com/go-core-fx/logger"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"activity",
		logger.WithNamedLogger("activity"),
		fx.Provide(NewRepository, fx.Private),
		fx.Provide(New),
	)
}
