package interceptors

import (
	"context"
	"fmt"
	"log/slog"

	"connectrpc.com/connect"
	"github.com/n-creativesystem/go-packages/lib/logging"
)

type recovery struct{}

func NewRecovery() connect.Interceptor {
	return &recovery{}
}

func (i *recovery) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (resp connect.AnyResponse, err error) {
		defer func() {
			if r := recover(); r != nil {
				slog.With(logging.WithStack(fmt.Errorf("%v", r))).ErrorContext(ctx, fmt.Sprintf("%+v\n", r))
				err = connect.NewError(connect.CodeInternal, fmt.Errorf("unexpected error"))
			}
		}()
		return next(ctx, req)
	}
}

func (i *recovery) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return func(ctx context.Context, spec connect.Spec) connect.StreamingClientConn {
		return next(ctx, spec)
	}
}

func (i *recovery) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return func(ctx context.Context, conn connect.StreamingHandlerConn) (err error) {
		defer func() {
			if r := recover(); r != nil {
				slog.With(logging.WithStack(fmt.Errorf("%v", r))).ErrorContext(ctx, fmt.Sprintf("%+v\n", r))
				err = connect.NewError(connect.CodeInternal, fmt.Errorf("unexpected error"))
			}
		}()
		return next(ctx, conn)
	}
}
