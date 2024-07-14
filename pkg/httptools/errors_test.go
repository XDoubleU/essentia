package httptools_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/XDoubleU/essentia/internal/mocks"
	"github.com/XDoubleU/essentia/pkg/contexttools"
	"github.com/XDoubleU/essentia/pkg/httptools"
	"github.com/XDoubleU/essentia/pkg/middleware"
	"github.com/XDoubleU/essentia/pkg/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"nhooyr.io/websocket"
)

func testErrorStatusCode(t *testing.T, handler http.HandlerFunc) int {
	t.Helper()

	req, _ := http.NewRequest(http.MethodGet, "", nil)
	res := httptest.NewRecorder()

	sentryMiddleware, err := middleware.Sentry(
		true,
		*mocks.GetMockedSentryClientOptions(),
	)
	require.Nil(t, err)

	sentryMiddleware(handler).ServeHTTP(res, req)

	return res.Result().StatusCode
}

func testError(t *testing.T, handler http.HandlerFunc) (int, httptools.ErrorDto) {
	t.Helper()

	req, _ := http.NewRequest(http.MethodGet, "", nil)
	return testErrorWithReq(t, handler, req)
}

func testErrorWithReq(
	t *testing.T,
	handler http.HandlerFunc,
	req *http.Request,
) (int, httptools.ErrorDto) {
	t.Helper()

	res := httptest.NewRecorder()

	sentryMiddleware, err := middleware.Sentry(
		true,
		*mocks.GetMockedSentryClientOptions(),
	)
	require.Nil(t, err)

	sentryMiddleware(handler).ServeHTTP(res, req)

	var errorDto httptools.ErrorDto
	err = httptools.ReadJSON(res.Result().Body, &errorDto)
	require.Nil(t, err)

	return res.Result().StatusCode, errorDto
}

func testErrorWS(
	t *testing.T,
	handler func(
		w http.ResponseWriter,
		r *http.Request,
		conn *websocket.Conn,
		msg TestSubjectMsg),
) {
	t.Helper()

	wsHandler := httptools.CreateWebsocketHandler[TestSubjectMsg]("localhost")
	wsHandler.AddSubjectHandler("subject", handler)

	sentryMiddleware, err := middleware.Sentry(
		true,
		*mocks.GetMockedSentryClientOptions(),
	)
	require.Nil(t, err)

	tWeb := test.CreateWebsocketTester(
		sentryMiddleware(wsHandler.GetHandler()),
	)
	tWeb.SetInitialMessage(TestSubjectMsg{Subject: "subject"})

	var errorDto httptools.ErrorDto
	err = tWeb.Do(t, &errorDto, nil)

	assert.ErrorContains(t, err, "status = StatusInternalError")
	assert.ErrorContains(
		t,
		err,
		"reason = \"the server encountered a problem and could not process your request\"",
	)
}

func TestServerErrorResponseObfuscated(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		httptools.ServerErrorResponse(w, r, errors.New("test"))
	}

	statusCode, errorDto := testError(t, handler)

	assert.Equal(t, http.StatusInternalServerError, statusCode)
	assert.Equal(t, httptools.MessageInternalServerError, errorDto.Message)
}

func TestServerErrorResponseShown(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		httptools.ServerErrorResponse(w, r, errors.New("test"))
	}

	req, _ := http.NewRequest(http.MethodGet, "", nil)
	req = contexttools.SetShowErrors(req)
	statusCode, errorDto := testErrorWithReq(t, handler, req)

	assert.Equal(t, http.StatusInternalServerError, statusCode)
	assert.Equal(t, "test", errorDto.Message)
}

func TestBadRequestResponse(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		httptools.BadRequestResponse(w, r, errors.New("test"))
	}

	statusCode, errorDto := testError(t, handler)

	assert.Equal(t, http.StatusBadRequest, statusCode)
	assert.Equal(t, "test", errorDto.Message)
}

func TestRateLimitExceededResponse(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		httptools.RateLimitExceededResponse(w, r)
	}

	statusCode, errorDto := testError(t, handler)

	assert.Equal(t, http.StatusTooManyRequests, statusCode)
	assert.Equal(t, httptools.MessageTooManyRequests, errorDto.Message)
}

func TestUnauthorizedResponse(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		httptools.UnauthorizedResponse(w, r, "test")
	}

	statusCode, errorDto := testError(t, handler)

	assert.Equal(t, http.StatusUnauthorized, statusCode)
	assert.Equal(t, "test", errorDto.Message)
}

func TestForbiddenResponse(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		httptools.ForbiddenResponse(w, r)
	}

	statusCode, errorDto := testError(t, handler)

	assert.Equal(t, http.StatusForbidden, statusCode)
	assert.Equal(t, httptools.MessageForbidden, errorDto.Message)
}

func TestConflictResponse(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		httptools.ConflictResponse(
			w,
			r,
			httptools.ErrResourceUniqueValue,
			"resource",
			"value",
			"field",
		)
	}

	statusCode, errorDto := testError(t, handler)

	assert.Equal(t, http.StatusConflict, statusCode)
	assert.Equal(t, map[string]any{
		"field": "resource with field 'value' already exists",
	}, errorDto.Message)
}

func TestNotFoundResponse(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		httptools.NotFoundResponse(
			w,
			r,
			httptools.ErrResourceNotFound,
			"resource",
			"value",
			"field",
		)
	}

	statusCode, errorDto := testError(t, handler)

	assert.Equal(t, http.StatusNotFound, statusCode)
	assert.Equal(t, map[string]any{
		"field": "resource with field 'value' doesn't exist",
	}, errorDto.Message)
}

func TestFailedValidationResponse(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		httptools.FailedValidationResponse(w, r, map[string]string{
			"field": "invalid value",
		})
	}

	statusCode, errorDto := testError(t, handler)

	assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	assert.Equal(t, map[string]any{
		"field": "invalid value",
	}, errorDto.Message)
}

func TestWSUpgradeErrorResponse(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		httptools.WSUpgradeErrorResponse(w, r, errors.New("test"))
	}

	statusCode := testErrorStatusCode(t, handler)

	assert.Equal(t, http.StatusInternalServerError, statusCode)
}

func TestWSErrorResponse(t *testing.T) {
	handler := func(
		w http.ResponseWriter,
		r *http.Request,
		conn *websocket.Conn,
		_ TestSubjectMsg) {
		httptools.WSErrorResponse(
			w,
			r,
			conn,
			func(_ *websocket.Conn) {},
			errors.New("test"),
		)
	}

	testErrorWS(t, handler)
}
