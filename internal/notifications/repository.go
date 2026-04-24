package notifications

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/uptrace/bun"
)

type Repository struct {
	db *bun.DB
}

func NewRepository(db *bun.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Add(ctx context.Context, n Notification) error {
	if _, err := r.db.NewInsert().Model(newNotificationModel(n)).Exec(ctx); err != nil {
		return fmt.Errorf("insert notification: %w", err)
	}
	return nil
}

func (r *Repository) LastSentAt(
	ctx context.Context,
	userID string,
	notificationType NotificationType,
	recipient string,
) (time.Time, error) {
	var notification notificationModel
	err := r.db.NewSelect().
		Model(&notification).
		Where("user_id = ?", userID).
		Where("type = ?", notificationType).
		Where("recipient = ?", recipient).
		Order("sent_at DESC").
		Limit(1).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return time.Time{}, nil
		}
		return time.Time{}, fmt.Errorf("select last notification: %w", err)
	}

	return notification.SentAt, nil
}
