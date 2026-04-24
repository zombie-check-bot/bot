package users

import (
	"time"

	"github.com/uptrace/bun"
	"github.com/zombie-check-bot/bot/internal/db"
)

type userModel struct {
	bun.BaseModel `bun:"table:users,alias:u"`
	db.TimedModel

	ID     string `bun:"id,pk"`
	Status Status `bun:"status"`

	Identities []identity `bun:"rel:has-many,join:id=user_id"`
}

func newUserModel(id string) *userModel {
	return &userModel{
		BaseModel: bun.BaseModel{},
		TimedModel: db.TimedModel{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},

		ID:     id,
		Status: StatusActive,

		Identities: nil,
	}
}

type identity struct {
	bun.BaseModel `bun:"table:user_identities,alias:ui"`
	db.TimedModel

	ID     int64  `bun:"id,pk,autoincrement"`
	UserID string `bun:"user_id"`

	Provider     string `bun:"provider"`
	ProviderID   string `bun:"provider_id"`
	ProviderData string `bun:"provider_data,nullzero"`
}

func newIdentity(userID string, provider string, providerID string, providerData string) *identity {
	return &identity{
		BaseModel: bun.BaseModel{},
		TimedModel: db.TimedModel{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},

		ID:     0,
		UserID: userID,

		Provider:     provider,
		ProviderID:   providerID,
		ProviderData: providerData,
	}
}

func (u *userModel) toDomain() *User {
	if u == nil {
		return nil
	}

	return &User{
		ID:     u.ID,
		Status: u.Status,

		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func (u *identity) toDomain() *Identity {
	if u == nil {
		return nil
	}

	return &Identity{
		Provider:     Provider(u.Provider),
		ProviderID:   u.ProviderID,
		ProviderData: u.ProviderData,
	}
}
