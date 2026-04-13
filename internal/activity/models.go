package activity

import (
	"encoding/json"
	"time"

	"github.com/uptrace/bun"
	"github.com/zombie-check-bot/bot/internal/db"
)

type stateModel struct {
	bun.BaseModel `bun:"table:activity_states,alias:as"`
	db.TimedModel
	UserID            string `bun:"user_id,pk"`
	LastAlive         time.Time
	CheckIntervalDays int
	TimeoutDays       int
	RemindersJSON     string `bun:"reminders_json"`
	IsNotified        bool   `bun:"is_notified"`
	NotifiedAt        *time.Time
}

func newStateModel(userID string, now time.Time, cfg Config) *stateModel {
	return &stateModel{UserID: userID, LastAlive: now, CheckIntervalDays: cfg.DefaultCheckIntervalDays, TimeoutDays: cfg.DefaultTimeoutDays, RemindersJSON: mustMarshalReminders(DefaultReminders())}
}

func (m *stateModel) toDomain() (*State, error) {
	var reminders []Reminder
	if err := json.Unmarshal([]byte(m.RemindersJSON), &reminders); err != nil {
		return nil, err
	}
	return &State{UserID: m.UserID, Settings: Settings{CheckIntervalDays: m.CheckIntervalDays, TimeoutDays: m.TimeoutDays, Reminders: reminders}, LastAlive: m.LastAlive, IsNotified: m.IsNotified, NotifiedAt: m.NotifiedAt, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt}, nil
}

func mustMarshalReminders(v []Reminder) string {
	payload, _ := json.Marshal(v)
	return string(payload)
}

func DefaultReminders() []Reminder {
	return []Reminder{{DaysBefore: 3, RepeatCount: 1}, {DaysBefore: 1, RepeatCount: 1}}
}
