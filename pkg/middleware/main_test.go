package middleware_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/XDoubleU/essentia/internal/mocks"
	"github.com/XDoubleU/essentia/pkg/contexttools"
	"github.com/XDoubleU/essentia/pkg/logger"
	"github.com/XDoubleU/essentia/pkg/middleware"
	"github.com/getsentry/sentry-go"
	"github.com/stretchr/testify/assert"
)

func testCorsHeaders(
	t *testing.T,
	handler func(next http.Handler) http.Handler,
	method string,
	allowedOrigins []string,
	exHeaders []string,
) {
	t.Helper()

	expectedMethods := []string{"GET", "POST", "PATCH", "DELETE"}

	req, _ := http.NewRequest(method, "http://example.com/foo", nil)
	req.Header["Origin"] = []string{"http://example.com"}

	if method == http.MethodOptions {
		req.Header["Access-Control-Request-Method"] = expectedMethods
		req.Header["Access-Control-Request-Headers"] = exHeaders
	}

	res := testMiddleware(t, handler, req, nil)

	origins := res.Header()["Access-Control-Allow-Origin"]
	assert.Equal(t, allowedOrigins, origins)

	creds, _ := strconv.ParseBool(res.Header()["Access-Control-Allow-Credentials"][0])
	assert.True(t, creds)

	if method != http.MethodOptions {
		return
	}

	methods := res.Header()["Access-Control-Allow-Methods"]
	assert.Equal(t, expectedMethods, methods)

	headers := res.Header()["Access-Control-Allow-Headers"]
	assert.Equal(t, exHeaders, headers)
}

func testMiddleware(
	t *testing.T,
	middleware func(next http.Handler) http.Handler,
	req *http.Request,
	innerHandler func(w http.ResponseWriter, r *http.Request),
) *httptest.ResponseRecorder {
	t.Helper()

	var testResponse = []byte("bar")
	var testHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if innerHandler != nil {
			innerHandler(w, r)
		}
		_, _ = w.Write(testResponse)
		w.WriteHeader(http.StatusOK)
	})

	res := httptest.NewRecorder()
	middleware(testHandler).ServeHTTP(res, req)

	return res
}

func TestCors(t *testing.T) {
	allowedOrigins := []string{"http://example.com"}

	sentryHeaders := []string{"Content-Type", "Baggage", "Sentry-Trace"}
	noSentryHeaders := []string{"Content-Type"}

	corsSentry := middleware.Cors(allowedOrigins, true)
	corsNoSentry := middleware.Cors(allowedOrigins, false)

	testCorsHeaders(t, corsSentry, http.MethodGet, allowedOrigins, sentryHeaders)
	testCorsHeaders(t, corsSentry, http.MethodOptions, allowedOrigins, sentryHeaders)

	testCorsHeaders(t, corsNoSentry, http.MethodGet, allowedOrigins, noSentryHeaders)
	testCorsHeaders(
		t,
		corsNoSentry,
		http.MethodOptions,
		allowedOrigins,
		noSentryHeaders,
	)
}

func TestErrors(t *testing.T) {
	obfuscatedErrors := middleware.ErrorObfuscater(false)
	shownErrors := middleware.ErrorObfuscater(true)

	req, _ := http.NewRequest(http.MethodGet, "http://example.com/foo", nil)
	testMiddleware(
		t,
		obfuscatedErrors,
		req,
		func(_ http.ResponseWriter, r *http.Request) {
			assert.False(t, r.Context().Value(contexttools.ShowErrorsContextKey).(bool))
		},
	)
	testMiddleware(t, shownErrors, req, func(_ http.ResponseWriter, r *http.Request) {
		assert.True(t, r.Context().Value(contexttools.ShowErrorsContextKey).(bool))
	})
}

func TestLogger(t *testing.T) {
	mockedLogger := mocks.MockedLogger{}
	logger.SetLogger(mockedLogger.GetLogger())

	req, _ := http.NewRequest(http.MethodGet, "http://example.com/foo", nil)
	testMiddleware(
		t,
		middleware.Logger,
		req,
		func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
		},
	)

	timeStr := time.Now().Format("2006/01/02")
	assert.Contains(t, mockedLogger.Buffer.String(), timeStr)
	assert.Contains(t, mockedLogger.Buffer.String(), fmt.Sprintf("[%d]", http.StatusOK))
}

func TestRateLimit(t *testing.T) {
	bucketSize := 30
	rateLimiter := middleware.RateLimit(
		10,
		bucketSize,
		time.Second,
		500*time.Millisecond,
	)

	singleRequest := func() *httptest.ResponseRecorder {
		req, _ := http.NewRequest(http.MethodGet, "http://example.com/foo", nil)
		req.RemoteAddr = "127.0.0.1:80"

		return testMiddleware(t, rateLimiter, req, nil)
	}

	for i := 0; i < bucketSize; i++ {
		res := singleRequest()
		assert.Equal(t, http.StatusOK, res.Result().StatusCode)
	}
	time.Sleep(10 * time.Millisecond)

	res := singleRequest()
	assert.Equal(t, http.StatusTooManyRequests, res.Result().StatusCode)
	time.Sleep(2 * time.Second)

	res = singleRequest()
	assert.Equal(t, http.StatusOK, res.Result().StatusCode)
}

func TestRecover(t *testing.T) {
	mockedLogger := mocks.MockedLogger{}
	logger.SetLogger(mockedLogger.GetLogger())

	req, _ := http.NewRequest(http.MethodGet, "http://example.com/foo", nil)
	res := testMiddleware(
		t,
		middleware.Recover,
		req,
		func(_ http.ResponseWriter, _ *http.Request) {
			panic("test")
		},
	)

	assert.Equal(t, http.StatusInternalServerError, res.Result().StatusCode)
	assert.Equal(t, "close", res.Header()["Connection"][0])
	assert.Contains(t, mockedLogger.Buffer.String(), "PANIC")
}

func TestSentry(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://example.com/foo", nil)
	testMiddleware(
		t,
		middleware.Sentry(true, *mocks.GetMockedSentryClientOptions()),
		req,
		func(_ http.ResponseWriter, r *http.Request) {
			assert.NotNil(t, sentry.GetHubFromContext(r.Context()))
		},
	)
}
