package trace

import (
	"context"
	"log/slog"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

func NewTracerProvider(ctx context.Context, opts ...Option) (trace.TracerProvider, context.CancelFunc, error) {
	opt := defaultOption
	for _, o := range opts {
		o.apply(opt)
	}
	var (
		tp      trace.TracerProvider
		cleanup func()
	)
	if !opt.enabled {
		tp = noop.NewTracerProvider()
		cleanup = func() {}
	} else {
		traceClient := otlptracegrpc.NewClient(
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithEndpoint(opt.agentAddr),
		)
		exporter, err := otlptrace.New(ctx, traceClient)
		if err != nil {
			return nil, nil, err
		}
		r := newResource(opt.serviceName, opt.version, opt.environment)

		sdkTP := sdktrace.NewTracerProvider(
			append([]sdktrace.TracerProviderOption{
				sdktrace.WithBatcher(exporter),
				sdktrace.WithResource(r),
			}, opt.sdktraceOptions...)...,
		)
		pp := newPropagator()
		otel.SetTextMapPropagator(pp)
		cleanup = func() {
			f := func(fn func(ctx context.Context) error) {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				if err := fn(ctx); err != nil {
					slog.Error(err.Error())
				}
				cancel()
			}
			f(sdkTP.ForceFlush)
			f(sdkTP.Shutdown)
			f(exporter.Shutdown)
		}
		tp = sdkTP
	}
	otel.SetTracerProvider(tp)
	return tp, cleanup, nil
}

func newResource(serviceName string, version string, environment string) *resource.Resource {
	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(serviceName),
		semconv.ServiceVersionKey.String(version),
		semconv.DeploymentEnvironmentKey.String(environment),
		attribute.String("environment", environment),
	)
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}
