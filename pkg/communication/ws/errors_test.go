package ws_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	wstools "github.com/xdoubleu/essentia/pkg/communication/ws"
	"github.com/xdoubleu/essentia/pkg/config"
	errortools "github.com/xdoubleu/essentia/pkg/errors"
	sentrytools "github.com/xdoubleu/essentia/pkg/sentry"
	"github.com/xdoubleu/essentia/pkg/test"
)

func testErrorStatusCode(t *testing.T, handler http.HandlerFunc) int {
	t.Helper()

	req, _ := http.NewRequest(http.MethodGet, "", nil)
	res := httptest.NewRecorder()

	sentryMiddleware, err := sentrytools.Middleware(
		config.TestEnv,
		sentrytools.MockedSentryClientOptions(),
	)
	require.Nil(t, err)

	sentryMiddleware(handler).ServeHTTP(res, req)

	return res.Result().StatusCode
}

func setupWS(t *testing.T) http.Handler {
	t.Helper()

	wsHandler := wstools.CreateWebSocketHandler[TestSubscribeMsg](
		1,
		10,
		[]string{"http://localhost"},
	)
	_, err := wsHandler.AddTopic("topic", nil)
	require.Nil(t, err)

	sentryMiddleware, err := sentrytools.Middleware(
		config.TestEnv,
		sentrytools.MockedSentryClientOptions(),
	)
	require.Nil(t, err)

	return sentryMiddleware(wsHandler.Handler())
}

func TestUpgradeErrorResponse(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		wstools.UpgradeErrorResponse(w, r, errors.New("test"))
	}

	statusCode := testErrorStatusCode(t, handler)

	assert.Equal(t, http.StatusInternalServerError, statusCode)
}

func TestErrorResponse(t *testing.T) {
	handler := setupWS(t)

	tWeb := test.CreateWebSocketTester(handler)
	tWeb.SetInitialMessage(TestSubscribeMsg{TopicName: "unknown"})

	var errorDto errortools.ErrorDto
	err := tWeb.Do(t, &errorDto, nil)
	require.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, errorDto.Status)
	assert.Equal(t, http.StatusText(http.StatusBadRequest), errorDto.Error)
	assert.Equal(t, "topic 'unknown' doesn't exist", errorDto.Message)
}
