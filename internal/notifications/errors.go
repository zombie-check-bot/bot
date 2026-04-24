package notifications

import "errors"

var (
	ErrValidationFailed = errors.New("validation failed")
	ErrUnsupportedType  = errors.New("unsupported notification type")
)
