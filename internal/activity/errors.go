package activity

import "errors"

var (
	ErrNotFound         = errors.New("activity state not found")
	ErrValidationFailed = errors.New("validation failed")
)
