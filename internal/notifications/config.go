package notifications

import (
	"fmt"
	"time"
)

type Config struct {
	AliveCheckCooldown   time.Duration
	TrustedAlertCooldown time.Duration
}

func (c Config) Validate() error {
	if c.AliveCheckCooldown <= 0 {
		return fmt.Errorf("%w: alive check cooldown must be greater than 0", ErrValidationFailed)
	}
	if c.TrustedAlertCooldown <= 0 {
		return fmt.Errorf("%w: trusted alert cooldown must be greater than 0", ErrValidationFailed)
	}
	return nil
}

func (c Config) CooldownByType(t NotificationType) time.Duration {
	switch t {
	case NotificationTypeAliveCheck:
		return c.AliveCheckCooldown
	case NotificationTypeTrustedAlert:
		return c.TrustedAlertCooldown
	default:
		return 0
	}
}
