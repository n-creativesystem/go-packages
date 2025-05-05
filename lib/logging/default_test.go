package logging

import (
	"context"
	"log/slog"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvLogLevel(t *testing.T) {
	// 元の環境変数を保存
	origLogLevel := os.Getenv("LOG_LEVEL")
	defer os.Setenv("LOG_LEVEL", origLogLevel)

	// 環境変数のテストケース
	testCases := []struct {
		env      string
		expected string
	}{
		{"DEBUG", "DEBUG"},
		{"INFO", "INFO"},
		{"WARN", "WARN"},
		{"WARNING", "WARN"}, // "WARNING"は"WARN"に変換される
		{"ERROR", "ERROR"},
		{"", "INFO"},        // 未設定の場合はINFO
		{"UNKNOWN", "INFO"}, // 不明な値の場合はINFO
	}

	for _, tc := range testCases {
		t.Run("ENV_"+tc.env, func(t *testing.T) {
			os.Setenv("LOG_LEVEL", tc.env)
			level := envLogLevel()
			assert.Equal(t, strings.ToUpper(tc.expected), strings.ToUpper(level.String()))
		})
	}
}

func TestDefaultLoggerFunctions(t *testing.T) {
	buf := new(strings.Builder)
	defaultLogger = slog.New(NewProcessHandler(slog.NewJSONHandler(buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))

	// 各レベルのログ関数をテスト
	Debug("Debug message")
	Info("Info message")
	Warn("Warning message")
	Error("Error message")

	// コンテキスト付きのロガー関数をテスト
	ctx := context.Background()
	DebugContext(ctx, "Debug with context")
	InfoContext(ctx, "Info with context")
	WarnContext(ctx, "Warn with context")
	ErrorContext(ctx, "Error with context")

	// 出力内容の確認
	output := buf.String()
	assert.Contains(t, output, "Debug message")
	assert.Contains(t, output, "Info message")
	assert.Contains(t, output, "Warning message")
	assert.Contains(t, output, "Error message")
	assert.Contains(t, output, "Debug with context")
	assert.Contains(t, output, "Info with context")
	assert.Contains(t, output, "Warn with context")
	assert.Contains(t, output, "Error with context")
}
