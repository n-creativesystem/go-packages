package trace

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
}

func TestDefaultOption(t *testing.T) {
	assert.Equal(t, "github.com/n-creativesystem/go-packages", defaultOption.serviceName)
	assert.Equal(t, "local", defaultOption.environment)
}
