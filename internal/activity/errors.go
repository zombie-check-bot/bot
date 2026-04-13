package activity

import "errors"

var (
	ErrValidationFailed = errors.New("validation failed")
	ErrNotFound         = errors.New("activity not found")
)
