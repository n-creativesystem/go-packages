package logging

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoggingHandle(t *testing.T) {
	// LoggingHandleのString()メソッドをテスト
	handle := JsonHandler
	assert.Equal(t, "json", handle.String())

	handle = TextHandler
	assert.Equal(t, "text", handle.String())

	handle = SentryHandler
	assert.Equal(t, "sentry", handle.String())

	handle = RollbarHandler
	assert.Equal(t, "rollbar", handle.String())

	handle = DatadogHandler
	assert.Equal(t, "datadog", handle.String())

	// カスタムハンドル名のテスト
	customHandle := LoggingHandle("custom")
	assert.Equal(t, "custom", customHandle.String())
}

func TestLoggingHandlersToInf(t *testing.T) {
	// LoggingHandlersToInf関数のテスト
	result := LoggingHandlersToInf()

	// 結果が配列であることを確認
	assert.IsType(t, []any{}, result)

	// すべてのハンドラーが含まれていることを確認
	assert.Len(t, result, len(LoggingHandlers))

	// 各ハンドラーが正しく変換されていることを確認
	for i, handle := range LoggingHandlers {
		assert.Equal(t, handle, result[i])
	}
}
