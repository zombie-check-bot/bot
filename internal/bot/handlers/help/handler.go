package help

import (
	"fmt"

	"github.com/go-core-fx/telegofx"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/zombie-check-bot/bot/internal/bot/handler"
)

type Handler struct{}

func New() handler.Handler {
	return &Handler{}
}

func (h *Handler) Register(router *telegofx.Router) {
	router.Handle(
		h.handleHelp,
		th.CommandEqual("help"),
		th.AnyMessageWithFrom(),
	)
}

func (h *Handler) handleHelp(ctx *th.Context, update telego.Update) error {
	if update.Message == nil {
		return nil
	}

	text := "Доступные команды:\n" +
		"/start — регистрация пользователя\n" +
		"/profile — показать профиль\n" +
		"/active — отметить свою активность\n" +
		"/contacts — управление доверенными контактами\n" +
		"/help — эта справка\n\n" +
		"Контакты:\n" +
		"`/contacts list`\n" +
		"`/contacts add`\n" +
		"`/contacts remove <id>`\n" +
		"`/contacts activate <id>`\n" +
		"`/contacts deactivate <id>`"

	_, err := ctx.Bot().SendMessage(ctx, tu.Message(
		tu.ID(update.Message.Chat.ID),
		text,
	).WithParseMode(telego.ModeMarkdownV2))
	if err != nil {
		return fmt.Errorf("send help message: %w", err)
	}
	return nil
}
