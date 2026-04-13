package activity

import (
	"fmt"

	"github.com/go-core-fx/telegofx"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/zombie-check-bot/bot/internal/activity"
	"github.com/zombie-check-bot/bot/internal/bot/handler"
	"github.com/zombie-check-bot/bot/internal/bot/middlewares/userauth"
	"go.uber.org/zap"
)

type Handler struct {
	activitySvc *activity.Service

	logger *zap.Logger
}

func New(activitySvc *activity.Service, logger *zap.Logger) handler.Handler {
	return &Handler{activitySvc: activitySvc, logger: logger}
}

func (h *Handler) Register(router *telegofx.Router) {
	router.Handle(
		h.handleActive,
		th.CommandEqual("active"),
		th.AnyMessageWithFrom(),
	)
}

func (h *Handler) handleActive(ctx *th.Context, update telego.Update) error {
	if update.Message == nil || update.Message.From == nil {
		h.logger.Warn("received /active update without message sender")
		return nil
	}

	user, err := userauth.User(ctx)
	if err != nil {
		h.logger.Error("failed to get user from context", zap.Error(err))
		return fmt.Errorf("get user from context: %w", err)
	}

	if markErr := h.activitySvc.MarkActive(ctx, user.ID); markErr != nil {
		h.logger.Error("failed to mark user as active", zap.String("user_id", user.ID), zap.Error(markErr))
		return fmt.Errorf("mark active: %w", markErr)
	}

	return h.reply(ctx, update.Message.Chat.ID,
		tu.Entity("✅ Активность отмечена!\n\n"),
		tu.Entity("Ваша активность успешно зафиксирована в логе."),
	)
}

func (h *Handler) reply(ctx *th.Context, chatID int64, entities ...tu.MessageEntityCollection) error {
	text, ents := tu.MessageEntities(entities...)
	_, err := ctx.Bot().SendMessage(ctx, tu.Message(tu.ID(chatID), text).WithEntities(ents...))
	if err != nil {
		return fmt.Errorf("send active message: %w", err)
	}
	return nil
}
