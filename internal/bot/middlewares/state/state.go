package state

import (
	"context"

	"github.com/zombie-check-bot/bot/internal/state"
	"go.uber.org/zap"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

type contextKey string

const (
	stateKey contextKey = "state"
)

func New(stateSvc *state.Service, logger *zap.Logger) th.Handler {
	return func(ctx *th.Context, update telego.Update) error {
		userID := userID(update)
		if userID != 0 {
			state, err := stateSvc.Get(ctx.Context(), userID)
			if err != nil {
				logger.Error("failed to get state", zap.Int64("user_id", userID), zap.Error(err))
			}
			ctx = ctx.WithValue(stateKey, &state)
		}

		if err := ctx.Next(update); err != nil {
			return err //nolint:wrapcheck // propagate error
		}

		if userID != 0 {
			state, ok := State(ctx)
			if !ok {
				return nil
			}

			if err := stateSvc.Set(ctx.Context(), userID, *state); err != nil {
				logger.Error("failed to set state", zap.Int64("user_id", userID), zap.Error(err))
			}
		}

		return nil
	}
}

func State(ctx context.Context) (*state.State, bool) {
	state, ok := ctx.Value(stateKey).(*state.State)

	return state, ok
}

func userID(update telego.Update) int64 {
	if update.Message != nil && update.Message.From != nil {
		return update.Message.From.ID
	}

	if update.CallbackQuery != nil {
		return update.CallbackQuery.From.ID
	}

	return 0
}
