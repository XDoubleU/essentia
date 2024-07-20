package middleware

import (
	"net/http"

	"github.com/XDoubleU/essentia/internal/mocks"
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
			return sentryHandler.Handle(useMockedHub(next))
		}, nil
	}

	return func(next http.Handler) http.Handler {
		return sentryHandler.Handle(next)
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

func useMockedHub(next http.Handler) http.Handler {
	mockedHub := mocks.MockedSentryHub()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sentry.SetHubOnContext(r.Context(), mockedHub)
		next.ServeHTTP(w, r)
	})
}
