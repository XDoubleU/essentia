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

type TestSubjectMsg struct {
	Subject string `json:"subject"`
}

func (s TestSubjectMsg) Validate() *validate.Validator {
	return validate.New()
}

func (s TestSubjectMsg) GetSubject() string {
	return s.Subject
}

func emptyHandler(
	_ http.ResponseWriter,
	_ *http.Request,
	_ *websocket.Conn,
	_ TestSubjectMsg,
) {
}

func setup(t *testing.T, onCloseCallBackIsCalled *bool) http.Handler {
	t.Helper()

	ws := httptools.CreateWebsocketHandler[TestSubjectMsg]("http://localhost")
	ws.SetOnCloseCallback(func(_ *websocket.Conn) {
		*onCloseCallBackIsCalled = true
	})
	ws.AddSubjectHandler(
		"exists",
		func(
			_ http.ResponseWriter,
			r *http.Request,
			conn *websocket.Conn,
			_ TestSubjectMsg) {
			err := wsjson.Write(r.Context(), conn, TestResponse{Ok: true})
			require.Nil(t, err)
		},
	)
	return ws.GetHandler()
}

func TestWebSocketExistingSubject(t *testing.T) {
	onCloseCallbackIsCalled := false
	wsHandler := setup(t, &onCloseCallbackIsCalled)

	tWeb := test.CreateWebsocketTester(wsHandler)
	tWeb.SetInitialMessage(TestSubjectMsg{Subject: "exists"})

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
	tWeb.SetInitialMessage(TestSubjectMsg{Subject: "unknown"})

	var initialResponse httptools.ErrorDto
	err := tWeb.Do(t, &initialResponse, nil)

	require.Nil(t, err)
	assert.Equal(t, "unknown subject", initialResponse.Error)
	assert.Equal(t, "no handler found for 'unknown'", initialResponse.Message)
	assert.True(t, onCloseCallbackIsCalled)
}

func TestWebSocketExistingHandler(t *testing.T) {
	ws := httptools.CreateWebsocketHandler[TestSubjectMsg]("localhost")
	ws.AddSubjectHandler(
		"exists",
		emptyHandler,
	)

	assert.PanicsWithValue(t, "subject and handler already in map", func() {
		ws.AddSubjectHandler(
			"exists",
			emptyHandler,
		)
	})
}
