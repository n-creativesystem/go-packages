package trace

import (
	"context"
	"errors"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace/noop"
)

func TestMustJoin(t *testing.T) {
	tests := []struct {
		name     string
		base     string
		elements []string
		want     string
	}{
		{
			name:     "single element",
			base:     "base",
			elements: []string{"path"},
			want:     "base/path",
		},
		{
			name:     "multiple elements",
			base:     "base",
			elements: []string{"path1", "path2", "path3"},
			want:     "base/path1/path2/path3",
		},
		{
			name:     "empty base",
			base:     "",
			elements: []string{"path"},
			want:     "path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mustJoin(tt.base, url.JoinPath, tt.elements...)
			assert.Equal(t, tt.want, result)
		})
	}

	// モンキーパッチを使ってurl.JoinPathを一時的に置き換えてパニックをテスト
	t.Run("invalid path should panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic but didn't get one")
			}
		}()

		// モンキーパッチ: 常にエラーを返す関数に置き換え
		urlJoinPath := func(base string, elem ...string) (string, error) {
			return "", errors.New("mock error for testing panic")
		}

		// これは必ずパニックするはず
		_ = mustJoin("base", urlJoinPath, "path")
	})
}

func TestStartSpan(t *testing.T) {
	// まずトレーサープロバイダーを無効なものにセット
	originalProvider := otel.GetTracerProvider()
	defer otel.SetTracerProvider(originalProvider)

	// 無効なトレーサーを使ってテスト
	otel.SetTracerProvider(noop.NewTracerProvider())

	ctx := context.Background()
	attrKey := attribute.Key("test-key")
	attrVal := "test-value"
	attrs := []attribute.KeyValue{attrKey.String(attrVal)}

	// スパンの作成
	spanCtx := StartSpan(ctx, "test-span", attrs...)
	assert.NotNil(t, spanCtx)

	// 属性を持つスパン
	spanCtxWithAttrs := StartSpan(ctx, "test-span-with-attrs", attrs...)
	assert.NotNil(t, spanCtxWithAttrs)
}

func TestEndSpan(t *testing.T) {
	// 無効なトレーサーをセット
	originalProvider := otel.GetTracerProvider()
	defer otel.SetTracerProvider(originalProvider)
	otel.SetTracerProvider(noop.NewTracerProvider())

	ctx := context.Background()
	testErr := errors.New("test error")

	// エラーなしでスパンを終了
	spanCtx := StartSpan(ctx, "test-span")
	EndSpan(spanCtx, nil)

	// エラーありでスパンを終了
	spanCtx = StartSpan(ctx, "test-span-with-error")
	EndSpan(spanCtx, testErr)
}

func TestSpanFromContext(t *testing.T) {
	// 無効なトレーサーをセット
	originalProvider := otel.GetTracerProvider()
	defer otel.SetTracerProvider(originalProvider)
	otel.SetTracerProvider(noop.NewTracerProvider())

	ctx := context.Background()

	// スパンなしのコンテキスト
	span := SpanFromContext(ctx)
	assert.NotNil(t, span) // 常に非nilの値が返る

	// スパンありのコンテキスト
	spanCtx := StartSpan(ctx, "test-span")
	span = SpanFromContext(spanCtx)
	assert.NotNil(t, span)
}
