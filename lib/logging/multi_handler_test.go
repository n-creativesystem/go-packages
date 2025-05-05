package logging

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMultiHandler(t *testing.T) {
	// 複数のバッファを用意
	buf1 := &bytes.Buffer{}
	buf2 := &bytes.Buffer{}

	// 複数のハンドラーを作成
	handler1 := slog.NewTextHandler(buf1, &slog.HandlerOptions{})
	handler2 := slog.NewTextHandler(buf2, &slog.HandlerOptions{})

	// マルチハンドラーの作成
	multiHandler := NewHandler(handler1, handler2)

	// ハンドラーからロガーを作成
	logger := slog.New(multiHandler)

	// ログを出力
	logger.Info("Test multi handler")

	// 両方のバッファに出力されていることを確認
	require.Contains(t, buf1.String(), "Test multi handler")
	require.Contains(t, buf2.String(), "Test multi handler")
}

func TestMultiHandlerEnabled(t *testing.T) {
	// テスト用のモックハンドラー
	mock1 := &mockHandler{enabled: true}
	mock2 := &mockHandler{enabled: false}
	mock3 := &mockHandler{enabled: true}

	// マルチハンドラーの作成
	multiHandler := NewHandler(mock1, mock2, mock3)

	// マルチハンドラーがENABLEDを返すことを確認（少なくとも一つのハンドラーがenabledなら）
	ctx := context.Background()
	require.True(t, multiHandler.Enabled(ctx, slog.LevelInfo))

	// すべてのハンドラーが無効な場合
	multiHandler2 := NewHandler(
		&mockHandler{enabled: false},
		&mockHandler{enabled: false},
	)
	require.False(t, multiHandler2.Enabled(ctx, slog.LevelInfo))
}

func TestMultiHandlerHandle(t *testing.T) {
	// エラーを返すハンドラー
	errHandler := &mockHandler{
		handleErr: errors.New("handle error"),
		enabled:   true,
	}

	// 正常に処理するハンドラー
	okHandler := &mockHandler{
		enabled: true,
	}

	// マルチハンドラーの作成
	multiHandler := NewHandler(errHandler, okHandler)

	// ハンドルメソッドを呼び出し、エラーがマージされることを確認
	ctx := context.Background()
	record := slog.Record{}
	record.Level = slog.LevelInfo

	err := multiHandler.Handle(ctx, record)
	require.Error(t, err)
	require.Contains(t, err.Error(), "handle error")
}

func TestMultiHandlerWithAttrs(t *testing.T) {
	buf1 := &bytes.Buffer{}
	buf2 := &bytes.Buffer{}

	handler1 := slog.NewTextHandler(buf1, &slog.HandlerOptions{})
	handler2 := slog.NewTextHandler(buf2, &slog.HandlerOptions{})

	multiHandler := NewHandler(handler1, handler2)

	// 属性を追加
	attrHandler := multiHandler.WithAttrs([]slog.Attr{
		slog.String("attr1", "value1"),
	})

	logger := slog.New(attrHandler)
	logger.Info("Test with attrs")

	// 両方のハンドラーが属性付きで出力されていることを確認
	require.Contains(t, buf1.String(), "attr1=value1")
	require.Contains(t, buf2.String(), "attr1=value1")
}

func TestMultiHandlerWithGroup(t *testing.T) {
	buf1 := &bytes.Buffer{}
	buf2 := &bytes.Buffer{}

	handler1 := slog.NewTextHandler(buf1, &slog.HandlerOptions{})
	handler2 := slog.NewTextHandler(buf2, &slog.HandlerOptions{})

	multiHandler := NewHandler(handler1, handler2)

	// グループを追加
	groupHandler := multiHandler.WithGroup("testgroup")

	logger := slog.New(groupHandler)
	logger.Info("Test with group", slog.String("key", "value"))

	// 両方のハンドラーがグループ付きで出力されていることを確認
	require.Contains(t, buf1.String(), "testgroup")
	require.Contains(t, buf2.String(), "testgroup")
}

func TestMultiHandlerClose(t *testing.T) {
	// クローズ呼び出しを追跡
	closeCalled1 := false
	closeCalled2 := false

	// クローザブルなモックハンドラー
	mockCloser1 := &mockCloseableHandler{
		closeFn: func() error {
			closeCalled1 = true
			return nil
		},
	}

	mockCloser2 := &mockCloseableHandler{
		closeFn: func() error {
			closeCalled2 = true
			return errors.New("close error")
		},
	}

	// マルチハンドラーの作成
	multiHandler := NewHandler(mockCloser1, mockCloser2)

	// クローズを呼び出し
	err := multiHandler.Close()

	// 両方のクローズが呼び出されたことを確認
	require.True(t, closeCalled1)
	require.True(t, closeCalled2)

	// エラーは無視されることを確認
	require.NoError(t, err)
}

// モックハンドラー実装
type mockHandler struct {
	enabled   bool
	handleErr error
}

func (m *mockHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return m.enabled
}

func (m *mockHandler) Handle(_ context.Context, _ slog.Record) error {
	return m.handleErr
}

func (m *mockHandler) WithAttrs([]slog.Attr) slog.Handler {
	return m
}

func (m *mockHandler) WithGroup(string) slog.Handler {
	return m
}

// クローズ可能なモックハンドラー
type mockCloseableHandler struct {
	mockHandler
	closeFn func() error
}

func (m *mockCloseableHandler) Close() error {
	return m.closeFn()
}
