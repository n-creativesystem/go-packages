package interceptors

import (
	"net/http"

	"connectrpc.com/connect"
)

// モック用のStreamingHandlerConn実装
type mockStreamingConn struct {
	header  http.Header
	trailer http.Header
	spec    connect.Spec
}

func (m *mockStreamingConn) Spec() connect.Spec {
	return m.spec
}

func (m *mockStreamingConn) RequestHeader() http.Header {
	return m.header
}

func (m *mockStreamingConn) ResponseHeader() http.Header {
	return http.Header{}
}

func (m *mockStreamingConn) ResponseTrailer() http.Header {
	return m.trailer
}

func (m *mockStreamingConn) Peer() connect.Peer {
	return connect.Peer{
		Addr:     "127.0.0.1",
		Protocol: "test",
	}
}

func (m *mockStreamingConn) Receive(interface{}) error {
	return nil
}

func (m *mockStreamingConn) Send(interface{}) error {
	return nil
}

func (m *mockStreamingConn) Close() error {
	return nil
}
