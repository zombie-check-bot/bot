package profiles

import (
	"time"

	"github.com/uptrace/bun"
	"github.com/zombie-check-bot/bot/internal/db"
)

type profileModel struct {
	bun.BaseModel `bun:"table:profiles,alias:p"`
	db.TimedModel

	UserID      string `bun:"user_id,pk"`
	Username    string `bun:"username"`
	DisplayName string `bun:"display_name"`
	Locale      string `bun:"locale"`
}

func newProfileModel(userID string, input Profile) *profileModel {
	return &profileModel{
		BaseModel: bun.BaseModel{},
		TimedModel: db.TimedModel{
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
		},

		UserID:      userID,
		Username:    input.Username,
		DisplayName: input.DisplayName,
		Locale:      input.Locale,
	}
}

func (p *profileModel) toDomain() *Profile {
	if p == nil {
		return nil
	}

	return &Profile{
		Username:    p.Username,
		DisplayName: p.DisplayName,
		Locale:      p.Locale,
	}
}
