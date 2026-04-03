package contacts

import "errors"

var (
	ErrValidationFailed = errors.New("validation failed")
	ErrNotFound         = errors.New("contact not found")
	ErrAlreadyExists    = errors.New("contact already exists")
	ErrLimitExceeded    = errors.New("contact limit exceeded")
)
