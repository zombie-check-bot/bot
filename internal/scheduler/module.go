package scheduler

import (
	"github.com/go-core-fx/fxutil"
	"github.com/go-core-fx/logger"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"scheduler",
		logger.WithNamedLogger("scheduler"),
		fx.Provide(NewService),
		fx.Invoke(fxutil.RegisterRunnable[*Service]()),
	)
}
