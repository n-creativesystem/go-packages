package interceptors

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLoggingInterceptor(t *testing.T) {
	interceptor := NewLoggingInterceptor()
	require.NotNil(t, interceptor, "インターセプターがnilであってはならない")
}

func TestLoggingIntercept_WrapUnary(t *testing.T) {
	interceptor := NewLoggingInterceptor()

	// シンプルなコネクトハンドラを作成
	handler := func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		return connect.NewResponse(&struct{}{}), nil
	}

	wrappedFunc := interceptor.WrapUnary(handler)
	require.NotNil(t, wrappedFunc, "ラップされた関数がnilであってはならない")

	// HTTP リクエストとレスポンスを準備
	httpReq := httptest.NewRequest("POST", "/test.api.v1.TestService/TestMethod", nil)
	httpReq.Header.Set("x-request-id", "test-request-id")

	// connectのリクエストを作成
	req := connect.NewRequest(&struct{}{})
	req.Header().Set("x-request-id", "test-request-id")

	// ラップされた関数を実行
	ctx := context.Background()
	resp, err := wrappedFunc(ctx, req)

	// 結果を検証
	require.NoError(t, err, "エラーが発生してはならない")
	require.NotNil(t, resp, "レスポンスがnilであってはならない")
}

func TestLoggingIntercept_WrapStreamingHandler(t *testing.T) {
	interceptor := NewLoggingInterceptor()

	// シンプルなストリーミングハンドラを作成
	handler := func(ctx context.Context, conn connect.StreamingHandlerConn) error {
		// リクエストを処理してレスポンスを返す
		// 実際のテストではここでデータの送受信をシミュレートすることもできる
		return nil
	}

	wrappedFunc := interceptor.WrapStreamingHandler(handler)
	require.NotNil(t, wrappedFunc, "ラップされた関数がnilであってはならない")

	// モックのストリーミング接続を作成
	header := http.Header{}
	header.Set("x-request-id", "test-streaming-id")
	conn := &mockStreamingConn{
		header:  header,
		trailer: http.Header{},
		spec: connect.Spec{
			Procedure: "/test.api.v1.TestService/TestStreamingMethod",
		},
	}

	// ラップされた関数を実行
	ctx := context.Background()
	err := wrappedFunc(ctx, conn)

	// 結果を検証
	require.NoError(t, err, "エラーが発生してはならない")
}

func TestGetRequestId(t *testing.T) {
	t.Run("ヘッダーにリクエストIDが含まれる場合", func(t *testing.T) {
		header := http.Header{}
		expected := "existing-request-id"
		header.Set("x-request-id", expected)

		actual := getRequestId(header)
		assert.Equal(t, expected, actual, "ヘッダーから正しいリクエストIDを取得する必要がある")
	})

	t.Run("ヘッダーにリクエストIDが含まれない場合", func(t *testing.T) {
		header := http.Header{}

		actual := getRequestId(header)
		_, err := uuid.Parse(actual)
		assert.NoError(t, err, "有効なUUIDを生成する必要がある")
	})
}

func TestSecToTime(t *testing.T) {
	testCases := []struct {
		duration time.Duration
		expected string
	}{
		{1 * time.Second, "00:00:01.000"},
		{1 * time.Minute, "00:01:00.000"},
		{1 * time.Hour, "01:00:00.000"},
		{1500 * time.Millisecond, "00:00:01.500"},
	}

	for _, tc := range testCases {
		t.Run(tc.duration.String(), func(t *testing.T) {
			actual := secToTime(tc.duration)
			assert.Equal(t, tc.expected, actual, "時間の文字列形式が正しくない")
		})
	}
}
