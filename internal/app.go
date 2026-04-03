package internal

import (
	"context"

	"github.com/go-core-fx/bunfx"
	"github.com/go-core-fx/fiberfx"
	"github.com/go-core-fx/goosefx"
	"github.com/go-core-fx/healthfx"
	"github.com/go-core-fx/logger"
	"github.com/go-core-fx/sqlfx"
	"github.com/go-core-fx/telegofx"
	"github.com/zombie-check-bot/bot/internal/bot"
	"github.com/zombie-check-bot/bot/internal/config"
	"github.com/zombie-check-bot/bot/internal/contacts"
	"github.com/zombie-check-bot/bot/internal/db"
	"github.com/zombie-check-bot/bot/internal/profiles"
	"github.com/zombie-check-bot/bot/internal/server"
	"github.com/zombie-check-bot/bot/internal/users"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Run(version healthfx.Version) {
	fx.New(
		// CORE MODULES
		logger.Module(),
		logger.WithFxDefaultLogger(),
		// badgerfx.Module(),
		bunfx.Module(),
		// cachefx.Module(),
		fiberfx.Module(),
		// gocqlfx.Module(),
		// gocqlxfx.Module(),
		sqlfx.Module(),
		goosefx.Module(),
		// gormfx.Module(),
		healthfx.Module(),
		// openrouterfx.Module(),
		// redisfx.Module(),
		// sqlxfx.Module(),
		telegofx.Module(true),
		// validatorfx.Module(),
		// watermillfx.Module(),
		//
		// APP MODULES
		config.Module(),
		db.Module(),
		server.Module(),
		bot.Module(),
		//
		// BUSINESS MODULES
		fx.Supply(version),
		users.Module(),
		profiles.Module(),
		contacts.Module(),
		//
		fx.Invoke(func(lc fx.Lifecycle, logger *zap.Logger) {
			lc.Append(fx.Hook{
				OnStart: func(_ context.Context) error {
					logger.Info("app started")
					return nil
				},
				OnStop: func(_ context.Context) error {
					logger.Info("app stopped")
					return nil
				},
			})
		}),
	).Run()
}
