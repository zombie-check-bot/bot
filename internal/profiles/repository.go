package profiles

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/uptrace/bun"
)

type Repository struct {
	db *bun.DB
}

func NewRepository(db *bun.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Upsert(ctx context.Context, userID string, profile Profile) (*Profile, error) {
	model := newProfileModel(userID, profile)

	err := r.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		// Use atomic upsert with ON CONFLICT to avoid check-then-insert race condition
		_, err := tx.NewInsert().
			Model(model).
			On("DUPLICATE KEY UPDATE").
			Set("username = VALUES(username), display_name = VALUES(display_name), locale = VALUES(locale)").
			Returning("*").
			Exec(ctx)
		if err != nil {
			return fmt.Errorf("upsert profile: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("upsert profile: %w", err)
	}
	return model.toDomain(), nil
}

func (r *Repository) Get(ctx context.Context, userID string) (*Profile, error) {
	model, err := r.get(ctx, r.db, userID, false)
	if err != nil {
		return nil, err
	}
	return model.toDomain(), nil
}

func (r *Repository) get(ctx context.Context, tx bun.IDB, userID string, forUpdate bool) (*profileModel, error) {
	var model profileModel
	query := tx.NewSelect().Model(&model).Where("user_id = ?", userID)
	if forUpdate {
		query = query.For("UPDATE")
	}
	err := query.Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}

		return nil, fmt.Errorf("get profile: %w", err)
	}

	return &model, nil
}
