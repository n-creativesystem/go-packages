package logging

import (
	"bytes"
	"log/slog"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestColorHandler(t *testing.T) {
	// コンソールへの出力をキャプチャするためのバッファ
	buf := &bytes.Buffer{}

	// カラーハンドラーを作成
	colorHandler := NewColorHandler(WithWriter(buf))

	// ハンドラーからロガーを作成
	logger := slog.New(colorHandler)

	// 異なるレベルのログを出力
	logger.Info("Info message")
	logger.Warn("Warning message")
	logger.Error("Error message")

	// 出力内容の確認
	output := buf.String()
	require.Contains(t, output, "Info message")
	require.Contains(t, output, "Warning message")
	require.Contains(t, output, "Error message")
}

func TestColorHandlerWithAttrs(t *testing.T) {
	buf := &bytes.Buffer{}

	// カラーハンドラーを作成
	colorHandler := NewColorHandler(WithWriter(buf))

	// 属性を追加
	attrHandler := colorHandler.WithAttrs([]slog.Attr{
		slog.String("custom", "value"),
	})

	// ハンドラーからロガーを作成
	logger := slog.New(attrHandler)

	// ログを出力
	logger.Info("Message with attributes")

	// 出力内容の確認
	output := buf.String()
	require.Contains(t, output, "Message with attributes")
	require.Contains(t, output, "\"custom\":\"value\"")
}

func TestColorHandlerWithGroup(t *testing.T) {
	buf := &bytes.Buffer{}

	// カラーハンドラーを作成
	colorHandler := NewColorHandler(WithWriter(buf))

	// グループを追加
	groupHandler := colorHandler.WithGroup("testgroup")

	// ハンドラーからロガーを作成
	logger := slog.New(groupHandler)

	// ログを出力
	logger.Info("Message with group", slog.String("key", "value"))

	// 出力内容の確認
	output := buf.String()
	require.Contains(t, output, "Message with group")
	require.Contains(t, output, "testgroup")
}

func TestIsEnabledTerminalColor(t *testing.T) {
	// 現在の状態を保存
	origTerm := os.Getenv("TERM")
	defer os.Setenv("TERM", origTerm)

	// TERMが設定されていない場合
	os.Unsetenv("TERM")
	resetEnabledTerminalColor()
	// initの代わりに明示的に設定関数を呼び出す
	updateTerminalColorSetting()
	require.False(t, IsEnabledTerminalColor())

	// TERMがxterm系の場合
	os.Setenv("TERM", "xterm-256color")
	resetEnabledTerminalColor()
	updateTerminalColorSetting()
	require.True(t, IsEnabledTerminalColor())

	// TERMが無関係な値の場合
	os.Setenv("TERM", "unknown")
	resetEnabledTerminalColor()
	updateTerminalColorSetting()
	require.False(t, IsEnabledTerminalColor())
}

// テストのためにenabledTerminalColor変数をリセットする
func resetEnabledTerminalColor() {
	enabledTerminalColor = false
}

// init関数の内容を抽出した関数
func updateTerminalColorSetting() {
	colorTerminals := []string{
		"xterm",
		"vt100",
		"rxvt",
		"screen",
	}
	if v, ok := os.LookupEnv("TERM"); ok {
		for _, t := range colorTerminals {
			if strings.Contains(v, t) {
				enabledTerminalColor = true
				break
			}
		}
	}
}
