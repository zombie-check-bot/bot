package storage

import (
	"context"
	"sync"

	"github.com/zombie-check-bot/bot/internal/state"
)

type Storage struct {
	mu     sync.RWMutex
	states map[string]*state.State
}

func New() state.Storage {
	return &Storage{
		mu:     sync.RWMutex{},
		states: make(map[string]*state.State),
	}
}

// Set implements state.Storage.
func (s *Storage) Set(_ context.Context, key string, st *state.State) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.states == nil {
		s.states = make(map[string]*state.State)
	}
	s.states[key] = st.Clone()
	return nil
}

// Get implements state.Storage.
func (s *Storage) Get(_ context.Context, key string) (*state.State, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.states == nil {
		return nil, state.ErrKeyNotFound
	}
	st, ok := s.states[key]
	if !ok {
		return nil, state.ErrKeyNotFound
	}
	return st.Clone(), nil
}

// Delete implements state.Storage.
func (s *Storage) Delete(_ context.Context, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.states == nil {
		return state.ErrKeyNotFound
	}
	if _, ok := s.states[key]; !ok {
		return state.ErrKeyNotFound
	}
	delete(s.states, key)
	return nil
}
