package state

import (
	"context"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

func WithStateEmpty() th.Predicate {
	return func(ctx context.Context, _ telego.Update) bool {
		state, ok := State(ctx)
		return !ok || state == nil || state.Name == ""
	}
}

func WithStateEqual(name string) th.Predicate {
	return func(ctx context.Context, _ telego.Update) bool {
		state, ok := State(ctx)
		return ok && state.Name == name
	}
}

func WithStatePrefix(prefix string) th.Predicate {
	return func(ctx context.Context, _ telego.Update) bool {
		state, ok := State(ctx)
		return ok && state.IsPrefixed(prefix)
	}
}
