package logging

import (
	"bytes"
	"log/slog"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProcessHandler(t *testing.T) {
	buf := &bytes.Buffer{}

	// ベースハンドラーの作成
	baseHandler := NewTextHandler(WithWriter(buf))

	// プロセスハンドラーの作成
	processHandler := NewProcessHandler(baseHandler)

	// ハンドラーからロガーを作成
	logger := slog.New(processHandler)

	// ログを出力
	logger.Info("Test process log")

	// 期待される値
	pid := os.Getpid()
	pidStr := strconv.Itoa(pid)

	// 出力内容の検証
	output := buf.String()
	require.Contains(t, output, "pid="+pidStr)
	require.Contains(t, output, "Test process log")
}

func TestProcessHandlerWithAttrs(t *testing.T) {
	buf := &bytes.Buffer{}

	// ベースハンドラーの作成
	baseHandler := NewTextHandler(WithWriter(buf))

	// プロセスハンドラーの作成
	processHandler := NewProcessHandler(baseHandler)

	// 属性付きハンドラー
	attrHandler := processHandler.WithAttrs([]slog.Attr{slog.String("custom", "value")})

	// ハンドラーからロガーを作成
	logger := slog.New(attrHandler)

	// ログを出力
	logger.Info("Test process with attrs")

	// 出力内容の検証
	output := buf.String()
	require.Contains(t, output, "custom=value")
	require.Contains(t, output, "pid=")
	require.Contains(t, output, "Test process with attrs")
}

func TestProcessHandlerWithGroup(t *testing.T) {
	buf := &bytes.Buffer{}

	// ベースハンドラーの作成
	baseHandler := NewTextHandler(WithWriter(buf))

	// プロセスハンドラーの作成
	processHandler := NewProcessHandler(baseHandler)

	// グループ付きハンドラー
	groupHandler := processHandler.WithGroup("testgroup")

	// ハンドラーからロガーを作成
	logger := slog.New(groupHandler)

	// ログを出力
	logger.Info("Test process with group")

	// 出力内容の検証
	output := buf.String()
	require.Contains(t, output, "testgroup")
	require.Contains(t, output, "pid=")
	require.Contains(t, output, "Test process with group")
}
