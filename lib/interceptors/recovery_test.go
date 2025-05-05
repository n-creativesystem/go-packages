package interceptors

import (
	"context"
	"net/http"
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRecovery(t *testing.T) {
	interceptor := NewRecovery()
	require.NotNil(t, interceptor, "インターセプターがnilであってはならない")
}

func TestRecovery_WrapUnary(t *testing.T) {
	interceptor := NewRecovery()

	t.Run("正常系", func(t *testing.T) {
		// 通常のハンドラーを作成
		handler := func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			return connect.NewResponse(&struct{}{}), nil
		}

		wrappedFunc := interceptor.WrapUnary(handler)
		require.NotNil(t, wrappedFunc, "ラップされた関数がnilであってはならない")

		// connectのリクエストを作成
		req := connect.NewRequest(&struct{}{})

		// ラップされた関数を実行
		ctx := context.Background()
		resp, err := wrappedFunc(ctx, req)

		// 結果を検証
		require.NoError(t, err, "エラーが発生してはならない")
		require.NotNil(t, resp, "レスポンスがnilであってはならない")
	})

	t.Run("パニック発生時", func(t *testing.T) {
		// パニックを発生させるハンドラーを作成
		handler := func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			panic("テスト用のパニック")
		}

		wrappedFunc := interceptor.WrapUnary(handler)
		require.NotNil(t, wrappedFunc, "ラップされた関数がnilであってはならない")

		// connectのリクエストを作成
		req := connect.NewRequest(&struct{}{})

		// ラップされた関数を実行
		ctx := context.Background()
		resp, err := wrappedFunc(ctx, req)

		// 結果を検証
		assert.Error(t, err, "エラーが発生するべき")
		assert.Equal(t, connect.CodeInternal, connect.CodeOf(err), "内部エラーのコードであるべき")
		assert.Nil(t, resp, "パニック時はレスポンスがnilであるべき")
	})
}

func TestRecovery_WrapStreamingHandler(t *testing.T) {
	interceptor := NewRecovery()

	t.Run("正常系", func(t *testing.T) {
		// 通常のストリーミングハンドラーを作成
		handler := func(ctx context.Context, conn connect.StreamingHandlerConn) error {
			return nil
		}

		wrappedFunc := interceptor.WrapStreamingHandler(handler)
		require.NotNil(t, wrappedFunc, "ラップされた関数がnilであってはならない")

		// モックのストリーミング接続を作成
		conn := &mockStreamingConn{
			header: make(http.Header),
			spec: connect.Spec{
				Procedure: "/test.api.v1.TestService/TestStreamingMethod",
			},
		}

		// ラップされた関数を実行
		ctx := context.Background()
		err := wrappedFunc(ctx, conn)

		// 結果を検証
		require.NoError(t, err, "エラーが発生してはならない")
	})

	t.Run("パニック発生時", func(t *testing.T) {
		// パニックを発生させるストリーミングハンドラーを作成
		handler := func(ctx context.Context, conn connect.StreamingHandlerConn) error {
			panic("テスト用のパニック")
		}

		wrappedFunc := interceptor.WrapStreamingHandler(handler)
		require.NotNil(t, wrappedFunc, "ラップされた関数がnilであってはならない")

		// モックのストリーミング接続を作成
		conn := &mockStreamingConn{
			header: make(http.Header),
			spec: connect.Spec{
				Procedure: "/test.api.v1.TestService/TestStreamingMethod",
			},
		}

		// ラップされた関数を実行
		ctx := context.Background()
		err := wrappedFunc(ctx, conn)

		// 結果を検証
		assert.Error(t, err, "エラーが発生するべき")
		assert.Equal(t, connect.CodeInternal, connect.CodeOf(err), "内部エラーのコードであるべき")
	})
}

func TestRecovery_WrapStreamingClient(t *testing.T) {
	interceptor := NewRecovery()

	// StreamingClientはパススルーのみで特別な処理をしないため、
	// 単純に実行して問題がないことを確認
	handler := func(ctx context.Context, spec connect.Spec) connect.StreamingClientConn {
		return nil // テストのため実際には何も返さない
	}

	wrappedFunc := interceptor.WrapStreamingClient(handler)
	require.NotNil(t, wrappedFunc, "ラップされた関数がnilであってはならない")

	// 簡易的な実行確認
	ctx := context.Background()
	spec := connect.Spec{Procedure: "/test.api.v1.TestService/TestStreamingMethod"}
	conn := wrappedFunc(ctx, spec)

	assert.Nil(t, conn, "このテスト環境では、接続はnilであるべき")
}
