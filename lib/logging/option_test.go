package logging

import (
	"bytes"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithOptions(t *testing.T) {
	// バッファを用意
	buf := &bytes.Buffer{}

	// オプションを設定
	opts := defaultOptions(
		WithWriter(buf),
		WithLevel(slog.LevelDebug),
		WithReplaceAttr(func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == "test_key" {
				return slog.String(a.Key, "modified")
			}
			return a
		}),
	)

	// オプションが正しく適用されていることを確認
	assert.Equal(t, buf, opts.writer)
	assert.Equal(t, slog.LevelDebug, opts.level)

	// ReplaceAttr関数が働くことを確認
	attr := slog.String("test_key", "original")
	modifiedAttr := opts.replaceAttr([]string{}, attr)
	assert.Equal(t, "modified", modifiedAttr.Value.String())
}

func TestWithHandler(t *testing.T) {
	// モックハンドラー
	mockHandler := &mockHandler{enabled: true}

	// オプションを設定
	opts := defaultOptions(WithHandler(mockHandler))

	// ハンドラーが正しく設定されていることを確認
	assert.Equal(t, mockHandler, opts.handler)
}

func TestDefaultOption(t *testing.T) {
	// デフォルトオプションの確認
	assert.Equal(t, os.Stdout, defaultOption.writer)
	assert.Equal(t, slog.LevelInfo, defaultOption.level)
	assert.NotNil(t, defaultOption.handler)
	assert.NotNil(t, defaultOption.replaceAttr)

	// デフォルトのreplaceAttr関数のテスト
	timeAttr := slog.String("time", "original")
	modifiedAttr := defaultOption.replaceAttr([]string{}, timeAttr)
	assert.Equal(t, "time", modifiedAttr.Key)
	assert.NotEqual(t, "original", modifiedAttr.Value.String())
}

func TestDefaultReplaceAttr(t *testing.T) {
	// 時間フィールドの置換
	timeAttr := slog.String("time", "original")
	modifiedTimeAttr := defaultReplaceAttr([]string{}, timeAttr)
	assert.Equal(t, "time", modifiedTimeAttr.Key)
	assert.NotEqual(t, "original", modifiedTimeAttr.Value.String())

	// 他のフィールドはそのまま
	otherAttr := slog.String("other", "value")
	modifiedOtherAttr := defaultReplaceAttr([]string{}, otherAttr)
	assert.Equal(t, otherAttr, modifiedOtherAttr)
}
