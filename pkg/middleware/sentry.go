package middleware

import (
	"net/http"

	"github.com/XDoubleU/essentia/internal/mocks"
	"github.com/XDoubleU/essentia/pkg/httptools"
	"github.com/XDoubleU/essentia/pkg/logger"
	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
)

func Sentry(isTestEnv bool, clientOptions sentry.ClientOptions) middleware {
	sentryHandler := getSentryHandler(clientOptions)

	if isTestEnv {
		return func(next http.Handler) http.Handler {
			return sentryHandler.Handle(useMockedHub(enrichSentryHub(next)))
		}
	}

	return func(next http.Handler) http.Handler {
		return sentryHandler.Handle(enrichSentryHub(next))
	}
}

func getSentryHandler(clientOptions sentry.ClientOptions) *sentryhttp.Handler {
	err := sentry.Init(clientOptions)

	if err != nil {
		logger.GetLogger().Printf("sentry initialization failed: %v\n", err)
		return nil
	}

	return sentryhttp.New(sentryhttp.Options{
		Repanic: true,
	})
}

func enrichSentryHub(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := httptools.NewResponseWriter(w)
		next.ServeHTTP(rw, r)

		transaction := sentry.TransactionFromContext(r.Context())
		transaction.Status = sentry.HTTPtoSpanStatus(rw.StatusCode())
	})
}

func useMockedHub(next http.Handler) http.Handler {
	mockedHub := mocks.GetMockedSentryHub()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sentry.SetHubOnContext(r.Context(), mockedHub)
		next.ServeHTTP(w, r)
	})
}
