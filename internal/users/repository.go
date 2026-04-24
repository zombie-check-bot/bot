package users

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/uptrace/bun"
	"github.com/zombie-check-bot/bot/internal/db"
)

type Repository struct {
	db *bun.DB
}

func NewRepository(db *bun.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) RegisterOrLogin(ctx context.Context, ident Identity) (*User, error) {
	// First, try to find an existing existingUser with this identity
	existingUser, err := r.Login(ctx, ident)
	if err == nil {
		return existingUser, nil
	}

	if !errors.Is(err, ErrNotFound) {
		return nil, err
	}

	// Create a new user
	newUser := newUserModel(uuid.NewString())

	// Create the identity
	newIdentity := newIdentity(
		newUser.ID,
		string(ident.Provider),
		ident.ProviderID,
		ident.ProviderData,
	)

	// Begin transaction
	err = r.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if _, err = tx.NewInsert().
			Model(newUser).
			Exec(ctx); err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}

		if _, err = tx.NewInsert().
			Model(newIdentity).
			Exec(ctx); err != nil {
			// Check if the error is a duplicate key error (unique constraint violation)
			if db.IsDuplicateKeyError(err) {
				// Another concurrent request created this identity, signal to fetch existing user
				return ErrIdentityExists
			}
			return fmt.Errorf("failed to create identity: %w", err)
		}

		return nil
	})

	// If the identity already exists (race condition), fetch and return the existing user
	if errors.Is(err, ErrIdentityExists) {
		return r.Login(ctx, ident)
	}

	if err != nil {
		return nil, fmt.Errorf("transaction failed: %w", err)
	}

	return newUser.toDomain(), nil
}

func (r *Repository) Login(ctx context.Context, ident Identity) (*User, error) {
	// Find user by identity
	var existingUser userModel
	err := r.db.NewSelect().
		Model(&existingUser).
		Where("id IN (?)", r.db.NewSelect().
			Model((*identity)(nil)).
			Column("user_id").
			Where("provider = ?", ident.Provider).
			Where("provider_id = ?", ident.ProviderID),
		).
		Scan(ctx, &existingUser)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return existingUser.toDomain(), nil
}

func (r *Repository) GetUser(ctx context.Context, userID string) (*User, error) {
	var existingUser userModel
	err := r.db.NewSelect().
		Model(&existingUser).
		Where("id = ?", userID).
		Scan(ctx, &existingUser)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return existingUser.toDomain(), nil
}

func (r *Repository) GetIdentity(ctx context.Context, userID string, provider Provider) (*Identity, error) {
	var existingIdentity identity
	err := r.db.NewSelect().
		Model(&existingIdentity).
		Where("user_id = ?", userID).
		Where("provider = ?", provider).
		Scan(ctx, &existingIdentity)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("failed to find identity: %w", err)
	}

	return existingIdentity.toDomain(), nil
}

func (r *Repository) ListActive(ctx context.Context, skip ...string) ([]User, error) {
	users := make([]userModel, 0)
	q := r.db.NewSelect().
		Model(&users).
		Where("status = ?", StatusActive)
	if len(skip) > 0 {
		q = q.Where("id NOT IN (?)", bun.List(skip))
	}
	if err := q.Scan(ctx); err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	return lo.Map(
		users,
		func(u userModel, _ int) User {
			return *u.toDomain()
		},
	), nil
}
