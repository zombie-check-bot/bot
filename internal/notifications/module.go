package notifications

import (
	"github.com/go-core-fx/logger"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"notifications",
		logger.WithNamedLogger("notifications"),
		fx.Provide(NewRepository, fx.Private),
		fx.Provide(
			fx.Annotate(
				New,
				fx.ParamTags("", "", `group:"notifiers"`),
			),
		),
	)
}

func AsNotifier(t any) any {
	return fx.Annotate(t, fx.ResultTags(`group:"notifiers"`))
}
