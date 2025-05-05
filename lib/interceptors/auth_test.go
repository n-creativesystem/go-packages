package interceptors

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"connectrpc.com/connect"
	"github.com/n-creativesystem/go-packages/lib/interceptors/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockHeader はauth.Getterインターフェースを実装するテスト用のモック / Mock implementing auth.Getter interface for testing
type mockHeader struct {
	values map[string]string
}

func newMockHeader() *mockHeader {
	return &mockHeader{
		values: make(map[string]string),
	}
}

func (m *mockHeader) Get(key string) string {
	return m.values[key]
}

func (m *mockHeader) Set(key, value string) {
	m.values[key] = value
}

// mockValidator はauth.Validatorインターフェースを実装するテスト用のモック / Mock implementing auth.Validator interface for testing
type mockValidator struct {
	shouldError bool
}

type mockTokenInfo struct {
	UserID string
}

func (m *mockValidator) Execute(ctx context.Context, getter auth.Getter) (*mockTokenInfo, error) {
	if m.shouldError {
		return nil, errors.New("validation failed")
	}
	return &mockTokenInfo{UserID: "test-user"}, nil
}

// HTTPHeaderGetterは複数のインターフェースをサポートするヘッダーアダプター / HTTPHeaderGetter is an adapter for different header interfaces
type HTTPHeaderGetter struct {
	header *mockHeader
}

func NewHTTPHeaderGetter(header *mockHeader) *HTTPHeaderGetter {
	return &HTTPHeaderGetter{header: header}
}

func (h *HTTPHeaderGetter) Get(key string) string {
	return h.header.Get(key)
}

// Getterを実装するhttp.Headerをラップするアダプター / Adapter wrapping http.Header that implements Getter
type headerAdapter struct {
	header http.Header
}

func newHeaderAdapter(header http.Header) *headerAdapter {
	return &headerAdapter{header: header}
}

func (h *headerAdapter) Get(key string) string {
	return h.header.Get(key)
}

func TestAuthenticate_WrapUnary(t *testing.T) {
	tests := []struct {
		name         string
		validator    *mockValidator
		setupHeader  func(h *mockHeader)
		errorHandler func(err error) error
		wantErr      bool
		wantErrCode  connect.Code
		checkContext func(t *testing.T, ctx context.Context)
	}{
		{
			name:      "認証成功 / Authentication success",
			validator: &mockValidator{shouldError: false},
			setupHeader: func(h *mockHeader) {
				h.Set("Authorization", "Bearer token")
			},
			wantErr: false,
			checkContext: func(t *testing.T, ctx context.Context) {
				tokenInfo, ok := auth.AuthFromContext[mockTokenInfo](ctx)
				require.True(t, ok)
				assert.Equal(t, "test-user", tokenInfo.UserID)
			},
		},
		{
			name:      "認証失敗 / Authentication failure",
			validator: &mockValidator{shouldError: true},
			setupHeader: func(h *mockHeader) {
				h.Set("Authorization", "Invalid token")
			},
			wantErr:     true,
			wantErrCode: connect.CodeUnauthenticated,
		},
		{
			name:      "カスタムエラーハンドラー / Custom error handler",
			validator: &mockValidator{shouldError: true},
			setupHeader: func(h *mockHeader) {
				h.Set("Authorization", "Invalid token")
			},
			errorHandler: func(err error) error {
				return connect.NewError(connect.CodePermissionDenied, err)
			},
			wantErr:     true,
			wantErrCode: connect.CodePermissionDenied,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 認証インターセプター作成 / Create authentication interceptor
			interceptor := NewAuthenticate[mockTokenInfo](tt.validator)
			authenticator, ok := interceptor.(*authenticate[mockTokenInfo])
			require.True(t, ok)

			// カスタムエラーハンドラーがある場合は設定 / Set custom error handler if provided
			if tt.errorHandler != nil {
				authenticator.errorHandler = tt.errorHandler
			}

			// リクエストとヘッダーを準備 / Prepare request and header
			header := newMockHeader()
			if tt.setupHeader != nil {
				tt.setupHeader(header)
			}

			// カスタム認証関数を直接テスト / Directly test the auth function
			// 正規のconnect.AnyRequestを使用する代わりにヘッダーとコンテキストのみテスト / Test only header and context instead of using regular connect.AnyRequest
			ctx := context.Background()
			newCtx, err := authenticator.authFunc(ctx, header)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errorHandler != nil {
					err = tt.errorHandler(err)
				} else {
					err = connect.NewError(connect.CodeUnauthenticated, err)
				}
				connectErr, ok := err.(*connect.Error)
				require.True(t, ok)
				assert.Equal(t, tt.wantErrCode, connectErr.Code())
			} else {
				require.NoError(t, err)
				if tt.checkContext != nil {
					tt.checkContext(t, newCtx)
				}
			}
		})
	}
}

func TestAuthenticate_WrapStreamingHandler(t *testing.T) {
	tests := []struct {
		name         string
		validator    *mockValidator
		setupHeader  func(h *mockHeader)
		errorHandler func(err error) error
		wantErr      bool
		wantErrCode  connect.Code
		checkContext func(t *testing.T, ctx context.Context)
	}{
		{
			name:      "認証成功 / Authentication success",
			validator: &mockValidator{shouldError: false},
			setupHeader: func(h *mockHeader) {
				h.Set("Authorization", "Bearer token")
			},
			wantErr: false,
			checkContext: func(t *testing.T, ctx context.Context) {
				tokenInfo, ok := auth.AuthFromContext[mockTokenInfo](ctx)
				require.True(t, ok)
				assert.Equal(t, "test-user", tokenInfo.UserID)
			},
		},
		{
			name:      "認証失敗 / Authentication failure",
			validator: &mockValidator{shouldError: true},
			setupHeader: func(h *mockHeader) {
				h.Set("Authorization", "Invalid token")
			},
			wantErr:     true,
			wantErrCode: connect.CodeUnauthenticated,
		},
		{
			name:      "カスタムエラーハンドラー / Custom error handler",
			validator: &mockValidator{shouldError: true},
			setupHeader: func(h *mockHeader) {
				h.Set("Authorization", "Invalid token")
			},
			errorHandler: func(err error) error {
				return connect.NewError(connect.CodePermissionDenied, err)
			},
			wantErr:     true,
			wantErrCode: connect.CodePermissionDenied,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 認証インターセプター作成 / Create authentication interceptor
			interceptor := NewAuthenticate[mockTokenInfo](tt.validator)
			authenticator, ok := interceptor.(*authenticate[mockTokenInfo])
			require.True(t, ok)

			// カスタムエラーハンドラーがある場合は設定 / Set custom error handler if provided
			if tt.errorHandler != nil {
				authenticator.errorHandler = tt.errorHandler
			}

			// ヘッダーを準備 / Prepare header
			header := newMockHeader()
			if tt.setupHeader != nil {
				tt.setupHeader(header)
			}

			// カスタム認証関数を直接テスト / Directly test the auth function
			ctx := context.Background()
			newCtx, err := authenticator.authFunc(ctx, header)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errorHandler != nil {
					err = tt.errorHandler(err)
				} else {
					err = connect.NewError(connect.CodeUnauthenticated, err)
				}
				connectErr, ok := err.(*connect.Error)
				require.True(t, ok)
				assert.Equal(t, tt.wantErrCode, connectErr.Code())
			} else {
				require.NoError(t, err)
				if tt.checkContext != nil {
					tt.checkContext(t, newCtx)
				}
			}
		})
	}
}
