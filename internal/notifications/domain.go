package notifications

import (
	"fmt"
	"time"

	"github.com/zombie-check-bot/bot/internal/contacts"
)

type NotificationType string

const (
	NotificationTypeAliveCheck   NotificationType = "alive_check"
	NotificationTypeTrustedAlert NotificationType = "trusted_alert"
)

type Notification struct {
	UserID    string
	Type      NotificationType
	Channel   contacts.ContactType
	Recipient string
	SentAt    time.Time
}

func (n Notification) Validate() error {
	if n.UserID == "" {
		return fmt.Errorf("%w: user id is required", ErrValidationFailed)
	}
	switch n.Type {
	case NotificationTypeAliveCheck, NotificationTypeTrustedAlert:
	default:
		return fmt.Errorf("%w: invalid type", ErrValidationFailed)
	}
	if !contacts.IsValidContactType(n.Channel) {
		return fmt.Errorf("%w: invalid channel", ErrValidationFailed)
	}
	if n.Recipient == "" {
		return fmt.Errorf("%w: recipient is required", ErrValidationFailed)
	}
	if n.SentAt.IsZero() {
		return fmt.Errorf("%w: sent at is required", ErrValidationFailed)
	}

	return nil
}
