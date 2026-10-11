package trace

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

// NewMeterProvider は OTLP(gRPC) でメトリクスをエクスポートする MeterProvider を生成し、グローバルに設定する。
// WithEnabled(false) の場合は noop.MeterProvider を返す。
func NewMeterProvider(ctx context.Context, opts ...Option) (metric.MeterProvider, context.CancelFunc, error) {
	opt := newOption(opts...)
	var (
		mp      metric.MeterProvider
		cleanup func()
	)
	if !opt.enabled {
		mp = noop.NewMeterProvider()
		cleanup = func() {}
	} else {
		exporter, err := otlpmetricgrpc.New(ctx,
			otlpmetricgrpc.WithInsecure(),
			otlpmetricgrpc.WithEndpoint(opt.agentAddr),
		)
		if err != nil {
			return nil, nil, err
		}
		r := newResource(opt.serviceName, opt.version, opt.environment)

		sdkMP := sdkmetric.NewMeterProvider(
			append([]sdkmetric.Option{
				sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter)),
				sdkmetric.WithResource(r),
			}, opt.sdkmetricOptions...)...,
		)
		// MeterProvider の Shutdown が exporter の Shutdown も行う
		cleanup = newCleanup(sdkMP.ForceFlush, sdkMP.Shutdown)
		mp = sdkMP
	}
	otel.SetMeterProvider(mp)
	return mp, cleanup, nil
}

// Meter はグローバルな MeterProvider から OpenTelemetryTracerName の Meter を返す。
func Meter(opts ...metric.MeterOption) metric.Meter {
	return otel.GetMeterProvider().Meter(OpenTelemetryTracerName, opts...)
}
