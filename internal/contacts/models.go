package contacts

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/schema"
	"github.com/zombie-check-bot/bot/internal/db"
)

type contactModel struct {
	bun.BaseModel `bun:"table:contacts,alias:tc"`
	db.TimedModel

	ID       string      `bun:"id,pk"`
	UserID   string      `bun:"user_id,notnull"`
	Name     string      `bun:"name"`
	Type     ContactType `bun:"type,notnull"`
	Value    string      `bun:"value,notnull"`
	IsActive bool        `bun:"is_active,notnull"`
}

func newContactModel(input ContactInput) *contactModel {
	return &contactModel{
		BaseModel: schema.BaseModel{},
		TimedModel: db.TimedModel{
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
		},

		ID:       uuid.NewString(),
		UserID:   input.UserID,
		Name:     input.Name,
		Type:     input.Type,
		Value:    input.Value,
		IsActive: true,
	}
}

func (m *contactModel) toDomain() *Contact {
	if m == nil {
		return nil
	}

	return &Contact{
		ContactInput: ContactInput{
			UserID: m.UserID,
			Name:   m.Name,
			Type:   m.Type,
			Value:  m.Value,
		},

		ID:       m.ID,
		IsActive: m.IsActive,
	}
}
