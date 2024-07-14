// Package middleware provides configureable middleware and predefined lists,
// such as [Minimal] and [Default].
package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/XDoubleU/essentia/internal/mocks"
	"github.com/getsentry/sentry-go"
	"github.com/goddtriffin/helmet"
	"github.com/justinas/alice"
)

type middleware = func(next http.Handler) http.Handler

// Minimal provides a predefined chain of useful middleware.
// Being:
//   - [Logger]
//   - [Recover]
func Minimal(logger *log.Logger) []alice.Constructor {
	return []alice.Constructor{
		Logger(logger),
		Recover(logger),
	}
}

// Default provides a predefined chain of useful middleware.
// Being:
//   - All middleware from [Minimal]
//   - [helmet.Helmet]
//   - [CORS]
//   - [RateLimit]
//   - [Sentry]
func Default(
	logger *log.Logger,
	isTestEnv bool,
	allowedOrigins []string,
	sentryClientOptions *sentry.ClientOptions,
) ([]alice.Constructor, error) {
	if isTestEnv {
		sentryClientOptions = mocks.GetMockedSentryClientOptions()
	}

	useSentry := sentryClientOptions != nil

	helmet := helmet.Default()

	handlers := Minimal(logger)
	handlers = append(handlers, helmet.Secure)
	handlers = append(handlers, CORS(allowedOrigins, useSentry))
	//nolint:mnd//no magic number
	handlers = append(handlers, RateLimit(10, 30, time.Minute, 3*time.Minute))

	if useSentry {
		sentryMiddleware, err := Sentry(isTestEnv, *sentryClientOptions)
		if err != nil {
			return nil, err
		}

		handlers = append(handlers, sentryMiddleware)
	}

	return handlers, nil
}
