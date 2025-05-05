package logging

import (
	"bytes"
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func TestDatadogHandler(t *testing.T) {
	// テスト用のバッファ
	buf := &bytes.Buffer{}

	// Datadog 引数の設定
	ddArgs := DDArgs{
		ServiceName: "test-service",
		Environment: "test",
		Version:     "1.0.0",
	}

	// ベースハンドラーの作成
	baseHandler := NewTextHandler(WithWriter(buf))

	// Datadog ハンドラーの作成
	ddHandler := NewDatadogHandler(ddArgs, baseHandler)

	// ハンドラーからロガーを作成
	logger := slog.New(ddHandler)

	// テスト用のトレースプロバイダーとトレーサーを設定
	tp := sdktrace.NewTracerProvider()
	defer func() {
		_ = tp.Shutdown(context.Background())
	}()

	// グローバルのプロパゲーターとトレーサーを設定
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))
	oldProvider := otel.GetTracerProvider()
	otel.SetTracerProvider(tp)
	defer otel.SetTracerProvider(oldProvider)

	// テスト用のスパンを作成
	ctx := context.Background()
	tracer := tp.Tracer("test-tracer")
	ctx, span := tracer.Start(ctx, "test-span")
	defer span.End()

	// ログを出力
	logger.InfoContext(ctx, "Test log message")

	// 出力内容の検証
	output := buf.String()
	require.Contains(t, output, "dd.trace_id=")
	require.Contains(t, output, "dd.span_id=")
	require.Contains(t, output, "dd.service=test-service")
	require.Contains(t, output, "dd.env=test")
	require.Contains(t, output, "dd.version=1.0.0")
}

func TestConvertTraceID(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "valid hex id",
			input:    "0000000000000000abcdef1234567890",
			expected: "12379813812177893520", // decimal of abcdef1234567890
		},
		{
			name:     "short id",
			input:    "1234567890",
			expected: "",
		},
		{
			name:     "invalid hex",
			input:    "0000000000000000abcdefghijklmno",
			expected: "",
		},
		{
			name:     "empty id",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertTraceID(tt.input)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestDatadogHandlerWithoutSpan(t *testing.T) {
	buf := &bytes.Buffer{}

	ddArgs := DDArgs{
		ServiceName: "test-service",
		Environment: "test",
		Version:     "1.0.0",
	}

	baseHandler := NewTextHandler(WithWriter(buf))
	ddHandler := NewDatadogHandler(ddArgs, baseHandler)
	logger := slog.New(ddHandler)

	// スパンなしのコンテキストでログを出力
	ctx := context.Background()
	logger.InfoContext(ctx, "Test log message without span")

	// スパン関連の属性なしでログが出力されていることを確認
	output := buf.String()
	require.Contains(t, output, "Test log message without span")
}

func TestDatadogHandlerWithGroup(t *testing.T) {
	buf := &bytes.Buffer{}

	ddArgs := DDArgs{
		ServiceName: "test-service",
		Environment: "test",
		Version:     "1.0.0",
	}

	baseHandler := NewTextHandler(WithWriter(buf))
	ddHandler := NewDatadogHandler(ddArgs, baseHandler)

	// グループ付きのハンドラー
	groupHandler := ddHandler.WithGroup("testgroup")
	logger := slog.New(groupHandler)

	logger.Info("Test group message")

	output := buf.String()
	require.Contains(t, output, "testgroup")
	require.Contains(t, output, "Test group message")
}

func TestDatadogHandlerWithAttrs(t *testing.T) {
	buf := &bytes.Buffer{}

	ddArgs := DDArgs{
		ServiceName: "test-service",
		Environment: "test",
		Version:     "1.0.0",
	}

	baseHandler := NewTextHandler(WithWriter(buf))
	ddHandler := NewDatadogHandler(ddArgs, baseHandler)

	// 属性付きのハンドラー
	attrHandler := ddHandler.WithAttrs([]slog.Attr{slog.String("attr", "value")})
	logger := slog.New(attrHandler)

	logger.Info("Test attr message")

	output := buf.String()
	require.Contains(t, output, "attr=value")
	require.Contains(t, output, "Test attr message")
}

func TestDatadogHandlerClose(t *testing.T) {
	// モック実装
	mockCloseCalled := false
	mockHandler := &mockCloseHandler{
		closeFn: func() error {
			mockCloseCalled = true
			return nil
		},
	}

	ddArgs := DDArgs{
		ServiceName: "test-service",
		Environment: "test",
		Version:     "1.0.0",
	}

	ddHandler := NewDatadogHandler(ddArgs, mockHandler)

	// Close を呼び出す
	err := ddHandler.Close()

	// モックの Close が呼ばれたことを確認
	require.NoError(t, err)
	require.True(t, mockCloseCalled)
}

// モック用のハンドラー
type mockCloseHandler struct {
	closeFn func() error
}

func (m *mockCloseHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return true
}

func (m *mockCloseHandler) Handle(ctx context.Context, record slog.Record) error {
	return nil
}

func (m *mockCloseHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return m
}

func (m *mockCloseHandler) WithGroup(name string) slog.Handler {
	return m
}

func (m *mockCloseHandler) Close() error {
	return m.closeFn()
}
