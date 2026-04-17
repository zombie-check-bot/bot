package contacts

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
	"github.com/samber/lo"
	"github.com/zombie-check-bot/bot/internal/bot/handler"
	"github.com/zombie-check-bot/bot/internal/bot/middlewares/state"
	"github.com/zombie-check-bot/bot/internal/bot/middlewares/userauth"
	"github.com/zombie-check-bot/bot/internal/contacts"
	"go.uber.org/zap"
)

const (
	callbackPrefix = "users:contacts:"

	callbackList       = callbackPrefix + "list"
	callbackRemovePref = callbackPrefix + "remove:"
	callbackActivate   = callbackPrefix + "activate:"
	callbackDeactivate = callbackPrefix + "deactivate:"

	statePrefix = "contacts:"

	stateAdd = statePrefix + "add"
)

type Handler struct {
	contactsSvc *contacts.Service

	logger *zap.Logger
}

func New(contactsSvc *contacts.Service, logger *zap.Logger) handler.Handler {
	return &Handler{
		contactsSvc: contactsSvc,

		logger: logger,
	}
}

func (h *Handler) Register(router *telegofx.Router) {
	router.Handle(
		h.handleContacts,
		state.WithStateEmpty(),
		th.CommandEqual("contacts"),
		th.AnyMessageWithFrom(),
		predicates.MessageWithChatType(telego.ChatTypePrivate),
	)
	router.Handle(h.handleContactsCallback, state.WithStateEmpty(), th.CallbackDataPrefix(callbackPrefix))
	router.Handle(
		h.handleUsersShared,
		state.WithStateEqual(stateAdd),
		th.AnyMessageWithFrom(),
		predicates.MessageWithUsersShared(),
	)
}

func (h *Handler) handleContacts(ctx *th.Context, update telego.Update) error {
	if update.Message == nil || update.Message.From == nil {
		return nil
	}

	currentUser, err := userauth.User(ctx)
	if err != nil {
		h.logger.Error("failed to get user from context", zap.Error(err))
		return fmt.Errorf("get user from context: %w", err)
	}

	_, _, parts := tu.ParseCommand(update.Message.Text)
	action := "list"
	if len(parts) > 0 {
		action = parts[0]
	}

	switch action {
	case "list":
		return h.sendContactsList(ctx, update.Message.Chat.ID, currentUser.ID)
	case "add":
		s, ok := state.State(ctx)
		if !ok {
			h.logger.Error("failed to get state from context")
			return h.reply(ctx, update.Message.Chat.ID, "Произошла ошибка. Попробуйте позже")
		}

		s.SetName(stateAdd)

		_, sendErr := ctx.Bot().SendMessage(
			ctx,
			tu.Message(
				tu.ID(update.Message.Chat.ID),
				"Отправьте контакт с помощью кнопки внизу",
			).WithReplyMarkup(
				tu.Keyboard(
					tu.KeyboardRow(
						tu.KeyboardButton("📲 Отправить контакт").WithRequestUsers(&telego.KeyboardButtonRequestUsers{
							RequestID:       0,
							UserIsBot:       lo.ToPtr(false),
							UserIsPremium:   nil,
							MaxQuantity:     1,
							RequestName:     lo.ToPtr(true),
							RequestUsername: lo.ToPtr(true),
							RequestPhoto:    nil,
						}),
					),
				),
			),
		)
		if sendErr != nil {
			return fmt.Errorf("send contacts message: %w", sendErr)
		}

		return nil
	case "remove", "activate", "deactivate":
		const partsCount = 2
		if len(parts) < partsCount {
			return h.reply(ctx, update.Message.Chat.ID, "Использование: /contacts <remove|activate|deactivate> <id>")
		}
		contactID := parts[1]
		return h.applyContactAction(ctx, update.Message.Chat.ID, currentUser.ID, action, contactID)
	default:
		return h.reply(
			ctx,
			update.Message.Chat.ID,
			"Неизвестная команда. Используйте /contacts list|add|remove|activate|deactivate",
		)
	}
}

func (h *Handler) handleContactsCallback(ctx *th.Context, update telego.Update) error {
	if update.CallbackQuery == nil || update.CallbackQuery.Message == nil {
		return nil
	}

	data := update.CallbackQuery.Data
	chatID := update.CallbackQuery.Message.GetChat().ID
	currentUser, err := userauth.User(ctx)
	if err != nil {
		h.logger.Error("failed to get user from context", zap.Error(err))
		return fmt.Errorf("get user from context: %w", err)
	}

	if data == callbackList {
		err = h.answerCallback(ctx, update.CallbackQuery.ID, "Обновляю список")
		if err != nil {
			return err
		}
		return h.sendContactsList(ctx, chatID, currentUser.ID)
	}

	action, contactID, ok := parseCallbackAction(data)
	if !ok {
		return nil
	}

	err = h.answerCallback(ctx, update.CallbackQuery.ID, "Выполняю действие")
	if err != nil {
		return err
	}
	return h.applyContactAction(ctx, chatID, currentUser.ID, action, contactID)
}

func (h *Handler) handleUsersShared(ctx *th.Context, update telego.Update) error {
	if update.Message.UsersShared == nil {
		h.logger.Warn("received update without users shared")
		return nil
	}

	user, err := userauth.User(ctx)
	if err != nil {
		h.logger.Error("failed to get user from context", zap.Error(err))
		return fmt.Errorf("get user from context: %w", err)
	}

	var loopErr error
	for _, contact := range update.Message.UsersShared.Users {
		addErr := h.contactsSvc.Add(ctx, contacts.ContactInput{
			UserID: user.ID,
			Name:   makeUserName(contact),
			Type:   contacts.ContactTypeTelegram,
			Value:  strconv.FormatInt(contact.UserID, 10),
		})
		if addErr != nil {
			h.logger.Error(
				"failed to upsert contact",
				zap.String("user_id", user.ID),
				zap.Any("contact", contact),
				zap.Error(addErr),
			)
			loopErr = errors.Join(loopErr, addErr)
		}
	}

	s, ok := state.State(ctx)
	if !ok {
		h.logger.Error("failed to get state from context")
	}
	s.Clear()

	if loopErr != nil {
		return h.reply(ctx, update.Message.Chat.ID, h.userFacingError(loopErr))
	}

	return h.reply(ctx, update.Message.Chat.ID, "Контакты успешно добавлены")
}

func parseCallbackAction(data string) (string, string, bool) {
	for action, prefix := range map[string]string{
		"remove":     callbackRemovePref,
		"activate":   callbackActivate,
		"deactivate": callbackDeactivate,
	} {
		if after, ok := strings.CutPrefix(data, prefix); ok {
			value := after
			return action, value, true
		}
	}
	return "", "", false
}

func (h *Handler) applyContactAction(
	ctx *th.Context,
	chatID int64,
	userID string,
	action string,
	contactID string,
) error {
	var err error
	switch action {
	case "remove":
		err = h.contactsSvc.Remove(ctx, userID, contactID)
	case "activate":
		err = h.contactsSvc.Activate(ctx, userID, contactID)
	case "deactivate":
		err = h.contactsSvc.Deactivate(ctx, userID, contactID)
	default:
		return h.reply(ctx, chatID, "Неизвестное действие")
	}
	if err != nil {
		return h.reply(ctx, chatID, h.userFacingError(err))
	}
	return h.sendContactsList(ctx, chatID, userID)
}

func (h *Handler) sendContactsList(ctx *th.Context, chatID int64, userID string) error {
	contacts, err := h.contactsSvc.List(ctx, userID)
	if err != nil {
		return h.reply(ctx, chatID, "Не удалось загрузить список контактов. Попробуйте позже")
	}

	lines := []string{"Доверенные контакты:"}
	if len(contacts) == 0 {
		lines = append(lines, "— список пуст")
	} else {
		for _, c := range contacts {
			status := "активен"
			if !c.IsActive {
				status = "неактивен"
			}
			lines = append(lines, fmt.Sprintf("%s [%s]", c.String(), status))
		}
	}
	lines = append(lines, "", "Команды: /contacts add @username | /contacts remove <id>")

	kb := h.contactsKeyboard(contacts)
	_, err = ctx.Bot().SendMessage(
		ctx,
		tu.Message(
			tu.ID(chatID),
			strings.Join(lines, "\n"),
		).WithReplyMarkup(kb),
	)
	if err != nil {
		return fmt.Errorf("send contacts list: %w", err)
	}
	return nil
}

func (h *Handler) contactsKeyboard(contacts []contacts.Contact) *telego.InlineKeyboardMarkup {
	buttons := [][]telego.InlineKeyboardButton{{
		tu.InlineKeyboardButton("🔄 Обновить").WithCallbackData(callbackList),
	}}

	for _, c := range contacts {
		row := []telego.InlineKeyboardButton{
			tu.InlineKeyboardButton(fmt.Sprintf("❌ %s", c.String())).
				WithCallbackData(fmt.Sprintf("%s%s", callbackRemovePref, c.ID)),
		}
		if c.IsActive {
			row = append(
				row,
				tu.InlineKeyboardButton("⏸").WithCallbackData(fmt.Sprintf("%s%s", callbackDeactivate, c.ID)),
			)
		} else {
			row = append(
				row,
				tu.InlineKeyboardButton("▶️").WithCallbackData(fmt.Sprintf("%s%s", callbackActivate, c.ID)),
			)
		}
		buttons = append(buttons, row)
	}

	return &telego.InlineKeyboardMarkup{InlineKeyboard: buttons}
}

func (h *Handler) answerCallback(ctx *th.Context, callbackID, text string) error {
	err := ctx.Bot().
		AnswerCallbackQuery(ctx, &telego.AnswerCallbackQueryParams{CallbackQueryID: callbackID, Text: text})
	if err != nil {
		return fmt.Errorf("answer callback: %w", err)
	}
	return nil
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
		return fmt.Errorf("send contacts message: %w", err)
	}
	return nil
}

func (h *Handler) userFacingError(err error) string {
	switch {
	case errors.Is(err, contacts.ErrAlreadyExists):
		return "Такой контакт уже добавлен"
	case errors.Is(err, contacts.ErrLimitExceeded):
		return "Достигнут лимит доверенных контактов"
	case errors.Is(err, contacts.ErrNotFound):
		return "Контакт не найден"
	case errors.Is(err, contacts.ErrValidationFailed):
		return "Некорректные данные. Проверьте формат и попробуйте снова"
	default:
		return "Не удалось выполнить операцию. Попробуйте позже"
	}
}

func makeUserName(user telego.SharedUser) string {
	builder := strings.Builder{}
	if user.FirstName != "" {
		builder.WriteString(user.FirstName)
	}
	if user.LastName != "" {
		if builder.Len() > 0 {
			builder.WriteString(" ")
		}
		builder.WriteString(user.LastName)
	}

	if user.Username != "" {
		hasName := builder.Len() > 0
		if hasName {
			builder.WriteString(" (")
		}
		builder.WriteString("@")
		builder.WriteString(user.Username)
		if hasName {
			builder.WriteString(")")
		}
	}

	return builder.String()
}
