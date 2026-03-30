package bot

import (
	"github.com/go-core-fx/logger"
	"github.com/go-core-fx/telegofx"
	"github.com/mymmrac/telego"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
	"github.com/zombie-check-bot/bot/internal/bot/handler"
	"github.com/zombie-check-bot/bot/internal/bot/handlers/start"
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
			fx.Annotate(start.New, fx.ResultTags(`group:"handlers"`)),
		),
		fx.Invoke(
			fx.Annotate(
				func(handlers []handler.Handler, r *telegofx.Router) {
					for _, h := range handlers {
						h.Register(r)
					}
				},
				fx.ParamTags(`group:"handlers"`),
			),
		),
	)
}
