package scheduler

import "time"

// Config holds scheduler configuration.
type Config struct {
	CheckInterval time.Duration // How often to poll database (e.g., every 5 minutes)
}
