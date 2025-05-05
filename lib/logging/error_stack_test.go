package logging

import (
	"errors"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWithStack(t *testing.T) {
	// エラーありの場合
	err := errors.New("test error")
	attr := WithStack(err)

	require.Equal(t, "stack", attr.Key)
	require.Contains(t, attr.Value.String(), "test error")
	require.Contains(t, attr.Value.String(), "TestWithStack")

	// nilエラーの場合
	nilAttr := WithStack(nil)
	require.Equal(t, slog.Attr{}, nilAttr)
}
