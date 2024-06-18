package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/XDoubleU/essentia/pkg/middleware"
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

	var testResponse = []byte("bar")
	var testHandler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(testResponse)
	})

	req, _ := http.NewRequest(method, "http://example.com/foo", nil)
	req.Header["Origin"] = []string{"http://example.com"}

	if method == http.MethodOptions {
		req.Header["Access-Control-Request-Method"] = expectedMethods
		req.Header["Access-Control-Request-Headers"] = exHeaders
	}

	res := httptest.NewRecorder()
	handler(testHandler).ServeHTTP(res, req)

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
