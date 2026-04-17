package state

import (
	"context"
	"errors"
)

type Storage interface {
	Get(ctx context.Context, key string) (*State, error)
	Set(ctx context.Context, key string, state *State) error
	Delete(ctx context.Context, key string) error
}

var (
	ErrKeyNotFound = errors.New("key not found")
)
