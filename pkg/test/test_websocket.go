package test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type TestWebsocket struct {
	handler           http.Handler
	timeout           time.Duration
	sleep             time.Duration
	initialMsg        any
	parallelOperation ParallelOperation
}

type ParallelOperation = func(t *testing.T, ts *httptest.Server)

func CreateTestWebsocket(handler http.Handler) TestWebsocket {
	return TestWebsocket{
		handler:    handler,
		timeout:    10 * time.Second,
		sleep:      time.Second,
		initialMsg: nil,
	}
}

func (tWeb *TestWebsocket) SetInitialMessage(msg any) {
	tWeb.initialMsg = msg
}

func (tWeb *TestWebsocket) SetParallelOperation(parallelOperation ParallelOperation) {
	tWeb.parallelOperation = parallelOperation
}

func (tWeb TestWebsocket) Do(t *testing.T, initialResponse any, parallelOperationResponse any) {
	t.Helper()

	var err error

	ts := httptest.NewServer(tWeb.handler)
	defer ts.Close()

	wsURL := "ws" + strings.TrimPrefix(ts.URL, "http")
	ws, err := dialWebsocket(wsURL, tWeb.timeout)
	require.Nil(t, err)

	if tWeb.initialMsg != nil {
		ctx, cancel := context.WithTimeout(context.Background(), tWeb.timeout)
		defer cancel()
		err := wsjson.Write(ctx, ws, tWeb.initialMsg)
		require.Nil(t, err)
	}

	if initialResponse != nil {
		ctx, cancel := context.WithTimeout(context.Background(), tWeb.timeout)
		defer cancel()
		err = wsjson.Read(ctx, ws, &initialResponse)
		require.Nil(t, err)
	}

	go func() {
		time.Sleep(tWeb.sleep)
		tWeb.parallelOperation(t, ts)
	}()

	ctx, cancel := context.WithTimeout(context.Background(), tWeb.timeout)
	defer cancel()
	err = wsjson.Read(ctx, ws, &parallelOperationResponse)
	require.Nil(t, err)
}

func dialWebsocket(url string, timeout time.Duration) (*websocket.Conn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ws, _, err := websocket.Dial(ctx, url, nil)
	if err != nil {
		return nil, err
	}

	return ws, nil
}