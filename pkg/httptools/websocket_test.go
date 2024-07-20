package httptools_test

import (
	"net/http"
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
	Ok bool `json:"ok"`
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

func emptyHandler(
	_ http.ResponseWriter,
	_ *http.Request,
	_ *websocket.Conn,
	_ TestSubscribeMsg,
) {
}

func setup(t *testing.T, onCloseCallBackIsCalled *bool) http.Handler {
	t.Helper()

	ws := httptools.CreateWebsocketHandler[TestSubscribeMsg](
		[]string{"http://localhost"},
	)
	ws.SetOnCloseCallback(func(_ *websocket.Conn) {
		*onCloseCallBackIsCalled = true
	})
	ws.AddTopicHandler(
		"exists",
		func(
			_ http.ResponseWriter,
			r *http.Request,
			conn *websocket.Conn,
			_ TestSubscribeMsg) {
			err := wsjson.Write(r.Context(), conn, TestResponse{Ok: true})
			require.Nil(t, err)
		},
	)
	return ws.Handler()
}

func TestWebSocketExistingSubject(t *testing.T) {
	onCloseCallbackIsCalled := false
	wsHandler := setup(t, &onCloseCallbackIsCalled)

	tWeb := test.CreateWebsocketTester(wsHandler)
	tWeb.SetInitialMessage(TestSubscribeMsg{Topic: "exists"})

	var initialResponse TestResponse
	err := tWeb.Do(t, &initialResponse, nil)

	require.Nil(t, err)
	assert.True(t, initialResponse.Ok)
	assert.True(t, onCloseCallbackIsCalled)
}

func TestWebSocketUnknownSubject(t *testing.T) {
	onCloseCallbackIsCalled := false
	wsHandler := setup(t, &onCloseCallbackIsCalled)

	tWeb := test.CreateWebsocketTester(wsHandler)
	tWeb.SetInitialMessage(TestSubscribeMsg{Topic: "unknown"})

	var initialResponse httptools.ErrorDto
	err := tWeb.Do(t, &initialResponse, nil)

	require.Nil(t, err)
	assert.Equal(t, "unknown topic", initialResponse.Error)
	assert.Equal(t, "no handler found for 'unknown'", initialResponse.Message)
	assert.True(t, onCloseCallbackIsCalled)
}

func TestWebSocketExistingHandler(t *testing.T) {
	ws := httptools.CreateWebsocketHandler[TestSubscribeMsg]([]string{"localhost"})
	ws.AddTopicHandler(
		"exists",
		emptyHandler,
	)

	assert.PanicsWithValue(t, "topic 'exists' already has a handler", func() {
		ws.AddTopicHandler(
			"exists",
			emptyHandler,
		)
	})
}
