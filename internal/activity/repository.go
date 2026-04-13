package activity

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/samber/lo"
	"github.com/uptrace/bun"
)

type Repository struct {
	db *bun.DB
}

func NewRepository(db *bun.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Add(ctx context.Context, activity Activity) error {
	_, err := r.db.NewInsert().Ignore().Model(newActivityModel(activity)).Exec(ctx)
	if err != nil {
		return fmt.Errorf("insert activity: %w", err)
	}
	return nil
}

func (r *Repository) GetLastByUser(ctx context.Context, userID string) (*Activity, error) {
	var activity activityModel
	err := r.db.NewSelect().Model(&activity).Where("user_id = ?", userID).Order("created_at DESC").Limit(1).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("select last activity: %w", err)
	}
	return activity.toDomain(), nil
}

func (r *Repository) ListActiveSince(ctx context.Context, since time.Time) ([]Activity, error) {
	activities := make([]activityModel, 0)
	if err := r.db.NewSelect().Model(&activities).
		ColumnExpr("user_id, MAX(created_at) as created_at").
		Where("created_at >= ?", since).
		Group("user_id").
		Scan(ctx); err != nil {
		return nil, fmt.Errorf("list activities: %w", err)
	}
	return lo.Map(
		activities,
		func(a activityModel, _ int) Activity {
			return *a.toDomain()
		},
	), nil
}
