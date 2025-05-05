package auth

import "context"

type Getter interface {
	Get(string) string
}

type Validator[T any] interface {
	Execute(ctx context.Context, getter Getter) (*T, error)
}
