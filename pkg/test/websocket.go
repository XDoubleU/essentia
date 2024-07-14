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

// WebsocketTester is used for testing a websocket.
type WebsocketTester struct {
	handler           http.Handler
	timeout           time.Duration
	sleep             time.Duration
	initialMsg        any
	parallelOperation ParallelOperation
}

// ParallelOperation is a operation executed while
// waiting for a new message coming from the server.
type ParallelOperation = func(t *testing.T, ts *httptest.Server, conn *websocket.Conn)

// CreateWebsocketTester creates a new [WebsocketTester].
func CreateWebsocketTester(handler http.Handler) WebsocketTester {
	return WebsocketTester{
		handler:           handler,
		timeout:           10 * time.Second, //nolint:mnd //no magic number
		sleep:             time.Second,
		initialMsg:        nil,
		parallelOperation: nil,
	}
}

// SetInitialMessage sets the message to be sent when the connection is created.
func (tWeb *WebsocketTester) SetInitialMessage(msg any) {
	tWeb.initialMsg = msg
}

// SetParallelOperation sets operation to be executed while
// waiting for a new message coming from the server.
func (tWeb *WebsocketTester) SetParallelOperation(parallelOperation ParallelOperation) {
	tWeb.parallelOperation = parallelOperation
}

// Do executes a [WebsocketTester] returning the response
// of the initialRequest and parallelOperation.
// An error is returned if the initialResponse can't be read.
func (tWeb WebsocketTester) Do(
	t *testing.T,
	initialResponse any,
	parallelOperationResponse any,
) error {
	t.Helper()

	var err error

	ts := httptest.NewServer(tWeb.handler)
	defer ts.Close()

	wsURL := "ws" + strings.TrimPrefix(ts.URL, "http")
	ws, err := dialWebsocket(wsURL, tWeb.timeout)
	require.Nil(t, err)

	defer func() {
		err = ws.CloseNow()
		require.Nil(t, err)
	}()

	if tWeb.initialMsg != nil {
		ctx, cancel := context.WithTimeout(context.Background(), tWeb.timeout)
		defer cancel()
		err = wsjson.Write(ctx, ws, tWeb.initialMsg)
		require.Nil(t, err)
	}

	if initialResponse != nil {
		ctx, cancel := context.WithTimeout(context.Background(), tWeb.timeout)
		defer cancel()
		err = wsjson.Read(ctx, ws, &initialResponse)

		if err != nil {
			return err
		}
	}

	if tWeb.parallelOperation != nil {
		go func() {
			time.Sleep(tWeb.sleep)
			tWeb.parallelOperation(t, ts, ws)
		}()

		ctx, cancel := context.WithTimeout(context.Background(), tWeb.timeout)
		defer cancel()
		err = wsjson.Read(ctx, ws, &parallelOperationResponse)
		require.Nil(t, err)
	}

	return nil
}

func dialWebsocket(url string, timeout time.Duration) (*websocket.Conn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	//nolint:bodyclose //don't close ws bodies
	ws, _, err := websocket.Dial(ctx, url, nil)
	if err != nil {
		return nil, err
	}

	return ws, nil
}
