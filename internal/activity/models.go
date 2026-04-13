package activity

import (
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/schema"
)

type activityModel struct {
	bun.BaseModel `bun:"table:activity,alias:a"`

	UserID    string    `bun:"user_id,pk"`
	CreatedAt time.Time `bun:"created_at,nullzero,notnull,pk"`
}

func newActivityModel(activity Activity) *activityModel {
	return &activityModel{
		BaseModel: schema.BaseModel{},

		UserID:    activity.UserID,
		CreatedAt: activity.CreatedAt,
	}
}

func (a *activityModel) toDomain() *Activity {
	return &Activity{
		UserID:    a.UserID,
		CreatedAt: a.CreatedAt,
	}
}
