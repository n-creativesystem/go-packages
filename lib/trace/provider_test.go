package trace

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTracerProvider_NoopProvider(t *testing.T) {
	ctx := context.Background()

	// enabled: falseの場合、noop.TracerProviderが返される
	tp, cleanup, err := NewTracerProvider(ctx, WithEnabled(false))
	defer cleanup()

	assert.NoError(t, err)
	assert.NotNil(t, tp)
	assert.NotNil(t, cleanup)

	// クリーンアップ関数が正常に呼び出せることを確認
	cleanup()
}

func TestNewResource(t *testing.T) {
	serviceName := "test-service"
	version := "v1.0.0"
	environment := "test"

	resource := newResource(serviceName, version, environment)
	assert.NotNil(t, resource)
}

func TestNewPropagator(t *testing.T) {
	propagator := newPropagator()
	assert.NotNil(t, propagator)
}
