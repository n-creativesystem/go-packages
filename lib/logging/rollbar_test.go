package logging

import (
	"io"
	"log/slog"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockTransport struct {
	fn func(req *http.Request) (*http.Response, error)
}

func (tr *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return tr.fn(req)
}

func newMockTransport(fn func(req *http.Request) (*http.Response, error)) http.RoundTripper {
	return &mockTransport{fn: fn}
}

func TestRollbar(t *testing.T) {
	done := make(chan struct{})
	called := false
	defer func() {
		<-done
		assert.True(t, called)
	}()
	transport := func(req *http.Request) (*http.Response, error) {
		called = true
		done <- struct{}{}
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("")),
		}, nil
	}
	client := *http.DefaultClient
	client.Transport = newMockTransport(transport)
	conf := RollbarConfig{
		Level:  "error",
		Token:  "DUMMY",
		Client: &client,
	}
	conf.Init("local", "v1", "test")
	h := NewRollbarHandler(&conf)
	defer h.Close()
	log := slog.New(h)
	slog.SetDefault(log)
	slog.Error("This is test")
}
