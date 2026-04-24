package notifications

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/samber/lo"
	"github.com/zombie-check-bot/bot/internal/contacts"
	"github.com/zombie-check-bot/bot/internal/profiles"
	"github.com/zombie-check-bot/bot/internal/users"
	"go.uber.org/zap"
)

type Service struct {
	config Config

	notifications *Repository
	notifiers     map[contacts.ContactType]Notifier

	usersSvc    *users.Service
	contactsSvc *contacts.Service
	profilesSvc *profiles.Service

	logger *zap.Logger
}

func New(
	config Config,
	notifications *Repository,
	notifiers []RegistrationMetadata,
	usersSvc *users.Service,
	contactsSvc *contacts.Service,
	profilesSvc *profiles.Service,
	logger *zap.Logger,
) (*Service, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &Service{
		config: config,

		notifications: notifications,
		notifiers: lo.Associate(
			notifiers,
			func(n RegistrationMetadata) (contacts.ContactType, Notifier) {
				return n.Channel, n.Notifier
			},
		),

		usersSvc:    usersSvc,
		contactsSvc: contactsSvc,
		profilesSvc: profilesSvc,

		logger: logger,
	}, nil
}

func (s *Service) SendAliveCheck(ctx context.Context, userID string) (bool, error) {
	identity, err := s.usersSvc.GetIdentity(ctx, userID, users.ProviderTelegram)
	if err != nil {
		return false, fmt.Errorf("get identity: %w", err)
	}

	n := Notification{
		UserID:    userID,
		Type:      NotificationTypeAliveCheck,
		Channel:   contacts.ContactTypeTelegram,
		Recipient: identity.ProviderID,
		SentAt:    time.Now(),
	}

	return s.send(ctx, n)
}

func (s *Service) SendTrustedAlert(ctx context.Context, userID string) (int, error) {
	contacts, err := s.contactsSvc.List(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("list trusted contacts: %w", err)
	}

	var count int
	var errs error
	for _, c := range contacts {
		if !c.IsActive {
			continue
		}

		n := Notification{
			UserID:    userID,
			Type:      NotificationTypeTrustedAlert,
			Channel:   c.Type,
			Recipient: c.Value,
			SentAt:    time.Now(),
		}

		if ok, sendErr := s.send(ctx, n); sendErr != nil {
			s.logger.Error("failed to send trusted alert",
				zap.String("user_id", userID),
				zap.String("contact_id", c.ID),
				zap.String("channel", string(c.Type)),
				zap.Error(sendErr))
			errs = errors.Join(errs, sendErr)
		} else if ok {
			count++
		}
	}

	return count, errs
}

func (s *Service) send(ctx context.Context, n Notification) (bool, error) {
	if err := n.Validate(); err != nil {
		return false, err
	}

	allowed, err := s.isAllowed(ctx, n)
	if err != nil {
		return false, err
	}
	if !allowed {
		return false, nil
	}

	notifier, ok := s.notifiers[n.Channel]
	if !ok {
		return false, fmt.Errorf("%w: unknown channel", ErrValidationFailed)
	}

	message, err := s.messageFor(ctx, n)
	if err != nil {
		return false, err
	}

	if err = notifier.Notify(ctx, n.Type, n.Recipient, message); err != nil {
		return false, fmt.Errorf("notify: %w", err)
	}

	if err = s.notifications.Add(ctx, n); err != nil {
		return false, err
	}

	return true, nil
}

func (s *Service) isAllowed(ctx context.Context, n Notification) (bool, error) {
	last, err := s.notifications.LastSentAt(ctx, n.UserID, n.Type, n.Recipient)
	if err != nil {
		return false, err
	}
	if last.IsZero() {
		return true, nil
	}

	cooldown := s.config.CooldownByType(n.Type)
	if cooldown <= 0 {
		return false, fmt.Errorf("%w: unknown notification type", ErrValidationFailed)
	}
	return n.SentAt.Sub(last) >= cooldown, nil
}

func (s *Service) messageFor(ctx context.Context, n Notification) (string, error) {
	switch n.Type {
	case NotificationTypeAliveCheck:
		return "⏰ Please confirm you are alive.", nil
	case NotificationTypeTrustedAlert:
		profile, err := s.profilesSvc.Get(ctx, n.UserID)
		if err != nil {
			return "", fmt.Errorf("get profile: %w", err)
		}

		return fmt.Sprintf(
			"🚨 Alert: User %s did not confirm being alive as of %s. Please check on them.",
			profile,
			n.SentAt.Format(time.DateTime),
		), nil
	default:
		return "", ErrUnsupportedType
	}
}
