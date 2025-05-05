package auth

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type authInfo struct{}

func TestValidContext(t *testing.T) {
	ctx := SetContext[authInfo](t.Context(), &authInfo{})
	info, ok := AuthFromContext[authInfo](ctx)
	require.True(t, ok)
	require.NotNil(t, info)
}

func TestInvalidContext(t *testing.T) {
	info, ok := AuthFromContext[authInfo](t.Context())
	require.False(t, ok)
	require.Nil(t, info)
}
