package trace

import (
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/log/noop"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type recordProcessor struct {
	records []sdklog.Record
}

func (p *recordProcessor) OnEmit(_ context.Context, record *sdklog.Record) error {
	p.records = append(p.records, record.Clone())
	return nil
}

func (p *recordProcessor) Enabled(context.Context, sdklog.EnabledParameters) bool { return true }
func (p *recordProcessor) Shutdown(context.Context) error                         { return nil }
func (p *recordProcessor) ForceFlush(context.Context) error                       { return nil }

func TestNewLoggerProvider_NoopProvider(t *testing.T) {
	ctx := context.Background()

	// enabled: falseの場合、noop.LoggerProviderが返される
	lp, cleanup, err := NewLoggerProvider(ctx, WithEnabled(false))
	require.NoError(t, err)
	defer cleanup()

	assert.IsType(t, noop.LoggerProvider{}, lp)
	assert.Equal(t, lp, global.GetLoggerProvider())
}

func TestNewSlogHandler(t *testing.T) {
	processor := &recordProcessor{}
	lp := sdklog.NewLoggerProvider(sdklog.WithProcessor(processor))
	oldProvider := global.GetLoggerProvider()
	global.SetLoggerProvider(lp)
	defer global.SetLoggerProvider(oldProvider)

	tp := sdktrace.NewTracerProvider()
	ctx, span := tp.Tracer("test").Start(context.Background(), "test-span")
	defer span.End()

	logger := slog.New(NewSlogHandler())
	logger.InfoContext(ctx, "hello", slog.String("key", "value"))

	require.Len(t, processor.records, 1)
	record := processor.records[0]
	assert.Equal(t, "hello", record.Body().AsString())
	assert.Equal(t, OpenTelemetryTracerName, record.InstrumentationScope().Name)
	// span の trace_id / span_id がログに紐付けられている
	assert.Equal(t, span.SpanContext().TraceID(), record.TraceID())
	assert.Equal(t, span.SpanContext().SpanID(), record.SpanID())
}
