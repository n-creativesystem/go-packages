package trace

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/log/noop"
	sdklog "go.opentelemetry.io/otel/sdk/log"
)

// NewLoggerProvider は OTLP(gRPC) でログをエクスポートする LoggerProvider を生成し、グローバルに設定する。
// WithEnabled(false) の場合は noop.LoggerProvider を返す。
func NewLoggerProvider(ctx context.Context, opts ...Option) (log.LoggerProvider, context.CancelFunc, error) {
	opt := newOption(opts...)
	var (
		lp      log.LoggerProvider
		cleanup func()
	)
	if !opt.enabled {
		lp = noop.NewLoggerProvider()
		cleanup = func() {}
	} else {
		exporter, err := otlploggrpc.New(ctx,
			otlploggrpc.WithInsecure(),
			otlploggrpc.WithEndpoint(opt.agentAddr),
		)
		if err != nil {
			return nil, nil, err
		}
		r := newResource(opt.serviceName, opt.version, opt.environment)

		sdkLP := sdklog.NewLoggerProvider(
			append([]sdklog.LoggerProviderOption{
				sdklog.WithProcessor(sdklog.NewBatchProcessor(exporter)),
				sdklog.WithResource(r),
			}, opt.sdklogOptions...)...,
		)
		// LoggerProvider の Shutdown が exporter の Shutdown も行う
		cleanup = newCleanup(sdkLP.ForceFlush, sdkLP.Shutdown)
		lp = sdkLP
	}
	global.SetLoggerProvider(lp)
	return lp, cleanup, nil
}

// NewSlogHandler はグローバルな LoggerProvider へログを送る slog.Handler を返す。
// ctx に span が含まれている場合、trace_id / span_id がログに紐付けられる。
func NewSlogHandler(opts ...otelslog.Option) slog.Handler {
	opts = append([]otelslog.Option{otelslog.WithLoggerProvider(global.GetLoggerProvider())}, opts...)
	return otelslog.NewHandler(OpenTelemetryTracerName, opts...)
}
