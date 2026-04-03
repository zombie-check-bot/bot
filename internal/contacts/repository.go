package contacts

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/samber/lo"
	"github.com/uptrace/bun"
	"github.com/zombie-check-bot/bot/internal/db"
)

type Repository struct {
	db *bun.DB
}

func NewRepository(db *bun.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Add(ctx context.Context, contact ContactInput) error {
	model := newContactModel(contact)
	_, err := r.db.NewInsert().Model(model).Exec(ctx)
	if err != nil {
		if db.IsDuplicateKeyError(err) {
			return ErrAlreadyExists
		}
		return fmt.Errorf("insert trusted contact: %w", err)
	}
	return nil
}

func (r *Repository) GetByID(ctx context.Context, userID, contactID string) (*Contact, error) {
	var contact contactModel
	err := r.db.NewSelect().Model(&contact).Where("id = ? AND user_id = ?", contactID, userID).Limit(1).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("select trusted contact by id: %w", err)
	}
	return contact.toDomain(), nil
}

func (r *Repository) List(ctx context.Context, userID string) ([]Contact, error) {
	contacts := make([]contactModel, 0)
	if err := r.db.NewSelect().Model(&contacts).Where("user_id = ?", userID).Order("id ASC").Scan(ctx); err != nil {
		return nil, fmt.Errorf("list trusted contacts: %w", err)
	}
	return lo.Map(
		contacts,
		func(c contactModel, _ int) Contact {
			return *c.toDomain()
		},
	), nil
}

func (r *Repository) Count(ctx context.Context, userID string) (int, error) {
	count, err := r.db.NewSelect().Model((*contactModel)(nil)).Where("user_id = ?", userID).Count(ctx)
	if err != nil {
		return 0, fmt.Errorf("count trusted contacts: %w", err)
	}
	return count, nil
}

func (r *Repository) Delete(ctx context.Context, userID, contactID string) error {
	res, err := r.db.NewDelete().
		Model((*contactModel)(nil)).
		Where("id = ? AND user_id = ?", contactID, userID).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("delete trusted contact: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete trusted contact rows affected: %w", err)
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) SetActive(ctx context.Context, userID, contactID string, active bool) error {
	res, err := r.db.NewUpdate().Model((*contactModel)(nil)).
		Set("is_active = ?", active).
		Where("id = ? AND user_id = ?", contactID, userID).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("set trusted contact active: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("set trusted contact active rows affected: %w", err)
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}
