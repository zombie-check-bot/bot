package cancel

import (
	"fmt"

	"github.com/go-core-fx/telegofx"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/zombie-check-bot/bot/internal/bot/handler"
	"github.com/zombie-check-bot/bot/internal/bot/middlewares/state"
	"go.uber.org/zap"
)

type Handler struct {
	logger *zap.Logger
}

func New(logger *zap.Logger) handler.Handler {
	return &Handler{logger: logger}
}

func (h *Handler) Register(router *telegofx.Router) {
	router.Handle(
		h.handleCancel,
		th.CommandEqual("cancel"),
	)
}

func (h *Handler) handleCancel(ctx *th.Context, update telego.Update) error {
	state, ok := state.State(ctx)
	if !ok {
		h.logger.Error("failed to get state from context")
	}

	state.Clear()

	return h.reply(ctx, update.Message.Chat.ID, "Операция отменена")
}

func (h *Handler) reply(ctx *th.Context, chatID int64, text string) error {
	_, err := ctx.Bot().SendMessage(
		ctx,
		tu.Message(
			tu.ID(chatID),
			text,
		).WithReplyMarkup(tu.ReplyKeyboardRemove()),
	)
	if err != nil {
		return fmt.Errorf("send cancel message: %w", err)
	}
	return nil
}
