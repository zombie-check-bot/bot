package notifications

import (
	"context"

	"github.com/zombie-check-bot/bot/internal/contacts"
)

type Notifier interface {
	Notify(ctx context.Context, typ NotificationType, address string, message string) error
}

type RegistrationMetadata struct {
	Channel  contacts.ContactType
	Notifier Notifier
}
