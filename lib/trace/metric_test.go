package trace

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric/noop"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

func TestNewMeterProvider_NoopProvider(t *testing.T) {
	ctx := context.Background()

	// enabled: falseの場合、noop.MeterProviderが返される
	mp, cleanup, err := NewMeterProvider(ctx, WithEnabled(false))
	require.NoError(t, err)
	defer cleanup()

	assert.IsType(t, noop.MeterProvider{}, mp)
	assert.Equal(t, mp, otel.GetMeterProvider())
}

func TestMeter(t *testing.T) {
	ctx := context.Background()

	reader := sdkmetric.NewManualReader()
	mp := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))
	oldProvider := otel.GetMeterProvider()
	otel.SetMeterProvider(mp)
	defer otel.SetMeterProvider(oldProvider)

	counter, err := Meter().Int64Counter("test.counter")
	require.NoError(t, err)
	counter.Add(ctx, 3)

	var rm metricdata.ResourceMetrics
	require.NoError(t, reader.Collect(ctx, &rm))
	require.Len(t, rm.ScopeMetrics, 1)
	assert.Equal(t, OpenTelemetryTracerName, rm.ScopeMetrics[0].Scope.Name)
	require.Len(t, rm.ScopeMetrics[0].Metrics, 1)

	sum, ok := rm.ScopeMetrics[0].Metrics[0].Data.(metricdata.Sum[int64])
	require.True(t, ok)
	require.Len(t, sum.DataPoints, 1)
	assert.Equal(t, int64(3), sum.DataPoints[0].Value)
}
