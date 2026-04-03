package userauth

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	"github.com/zombie-check-bot/bot/internal/users"
	"go.uber.org/zap"
)

type contextKey string

const (
	userKey contextKey = "user"
)

var (
	ErrUserIsNil = errors.New("user is nil")
)

func New(usersSvc *users.Service, logger *zap.Logger) th.Handler {
	return func(ctx *th.Context, update telego.Update) error {
		var telegramUser *telego.User
		switch {
		case update.Message != nil:
			telegramUser = update.Message.From
		case update.CallbackQuery != nil:
			telegramUser = &update.CallbackQuery.From
		default:
			return ctx.Next(update)
		}

		identity, err := identityFromTelegramUser(telegramUser)
		if err != nil {
			logger.Error("failed to map telegram user to identity", zap.Error(err))
			return ctx.Next(update)
		}

		user, err := usersSvc.RegisterOrLogin(ctx, *identity)
		if err != nil {
			logger.Error("failed to register or login telegram user", zap.Error(err))
			return ctx.Next(update)
		}

		ctx = ctx.WithValue(userKey, user)

		return ctx.Next(update)
	}
}

func User(ctx *th.Context) (*users.User, error) {
	user, ok := ctx.Value(userKey).(*users.User)
	if !ok || user == nil {
		return nil, ErrUserIsNil
	}
	return user, nil
}

func identityFromTelegramUser(user *telego.User) (*users.Identity, error) {
	if user == nil {
		return nil, ErrUserIsNil
	}

	providerData, _ := json.Marshal(map[string]any{
		"username":      user.Username,
		"first_name":    user.FirstName,
		"last_name":     user.LastName,
		"language_code": user.LanguageCode,
	})

	return &users.Identity{
		Provider:     users.ProviderTelegram,
		ProviderID:   strconv.FormatInt(user.ID, 10),
		ProviderData: string(providerData),
	}, nil
}
