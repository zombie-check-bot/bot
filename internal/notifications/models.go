package notifications

import (
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/schema"
)

type notificationModel struct {
	bun.BaseModel `bun:"table:notifications,alias:n"`

	ID        int64     `bun:"id,pk,autoincrement"`
	UserID    string    `bun:"user_id,notnull"`
	Type      string    `bun:"type,notnull"`
	Channel   string    `bun:"channel,notnull"`
	Recipient string    `bun:"recipient,notnull"`
	SentAt    time.Time `bun:"sent_at,notnull"`
}

func newNotificationModel(n Notification) *notificationModel {
	return &notificationModel{
		BaseModel: schema.BaseModel{},

		ID:        0,
		UserID:    n.UserID,
		Type:      string(n.Type),
		Channel:   string(n.Channel),
		Recipient: n.Recipient,
		SentAt:    n.SentAt,
	}
}

// func (n *notificationModel) toDomain() *Notification {
// 	if n == nil {
// 		return nil
// 	}

// 	return &Notification{
// 		UserID:    n.UserID,
// 		Type:      NotificationType(n.Type),
// 		Channel:   contacts.ContactType(n.Channel),
// 		Recipient: n.Recipient,
// 		SentAt:    n.SentAt,
// 	}
// }
