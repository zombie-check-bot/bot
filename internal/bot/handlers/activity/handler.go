package activity

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-core-fx/telegofx"
	"github.com/go-core-fx/telegofx/predicates"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/zombie-check-bot/bot/internal/activity"
	"github.com/zombie-check-bot/bot/internal/bot/handler"
	"github.com/zombie-check-bot/bot/internal/bot/middlewares/userauth"
	"go.uber.org/zap"
)

const (
	callbackPrefix         = "users:activity:"
	callbackSettings       = callbackPrefix + "settings"
	callbackMarkAlive      = callbackPrefix + "alive"
	callbackCycleInterval  = callbackPrefix + "interval"
	callbackCycleTimeout   = callbackPrefix + "timeout"
	callbackCycleReminders = callbackPrefix + "reminders"
)

type Handler struct {
	activitySvc *activity.Service
	logger      *zap.Logger
}

func New(activitySvc *activity.Service, logger *zap.Logger) handler.Handler {
	return &Handler{activitySvc: activitySvc, logger: logger}
}

func (h *Handler) Register(router *telegofx.Router) {
	router.Handle(h.handleSettings, th.CommandEqual("settings"), th.AnyMessageWithFrom(), predicates.MessageWithChatType(telego.ChatTypePrivate))
	router.Handle(h.handleStatus, th.CommandEqual("status"), th.AnyMessageWithFrom(), predicates.MessageWithChatType(telego.ChatTypePrivate))
	router.Handle(h.handleAlive, th.CommandEqual("alive"), th.AnyMessageWithFrom(), predicates.MessageWithChatType(telego.ChatTypePrivate))
	router.Handle(h.handleCallback, th.CallbackDataPrefix(callbackPrefix))
}

func (h *Handler) handleSettings(ctx *th.Context, update telego.Update) error {
	if update.Message == nil {
		return nil
	}
	user, err := userauth.User(ctx)
	if err != nil {
		return err
	}
	state, err := h.activitySvc.Ensure(ctx, user.ID)
	if err != nil {
		return err
	}
	return h.sendSettings(ctx, update.Message.Chat.ID, state)
}

func (h *Handler) handleStatus(ctx *th.Context, update telego.Update) error {
	if update.Message == nil {
		return nil
	}
	user, err := userauth.User(ctx)
	if err != nil {
		return err
	}
	state, err := h.activitySvc.Ensure(ctx, user.ID)
	if err != nil {
		return err
	}
	_, err = ctx.Bot().SendMessage(ctx, tu.Message(tu.ID(update.Message.Chat.ID), h.statusText(state)).WithReplyMarkup(h.inlineKeyboard(state)))
	return err
}

func (h *Handler) handleAlive(ctx *th.Context, update telego.Update) error {
	if update.Message == nil {
		return nil
	}
	return h.markAliveAndReply(ctx, update.Message.Chat.ID)
}

func (h *Handler) handleCallback(ctx *th.Context, update telego.Update) error {
	if update.CallbackQuery == nil || update.CallbackQuery.Message == nil {
		return nil
	}
	chatID := update.CallbackQuery.Message.GetChat().ID
	data := update.CallbackQuery.Data
	if err := h.answerCallback(ctx, update.CallbackQuery.ID, "Готово"); err != nil {
		return err
	}

	user, err := userauth.User(ctx)
	if err != nil {
		return err
	}

	state, err := h.activitySvc.Ensure(ctx, user.ID)
	if err != nil {
		return err
	}

	switch data {
	case callbackSettings:
		return h.sendSettings(ctx, chatID, state)
	case callbackMarkAlive:
		return h.markAliveAndReply(ctx, chatID)
	case callbackCycleInterval:
		state.CheckIntervalDays = cycleInt(state.CheckIntervalDays, []int{1, 3, 7, 14, 30})
	case callbackCycleTimeout:
		state.TimeoutDays = cycleInt(state.TimeoutDays, []int{7, 14, 21, 30, 60})
	case callbackCycleReminders:
		state.Reminders = cycleReminders(state.Reminders)
	default:
		return nil
	}

	updated, err := h.activitySvc.UpdateSettings(ctx, user.ID, state.Settings)
	if err != nil {
		if errors.Is(err, activity.ErrValidationFailed) {
			return h.reply(ctx, chatID, "Некорректные настройки. Таймаут должен быть не меньше интервала")
		}
		return err
	}
	return h.sendSettings(ctx, chatID, updated)
}

func (h *Handler) markAliveAndReply(ctx *th.Context, chatID int64) error {
	user, err := userauth.User(ctx)
	if err != nil {
		return err
	}
	state, err := h.activitySvc.Ensure(ctx, user.ID)
	if err != nil {
		return err
	}
	state, err = h.activitySvc.MarkAlive(ctx, user.ID)
	if err != nil {
		return err
	}
	text := fmt.Sprintf("✅ Отлично, отметил активность. Следующая плановая проверка через %d дн.", state.CheckIntervalDays)
	_, err = ctx.Bot().SendMessage(ctx, tu.Message(tu.ID(chatID), text).WithReplyMarkup(h.inlineKeyboard(state)))
	return err
}

func (h *Handler) sendSettings(ctx *th.Context, chatID int64, state *activity.State) error {
	_, err := ctx.Bot().SendMessage(ctx, tu.Message(tu.ID(chatID), "⚙️ Настройки Zombie Check\n\n"+h.statusText(state)).WithReplyMarkup(h.inlineKeyboard(state)))
	return err
}

func (h *Handler) statusText(state *activity.State) string {
	return strings.Join([]string{
		fmt.Sprintf("Интервал проверок: %d дн.", state.CheckIntervalDays),
		fmt.Sprintf("Таймаут неактивности: %d дн.", state.TimeoutDays),
		fmt.Sprintf("Напоминания: %s", remindersText(state.Reminders)),
		fmt.Sprintf("Последний /alive: %s", state.LastAlive.Format("02.01.2006 15:04 UTC")),
		fmt.Sprintf("Следующий check-in: %s (%s)", state.NextCheckAt().Format("02.01.2006"), leftText(state.NextCheckAt())),
		fmt.Sprintf("Дедлайн инцидента: %s (%s)", state.DeadlineAt().Format("02.01.2006"), leftText(state.DeadlineAt())),
	}, "\n")
}

func (h *Handler) inlineKeyboard(state *activity.State) *telego.InlineKeyboardMarkup {
	buttons := [][]telego.InlineKeyboardButton{
		{tu.InlineKeyboardButton("✅ Я жив!").WithCallbackData(callbackMarkAlive)},
		{tu.InlineKeyboardButton("Интервал: " + strconv.Itoa(state.CheckIntervalDays) + "д").WithCallbackData(callbackCycleInterval)},
		{tu.InlineKeyboardButton("Таймаут: " + strconv.Itoa(state.TimeoutDays) + "д").WithCallbackData(callbackCycleTimeout)},
		{tu.InlineKeyboardButton("Напоминания").WithCallbackData(callbackCycleReminders)},
		{tu.InlineKeyboardButton("🔄 Обновить").WithCallbackData(callbackSettings)},
	}
	return &telego.InlineKeyboardMarkup{InlineKeyboard: buttons}
}

func (h *Handler) answerCallback(ctx *th.Context, callbackID, text string) error {
	return ctx.Bot().AnswerCallbackQuery(ctx, &telego.AnswerCallbackQueryParams{CallbackQueryID: callbackID, Text: text})
}

func (h *Handler) reply(ctx *th.Context, chatID int64, text string) error {
	_, err := ctx.Bot().SendMessage(ctx, tu.Message(tu.ID(chatID), text))
	return err
}

func remindersText(reminders []activity.Reminder) string {
	if len(reminders) == 0 {
		return "выключены"
	}
	parts := make([]string, 0, len(reminders))
	for _, r := range reminders {
		parts = append(parts, fmt.Sprintf("за %dд x%d", r.DaysBefore, r.RepeatCount))
	}
	return strings.Join(parts, ", ")
}

func cycleReminders(rem []activity.Reminder) []activity.Reminder {
	if len(rem) == 0 {
		return activity.DefaultReminders()
	}
	if len(rem) == 2 && rem[0].DaysBefore == 3 && rem[1].DaysBefore == 1 {
		return []activity.Reminder{{DaysBefore: 1, RepeatCount: 1}}
	}
	return []activity.Reminder{}
}

func cycleInt(current int, options []int) int {
	for i, v := range options {
		if v == current {
			return options[(i+1)%len(options)]
		}
	}
	return options[0]
}

func leftText(t time.Time) string {
	d := time.Until(t)
	days := int(d.Hours() / 24)
	if days < 0 {
		return "просрочено"
	}
	if days == 0 {
		return "сегодня"
	}
	return fmt.Sprintf("через %d дн.", days)
}
