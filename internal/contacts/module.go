package contacts

import (
	"github.com/go-core-fx/logger"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"contacts",
		logger.WithNamedLogger("contacts"),
		fx.Provide(NewRepository, fx.Private),
		fx.Provide(New),
	)
}
