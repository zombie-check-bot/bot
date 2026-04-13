package activity

import (
	"fmt"
	"time"
)

type Config struct {
	Pending  time.Duration
	Deadline time.Duration
}

func (c Config) Validate() error {
	if c.Pending <= 0 {
		return fmt.Errorf("%w: pending time must be greater than 0", ErrValidationFailed)
	}
	if c.Deadline <= 0 {
		return fmt.Errorf("%w: deadline time must be greater than 0", ErrValidationFailed)
	}

	if c.Pending > c.Deadline {
		return fmt.Errorf("%w: pending time must be less than timeout", ErrValidationFailed)
	}

	return nil
}
