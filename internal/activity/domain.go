package activity

import (
	"fmt"
	"time"
)

type Reminder struct {
	DaysBefore  int `json:"days_before"`
	RepeatCount int `json:"repeat_count"`
}

type Settings struct {
	CheckIntervalDays int
	TimeoutDays       int
	Reminders         []Reminder
}

type State struct {
	UserID string
	Settings
	LastAlive  time.Time
	IsNotified bool
	NotifiedAt *time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (s Settings) Validate() error {
	if s.CheckIntervalDays <= 0 {
		return fmt.Errorf("%w: check interval must be greater than zero", ErrValidationFailed)
	}
	if s.TimeoutDays <= 0 {
		return fmt.Errorf("%w: timeout must be greater than zero", ErrValidationFailed)
	}
	if s.TimeoutDays < s.CheckIntervalDays {
		return fmt.Errorf("%w: timeout must be greater or equal to check interval", ErrValidationFailed)
	}
	for _, r := range s.Reminders {
		if r.DaysBefore <= 0 || r.RepeatCount <= 0 || r.DaysBefore >= s.TimeoutDays {
			return fmt.Errorf("%w: invalid reminders", ErrValidationFailed)
		}
	}
	return nil
}

func (s State) DeadlineAt() time.Time  { return s.LastAlive.AddDate(0, 0, s.TimeoutDays) }
func (s State) NextCheckAt() time.Time { return s.LastAlive.AddDate(0, 0, s.CheckIntervalDays) }
