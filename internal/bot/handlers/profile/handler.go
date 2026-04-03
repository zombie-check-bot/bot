package profile

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-core-fx/telegofx"
	"github.com/go-core-fx/telegofx/predicates"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/zombie-check-bot/bot/internal/bot/handler"
	"github.com/zombie-check-bot/bot/internal/bot/middlewares/userauth"
	"github.com/zombie-check-bot/bot/internal/contacts"
	"github.com/zombie-check-bot/bot/internal/profiles"
	"go.uber.org/zap"
)

type Handler struct {
	profilesSvc *profiles.Service
	contactsSvc *contacts.Service

	logger *zap.Logger
}

func New(profilesSvc *profiles.Service, contactsSvc *contacts.Service, logger *zap.Logger) handler.Handler {
	return &Handler{profilesSvc: profilesSvc, contactsSvc: contactsSvc, logger: logger}
}

func (h *Handler) Register(router *telegofx.Router) {
	router.Handle(
		h.handleProfile,
		th.CommandEqual("profile"),
		th.AnyMessageWithFrom(),
		predicates.MessageWithChatType(telego.ChatTypePrivate),
	)
}

func (h *Handler) handleProfile(ctx *th.Context, update telego.Update) error {
	if update.Message == nil {
		return nil
	}

	currentUser, err := userauth.User(ctx)
	if err != nil {
		h.logger.Error("failed to get user from context", zap.Error(err))
		return fmt.Errorf("get user from context: %w", err)
	}

	profile, err := h.profilesSvc.Get(ctx, currentUser.ID)
	switch {
	case errors.Is(err, profiles.ErrNotFound):
		return h.reply(ctx, update.Message.Chat.ID, tu.Entity("Профиль не зарегистрирован. Используйте /start"))
	case err != nil:
		return fmt.Errorf("get profile: %w", err)
	}

	displayName := strings.TrimSpace(profile.DisplayName)
	if displayName == "" {
		displayName = "—"
	}
	username := strings.TrimSpace(profile.Username)
	if username == "" {
		username = "—"
	}
	locale := strings.TrimSpace(profile.Locale)

	contactsCount, err := h.contactsSvc.Count(ctx, currentUser.ID)
	if err != nil {
		return fmt.Errorf("count contacts: %w", err)
	}

	return h.reply(ctx, update.Message.Chat.ID,
		tu.Entity("Профиль:\n\n").Bold(),
		tu.Entity("ID: ").Bold(), tu.Entity(currentUser.ID).Code(), tu.Entity("\n"),
		tu.Entity("Логин: ").Bold(), tu.Entity(username).Code(), tu.Entity("\n"),
		tu.Entity("Отображаемое имя: ").Bold(), tu.Entity(displayName), tu.Entity("\n"),
		tu.Entity("Язык: ").Bold(), tu.Entity(locale), tu.Entity("\n"),
		tu.Entity("Доверенных контактов: ").Bold(), tu.Entity(strconv.Itoa(contactsCount)),
	)
}

func (h *Handler) reply(ctx *th.Context, chatID int64, entities ...tu.MessageEntityCollection) error {
	text, ents := tu.MessageEntities(entities...)
	_, err := ctx.Bot().SendMessage(ctx, tu.Message(tu.ID(chatID), text).WithEntities(ents...))
	if err != nil {
		return fmt.Errorf("send profile message: %w", err)
	}
	return nil
}
