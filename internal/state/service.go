package state

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"go.uber.org/zap"
)

type Service struct {
	storage Storage

	logger *zap.Logger
}

func NewService(storage Storage, logger *zap.Logger) *Service {
	return &Service{
		storage: storage,
		logger:  logger,
	}
}

func (s *Service) Get(ctx context.Context, userID int64) (State, error) {
	item, err := s.storage.Get(ctx, strconv.FormatInt(userID, 10))
	if errors.Is(err, ErrKeyNotFound) {
		return State{Name: "", Data: map[string]string{}}, nil
	}

	if err != nil {
		return State{Name: "", Data: map[string]string{}}, fmt.Errorf("get state: %w", err)
	}

	if item == nil {
		return State{Name: "", Data: map[string]string{}}, nil
	}
	if item.Data == nil {
		item.Data = map[string]string{}
	}

	return *item, nil
}

func (s *Service) Set(ctx context.Context, userID int64, state State) error {
	if err := s.storage.Set(ctx, strconv.FormatInt(userID, 10), &state); err != nil {
		return fmt.Errorf("set state: %w", err)
	}

	return nil
}

func (s *Service) Delete(ctx context.Context, userID int64) error {
	if err := s.storage.Delete(ctx, strconv.FormatInt(userID, 10)); err != nil {
		if errors.Is(err, ErrKeyNotFound) {
			return nil
		}
		return fmt.Errorf("delete state: %w", err)
	}

	return nil
}
