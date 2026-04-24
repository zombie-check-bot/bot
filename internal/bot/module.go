package bot

import (
	"github.com/go-core-fx/logger"
	"github.com/go-core-fx/telegofx"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
	"github.com/zombie-check-bot/bot/internal/bot/handler"
	"github.com/zombie-check-bot/bot/internal/bot/handlers/activity"
	"github.com/zombie-check-bot/bot/internal/bot/handlers/cancel"
	"github.com/zombie-check-bot/bot/internal/bot/handlers/contacts"
	"github.com/zombie-check-bot/bot/internal/bot/handlers/help"
	"github.com/zombie-check-bot/bot/internal/bot/handlers/profile"
	"github.com/zombie-check-bot/bot/internal/bot/handlers/start"
	"github.com/zombie-check-bot/bot/internal/bot/middlewares/state"
	"github.com/zombie-check-bot/bot/internal/bot/middlewares/userauth"
	"github.com/zombie-check-bot/bot/internal/notifications"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"bot",
		logger.WithNamedLogger("bot"),
		fx.Provide(func() []telego.BotOption {
			return []telego.BotOption{
				telego.WithFastHTTPClient(&fasthttp.Client{Dial: fasthttpproxy.FasthttpProxyHTTPDialer()}),
			}
		}),
		fx.Provide(
			newNotifier,
			notifications.AsNotifier(
				ProvideNotifier,
			),
		),
		fx.Provide(
			fx.Annotate(userauth.New, fx.ResultTags(`name:"middlewares-userauth"`)),
			fx.Annotate(state.New, fx.ResultTags(`name:"middlewares-state"`)),

			fx.Annotate(start.New, fx.ResultTags(`group:"handlers"`)),
			fx.Annotate(profile.New, fx.ResultTags(`group:"handlers"`)),
			fx.Annotate(contacts.New, fx.ResultTags(`group:"handlers"`)),
			fx.Annotate(activity.New, fx.ResultTags(`group:"handlers"`)),
			fx.Annotate(cancel.New, fx.ResultTags(`group:"handlers"`)),
			fx.Annotate(help.New, fx.ResultTags(`group:"handlers"`)),
		),
		fx.Invoke(
			fx.Annotate(
				func(handlers []handler.Handler, usersauthMw th.Handler, stateMw th.Handler, r *telegofx.Router) {
					r.Use(
						usersauthMw,
						stateMw,
					)

					for _, h := range handlers {
						h.Register(r)
					}
				},
				fx.ParamTags(`group:"handlers"`, `name:"middlewares-userauth"`, `name:"middlewares-state"`),
			),
		),
	)
}
