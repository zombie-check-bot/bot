package contacts

import "fmt"

type ContactType string

const (
	ContactTypeTelegram ContactType = "telegram"
)

type ContactInput struct {
	UserID string
	Name   string
	Type   ContactType
	Value  string
}

type Contact struct {
	ContactInput

	ID       string
	IsActive bool
}

func (c ContactInput) Validate() error {
	if c.UserID == "" {
		return fmt.Errorf("%w: user id is required", ErrValidationFailed)
	}
	if !IsValidContactType(c.Type) {
		return fmt.Errorf("%w: invalid contact type", ErrValidationFailed)
	}
	if c.Value == "" {
		return fmt.Errorf("%w: contact value is required", ErrValidationFailed)
	}

	return nil
}

func (c Contact) String() string {
	display := c.Name
	if display == "" {
		display = c.Value
	}
	return fmt.Sprintf("%s (%s)", display, c.Type)
}

func IsValidContactType(t ContactType) bool {
	switch t {
	case ContactTypeTelegram:
		return true
	default:
		return false
	}
}
