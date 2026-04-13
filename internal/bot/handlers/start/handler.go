package start

import (
	"fmt"
	"strings"

	"github.com/go-core-fx/telegofx"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/zombie-check-bot/bot/internal/bot/handler"
	"github.com/zombie-check-bot/bot/internal/bot/middlewares/userauth"
	"github.com/zombie-check-bot/bot/internal/profiles"
	"go.uber.org/zap"
)

type Handler struct {
	profilesSvc *profiles.Service

	logger *zap.Logger
}

func New(profilesSvc *profiles.Service, logger *zap.Logger) handler.Handler {
	return &Handler{profilesSvc: profilesSvc, logger: logger}
}

func (h *Handler) Register(router *telegofx.Router) {
	router.Handle(
		h.handleStart,
		th.CommandEqual("start"),
		th.AnyMessageWithFrom(),
	)
}

func (h *Handler) handleStart(ctx *th.Context, update telego.Update) error {
	if update.Message == nil || update.Message.From == nil {
		h.logger.Warn("received /start update without message sender")
		return nil
	}

	user, err := userauth.User(ctx)
	if err != nil {
		h.logger.Error("failed to get user from context", zap.Error(err))
		return fmt.Errorf("get user from context: %w", err)
	}

	profile, err := h.profilesSvc.Upsert(ctx, user.ID, profiles.Profile{
		Username:    update.Message.From.Username,
		DisplayName: strings.TrimSpace(update.Message.From.FirstName + " " + update.Message.From.LastName),
		Locale:      update.Message.From.LanguageCode,
	})
	if err != nil {
		h.logger.Error("failed to upsert profile", zap.Error(err))
		return fmt.Errorf("upsert profile: %w", err)
	}

	username := profile.DisplayName
	if username == "" {
		username = "друг"
	}

	return h.reply(ctx, update.Message.Chat.ID,
		tu.Entity("Привет, "),
		tu.Entity(username).Bold(),
		tu.Entity("!\nПрофиль сохранён.\n\nКоманды:\n/profile\n/settings\n/status\n/alive\n/contacts\n/help"),
	)
}

func (h *Handler) reply(ctx *th.Context, chatID int64, entities ...tu.MessageEntityCollection) error {
	text, ents := tu.MessageEntities(entities...)
	_, err := ctx.Bot().SendMessage(ctx, tu.Message(tu.ID(chatID), text).WithEntities(ents...))
	if err != nil {
		return fmt.Errorf("send start message: %w", err)
	}
	return nil
}
