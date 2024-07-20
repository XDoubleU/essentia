package test_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/XDoubleU/essentia/pkg/httptools"
	"github.com/XDoubleU/essentia/pkg/test"
	"github.com/XDoubleU/essentia/pkg/validate"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type TestResponse struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

type TestSubscribeMsg struct {
	Topic string `json:"topic"`
}

func (s TestSubscribeMsg) Validate() *validate.Validator {
	return validate.New()
}

func (s TestSubscribeMsg) GetTopic() string {
	return s.Topic
}

func setup(t *testing.T) http.Handler {
	t.Helper()

	ws := httptools.CreateWebsocketHandler[TestSubscribeMsg](
		[]string{"http://localhost"},
	)
	ws.AddTopicHandler(
		"exists",
		func(
			_ http.ResponseWriter,
			r *http.Request,
			conn *websocket.Conn,
			_ TestSubscribeMsg) {
			err := wsjson.Write(
				r.Context(),
				conn,
				TestResponse{Ok: true, Message: "initial"},
			)
			require.Nil(t, err)

			var msg TestSubscribeMsg
			err = wsjson.Read(r.Context(), conn, &msg)
			require.Nil(t, err)

			err = wsjson.Write(
				r.Context(),
				conn,
				TestResponse{Ok: true, Message: "parallel"},
			)
			require.Nil(t, err)
		},
	)
	return ws.Handler()
}

func TestWebsocketTester(t *testing.T) {
	tWeb := test.CreateWebsocketTester(setup(t))
	tWeb.SetInitialMessage(TestSubscribeMsg{
		Topic: "exists",
	})
	tWeb.SetParallelOperation(
		func(t *testing.T, _ *httptest.Server, ws *websocket.Conn) {
			err := wsjson.Write(context.Background(), ws, TestSubscribeMsg{
				Topic: "exists",
			})
			require.Nil(t, err)
		},
	)

	var initialResponse, parallelResponse TestResponse
	err := tWeb.Do(t, &initialResponse, &parallelResponse)

	require.Nil(t, err)
	assert.Equal(t, TestResponse{Ok: true, Message: "initial"}, initialResponse)
	assert.Equal(t, TestResponse{Ok: true, Message: "parallel"}, parallelResponse)
}
