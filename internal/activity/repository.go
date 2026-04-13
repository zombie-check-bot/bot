package activity

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/uptrace/bun"
)

type Repository struct {
	db     *bun.DB
	config Config
}

func NewRepository(db *bun.DB, config Config) *Repository { return &Repository{db: db, config: config} }

func (r *Repository) Ensure(ctx context.Context, userID string, now time.Time) (*State, error) {
	state, err := r.Get(ctx, userID)
	if err == nil {
		return state, nil
	}
	if !errors.Is(err, ErrNotFound) {
		return nil, err
	}
	model := newStateModel(userID, now, r.config)
	if _, err = r.db.NewInsert().Model(model).Exec(ctx); err != nil {
		return nil, fmt.Errorf("insert activity state: %w", err)
	}
	return model.toDomain()
}

func (r *Repository) Get(ctx context.Context, userID string) (*State, error) {
	var model stateModel
	if err := r.db.NewSelect().Model(&model).Where("user_id = ?", userID).Limit(1).Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get activity state: %w", err)
	}
	state, err := model.toDomain()
	if err != nil {
		return nil, fmt.Errorf("decode reminders: %w", err)
	}
	return state, nil
}

func (r *Repository) UpdateSettings(ctx context.Context, userID string, settings Settings) error {
	remindersJSON, err := json.Marshal(settings.Reminders)
	if err != nil {
		return fmt.Errorf("marshal reminders: %w", err)
	}
	_, err = r.db.NewUpdate().Model((*stateModel)(nil)).Set("check_interval_days = ?", settings.CheckIntervalDays).Set("timeout_days = ?", settings.TimeoutDays).Set("reminders_json = ?", string(remindersJSON)).Where("user_id = ?", userID).Exec(ctx)
	if err != nil {
		return fmt.Errorf("update activity settings: %w", err)
	}
	return nil
}

func (r *Repository) MarkAlive(ctx context.Context, userID string, now time.Time) error {
	_, err := r.db.NewUpdate().Model((*stateModel)(nil)).Set("last_alive = ?", now).Set("is_notified = FALSE").Set("notified_at = NULL").Where("user_id = ?", userID).Exec(ctx)
	if err != nil {
		return fmt.Errorf("mark alive: %w", err)
	}
	return nil
}
