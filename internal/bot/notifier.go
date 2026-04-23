package bot

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-core-fx/telegofx"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/zombie-check-bot/bot/internal/bot/handlers/activity"
	"github.com/zombie-check-bot/bot/internal/contacts"
	"github.com/zombie-check-bot/bot/internal/notifications"
)

type notifier struct {
	bot *telegofx.Bot
}

func newNotifier(bot *telegofx.Bot) *notifier {
	return &notifier{bot: bot}
}

// Notify implements [notifications.Notifier].
func (n *notifier) Notify(
	ctx context.Context,
	typ notifications.NotificationType,
	address string,
	message string,
) error {
	switch typ {
	case notifications.NotificationTypeAliveCheck:
		return n.sendAliveCheck(ctx, address, message)
	case notifications.NotificationTypeTrustedAlert:
		return n.sendTrustedAlert(ctx, address, message)
	default:
		return fmt.Errorf("%w: %s", ErrNotificationTypeNotSupported, typ)
	}
}

func (n *notifier) sendAliveCheck(ctx context.Context, recipient string, text string) error {
	chatID, err := parseTelegramChatID(recipient)
	if err != nil {
		return err
	}

	msg := tu.Message(tu.ID(chatID), text).
		WithReplyMarkup(&telego.InlineKeyboardMarkup{InlineKeyboard: [][]telego.InlineKeyboardButton{{
			tu.InlineKeyboardButton("✅ I am alive").WithCallbackData(activity.AliveConfirmCallback),
		}}})
	if _, err = n.bot.SendMessage(ctx, msg); err != nil {
		return fmt.Errorf("send telegram notification: %w", err)
	}

	return nil
}

func (n *notifier) sendTrustedAlert(ctx context.Context, recipient string, text string) error {
	chatID, err := parseTelegramChatID(recipient)
	if err != nil {
		return err
	}

	if _, err = n.bot.SendMessage(ctx, tu.Message(tu.ID(chatID), text)); err != nil {
		return fmt.Errorf("send telegram notification: %w", err)
	}

	return nil
}

func parseTelegramChatID(recipient string) (int64, error) {
	chatID, err := strconv.ParseInt(recipient, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parse telegram recipient: %w", err)
	}
	return chatID, nil
}

func ProvideNotifier(notifier *notifier) notifications.RegistrationMetadata {
	return notifications.RegistrationMetadata{
		Channel:  contacts.ContactTypeTelegram,
		Notifier: notifier,
	}
}

var _ notifications.Notifier = (*notifier)(nil)
