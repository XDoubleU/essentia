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

type TestSubjectMsg struct {
	Subject string `json:"subject"`
}

func (s TestSubjectMsg) Validate() *validate.Validator {
	return validate.New()
}

func (s TestSubjectMsg) GetSubject() string {
	return s.Subject
}

func setup(t *testing.T) http.Handler {
	t.Helper()

	ws := httptools.CreateWebsocketHandler[TestSubjectMsg]([]string{"http://localhost"})
	ws.AddSubjectHandler(
		"exists",
		func(
			_ http.ResponseWriter,
			r *http.Request,
			conn *websocket.Conn,
			_ TestSubjectMsg) {
			err := wsjson.Write(
				r.Context(),
				conn,
				TestResponse{Ok: true, Message: "initial"},
			)
			require.Nil(t, err)

			var msg TestSubjectMsg
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
	return ws.GetHandler()
}

func TestWebsocketTester(t *testing.T) {
	tWeb := test.CreateWebsocketTester(setup(t))
	tWeb.SetInitialMessage(TestSubjectMsg{
		Subject: "exists",
	})
	tWeb.SetParallelOperation(
		func(t *testing.T, _ *httptest.Server, ws *websocket.Conn) {
			err := wsjson.Write(context.Background(), ws, TestSubjectMsg{
				Subject: "exists",
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
