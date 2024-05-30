package middleware

import (
	"net/http"

	"github.com/XDoubleU/essentia/pkg/http_tools"
	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
)

func Sentry(clientOptions sentry.ClientOptions) middleware {
	sentryHandler := getSentryHandler(clientOptions)

	return func(next http.Handler) http.Handler {
		return enrichSentryHub(sentryHandler.Handle(next))
	}
}

func getSentryHandler(clientOptions sentry.ClientOptions) *sentryhttp.Handler {
	err := sentry.Init(clientOptions)

	if err != nil {
		//todo: app.logger.Printf("sentry initialization failed: %v\n", err)
		return nil
	}

	return sentryhttp.New(sentryhttp.Options{
		Repanic: true,
	})
}

func enrichSentryHub(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := http_tools.NewResponseWriter(w)
		next.ServeHTTP(rw, r)

		transaction := sentry.TransactionFromContext(r.Context())
		transaction.Status = sentry.HTTPtoSpanStatus(rw.StatusCode())
	})
}
