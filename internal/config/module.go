package config

import (
	"github.com/go-core-fx/fiberfx"
	"github.com/go-core-fx/fiberfx/openapi"
	"github.com/go-core-fx/sqlfx"
	"github.com/go-core-fx/telegofx"
	"github.com/zombie-check-bot/bot/internal/activity"
	"github.com/zombie-check-bot/bot/internal/contacts"
	"github.com/zombie-check-bot/bot/internal/notifications"
	"github.com/zombie-check-bot/bot/internal/profiles"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"config",
		fx.Provide(New, fx.Private),
		fx.Provide(
			func(cfg Config) fiberfx.Config {
				return fiberfx.Config{
					Address:     cfg.HTTP.Address,
					ProxyHeader: cfg.HTTP.ProxyHeader,
					Proxies:     cfg.HTTP.Proxies,
				}
			},
			func(cfg Config) openapi.Config {
				return openapi.Config{
					Enabled:    cfg.HTTP.OpenAPI.Enabled,
					PublicHost: cfg.HTTP.OpenAPI.PublicHost,
					PublicPath: cfg.HTTP.OpenAPI.PublicPath,
				}
			},
			func(cfg Config) telegofx.Config {
				return telegofx.Config{Token: cfg.Telegram.Token}
			},
			func(cfg Config) sqlfx.Config {
				return sqlfx.Config{
					URL:             cfg.Database.URL,
					ConnMaxIdleTime: cfg.Database.ConnMaxIdleTime,
					ConnMaxLifetime: cfg.Database.ConnMaxLifetime,
					MaxOpenConns:    cfg.Database.MaxOpenConns,
					MaxIdleConns:    cfg.Database.MaxIdleConns,
				}
			},
		),
		fx.Provide(
			func(cfg Config) profiles.Config {
				return profiles.Config{
					DefaultLocale: cfg.Profiles.DefaultLocale,
				}
			},
			func(cfg Config) contacts.Config {
				return contacts.Config{
					MaxTrustedContacts: cfg.Contacts.MaxTrustedContacts,
				}
			},
			func(cfg Config) activity.Config {
				return activity.Config{
					Pending:  cfg.Activity.PendingTime,
					Deadline: cfg.Activity.DeadlineTime,
				}
			},
			func(cfg Config) notifications.Config {
				return notifications.Config{
					AliveCheckCooldown:   cfg.Notifications.AliveCheckCooldown,
					TrustedAlertCooldown: cfg.Notifications.TrustedAlertCooldown,
				}
			},
		),
	)
}
