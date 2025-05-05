package auth

import (
	"context"
)

type authInfoContextKey struct{}

var authInfoKey authInfoContextKey

func AuthFromContext[T any](ctx context.Context) (*T, bool) {
	t, ok := ctx.Value(authInfoKey).(*T)
	return t, ok
}

func SetContext[T any](ctx context.Context, info *T) context.Context {
	return context.WithValue(ctx, authInfoKey, info)
}
