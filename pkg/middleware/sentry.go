package middleware

import (
	"net/http"

	"github.com/XDoubleU/essentia/internal/mocks"
	"github.com/XDoubleU/essentia/pkg/httptools"
	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
)

// Sentry is middleware used to configure and enable Sentry.
// When isTestEnv is true, a mocked [sentry.Hub] will be used.
func Sentry(isTestEnv bool, clientOptions sentry.ClientOptions) (middleware, error) {
	sentryHandler, err := getSentryHandler(clientOptions)
	if err != nil {
		return nil, err
	}

	if isTestEnv {
		return func(next http.Handler) http.Handler {
			return sentryHandler.Handle(useMockedHub(enrichSentryHub(next)))
		}, nil
	}

	return func(next http.Handler) http.Handler {
		return sentryHandler.Handle(enrichSentryHub(next))
	}, nil
}

func getSentryHandler(clientOptions sentry.ClientOptions) (*sentryhttp.Handler, error) {
	err := sentry.Init(clientOptions)

	if err != nil {
		return nil, err
	}

	//nolint:exhaustruct //other fields are optional
	return sentryhttp.New(sentryhttp.Options{
		Repanic: true,
	}), nil
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
	mockedHub := mocks.MockedSentryHub()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sentry.SetHubOnContext(r.Context(), mockedHub)
		next.ServeHTTP(w, r)
	})
}
