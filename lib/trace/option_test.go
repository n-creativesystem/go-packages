package trace

import (
	"testing"

	"github.com/stretchr/testify/assert"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func TestOptionApply(t *testing.T) {
	t.Run("WithEnabled", func(t *testing.T) {
		opt := &option{enabled: false}
		WithEnabled(true).apply(opt)
		assert.True(t, opt.enabled)
	})

	t.Run("WithAgentAddr", func(t *testing.T) {
		opt := &option{agentAddr: ""}
		WithAgentAddr("localhost:4317").apply(opt)
		assert.Equal(t, "localhost:4317", opt.agentAddr)
	})

	t.Run("WithServiceName", func(t *testing.T) {
		opt := &option{serviceName: "old-service"}
		WithServiceName("new-service").apply(opt)
		assert.Equal(t, "new-service", opt.serviceName)
	})

	t.Run("WithEnvironment", func(t *testing.T) {
		opt := &option{environment: "production"}
		WithEnvironment("staging").apply(opt)
		assert.Equal(t, "staging", opt.environment)
	})

	t.Run("WithVersion", func(t *testing.T) {
		opt := &option{version: "v1.0.0"}
		WithVersion("v1.1.0").apply(opt)
		assert.Equal(t, "v1.1.0", opt.version)
	})

	t.Run("WithSdkTraceOption", func(t *testing.T) {
		opt := &option{sdktraceOptions: []sdktrace.TracerProviderOption{}}
		sampler := sdktrace.WithSampler(sdktrace.AlwaysSample())
		WithSdkTraceOption(sampler).apply(opt)
		assert.Len(t, opt.sdktraceOptions, 1)
	})

	t.Run("WithSdkMetricOption", func(t *testing.T) {
		opt := &option{}
		WithSdkMetricOption(sdkmetric.WithReader(sdkmetric.NewManualReader())).apply(opt)
		assert.Len(t, opt.sdkmetricOptions, 1)
	})

	t.Run("WithSdkLogOption", func(t *testing.T) {
		opt := &option{}
		WithSdkLogOption(sdklog.WithAttributeCountLimit(10)).apply(opt)
		assert.Len(t, opt.sdklogOptions, 1)
	})
}

func TestNewOption_DoesNotMutateDefault(t *testing.T) {
	opt := newOption(WithEnabled(true), WithServiceName("svc"), WithSdkTraceOption(sdktrace.WithSampler(sdktrace.AlwaysSample())))
	assert.True(t, opt.enabled)
	assert.Equal(t, "svc", opt.serviceName)
	assert.Len(t, opt.sdktraceOptions, 1)

	assert.False(t, defaultOption.enabled)
	assert.Equal(t, "github.com/n-creativesystem/go-packages", defaultOption.serviceName)
	assert.Empty(t, defaultOption.sdktraceOptions)
}

func TestDefaultOption(t *testing.T) {
	assert.Equal(t, "github.com/n-creativesystem/go-packages", defaultOption.serviceName)
	assert.Equal(t, "local", defaultOption.environment)
}
