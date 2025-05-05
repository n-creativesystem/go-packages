package trace

import (
	"context"
	"net/url"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var (
	OpenTelemetryTracerName = "github.com/moneyforward/go-apex-common"
)

func mustJoin(base string, urlJoinPath urlJoinPathFunc, elem ...string) string {
	value, err := urlJoinPath(base, elem...)
	if err != nil {
		panic(err)
	}
	return value
}

// StartSpan adds an OpenTelemetry span to the trace with the given name.
func StartSpan(ctx context.Context, name string, attr ...attribute.KeyValue) context.Context {
	return startSpanWithOption(ctx, name, WithAttribute(attr...))
}

func EndSpan(ctx context.Context, err error) {
	span := trace.SpanFromContext(ctx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
	}
	span.End()
}

func SpanFromContext(ctx context.Context) trace.Span {
	return trace.SpanFromContext(ctx)
}

func startSpanWithOption(ctx context.Context, name string, opts ...StartSpanOption) context.Context {
	opt := &startSpanOption{
		urlJoinPathFunc: url.JoinPath,
	}
	for _, o := range opts {
		o.apply(opt)
	}
	var spanOpt []trace.SpanStartOption
	if len(opt.attributes) > 0 {
		spanOpt = append(spanOpt, trace.WithAttributes(opt.attributes...))
	}
	ctx, _ = otel.GetTracerProvider().Tracer(OpenTelemetryTracerName).Start(ctx, mustJoin(OpenTelemetryTracerName, opt.urlJoinPathFunc, name), spanOpt...)
	return ctx
}
