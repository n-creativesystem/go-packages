package interceptors

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"connectrpc.com/connect"
	"github.com/google/uuid"
)

type loggingIntercept struct{}

func NewLoggingInterceptor() connect.Interceptor {
	return &loggingIntercept{}
}

var (
	_ connect.Interceptor = (*loggingIntercept)(nil)
)

func (l *loggingIntercept) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, request connect.AnyRequest) (connect.AnyResponse, error) {
		requestId := getRequestId(request.Header())
		requestIdWith := slog.String("request-id", requestId)
		start := time.Now()
		requestWith := []any{
			requestIdWith,
			slog.Time("request-time", start),
		}
		slog.With(requestWith...).Info(fmt.Sprintf("request calling: %s", request.Spec().Procedure))
		slog.Debug(fmt.Sprintf("request body: %v", request.Any()))
		response, err := next(ctx, request)
		end := time.Now()
		latency := end.Sub(start)
		responseWith := []any{
			requestIdWith,
			slog.Time("response-time", end),
			slog.String("latency", secToTime(latency)),
		}
		if err != nil {
			with := append(requestWith, responseWith...)
			slog.With(with...).ErrorContext(ctx, fmt.Errorf("error %w", err).Error())
		} else {
			slog.Debug(fmt.Sprintf("response body: %v", response.Any()))
		}
		slog.With(responseWith...).Info(fmt.Sprintf("response calling: %s", request.Spec().Procedure))
		return response, err
	}
}

func (l *loggingIntercept) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return next
}

func (l *loggingIntercept) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return func(ctx context.Context, conn connect.StreamingHandlerConn) error {
		requestId := getRequestId(conn.RequestHeader())
		requestIdWith := slog.String("request-id", requestId)
		start := time.Now()
		requestWith := []any{
			requestIdWith,
			slog.Time("request-time", start),
		}
		slog.With(requestWith...).Info(fmt.Sprintf("request calling: %s", conn.Spec().Procedure))
		err := next(ctx, conn)
		end := time.Now()
		latency := end.Sub(start)
		responseWith := []any{
			requestIdWith,
			slog.Time("response-time", end),
			slog.String("latency", secToTime(latency)),
		}
		if err != nil {
			with := append(requestWith, responseWith...)
			slog.With(with...).ErrorContext(ctx, fmt.Errorf("error %w", err).Error())
		}
		slog.With(responseWith...).Info(fmt.Sprintf("response calling: %s", conn.Spec().Procedure))
		return err
	}
}

func getRequestId(header http.Header) string {
	value := header.Get("x-request-id")
	if value == "" {
		return uuid.Must(uuid.NewV7()).String()
	}
	return value
}

func secToTime(sec time.Duration) string {
	const format = "15:04:05.000"
	tZero := time.Unix(0, 0).UTC()
	return tZero.Add(sec).Format(format)
}
