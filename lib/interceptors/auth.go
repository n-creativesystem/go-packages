package interceptors

import (
	"context"

	"connectrpc.com/connect"
	"github.com/n-creativesystem/go-packages/lib/interceptors/auth"
)

type authenticate[T any] struct {
	validate auth.Validator[T]

	errorHandler func(err error) error
}

var (
	_ connect.Interceptor = (*authenticate[any])(nil)
)

func NewAuthenticate[T any](validate auth.Validator[T]) connect.Interceptor {
	return newAuthenticate(validate)
}

func newAuthenticate[T any](validate auth.Validator[T]) *authenticate[T] {
	return &authenticate[T]{
		validate: validate,
	}
}

func (a *authenticate[T]) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		ctx, err := a.authFunc(ctx, req.Header())
		if err != nil {
			if a.errorHandler != nil {
				return nil, a.errorHandler(err)
			}
			return nil, connect.NewError(connect.CodeUnauthenticated, err)
		}
		return next(ctx, req)
	}
}

func (a *authenticate[T]) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return next
}

func (a *authenticate[T]) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return func(ctx context.Context, conn connect.StreamingHandlerConn) error {
		ctx, err := a.authFunc(ctx, conn.RequestHeader())
		if err != nil {
			if a.errorHandler != nil {
				return a.errorHandler(err)
			}
			return connect.NewError(connect.CodeUnauthenticated, err)
		}
		return next(ctx, conn)
	}
}

func (a *authenticate[T]) authFunc(ctx context.Context, getter auth.Getter) (context.Context, error) {
	tokenInfo, err := a.validate.Execute(ctx, getter)
	if err != nil {
		return nil, auth.ErrUnAuthorization
	}
	return auth.SetContext(ctx, tokenInfo), nil
}
