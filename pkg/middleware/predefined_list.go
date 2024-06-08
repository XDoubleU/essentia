package middleware

import (
	"net/http"

	"github.com/XDoubleU/essentia/internal/sentry_mock"
	"github.com/getsentry/sentry-go"
	"github.com/goddtriffin/helmet"
	"github.com/justinas/alice"
)

type middleware = func(next http.Handler) http.Handler

func Minimal(showErrors bool) []alice.Constructor {

	return []alice.Constructor{
		Logger,
		Recover,
		ErrorObfuscater(showErrors),
	}
}

func Default(isTestEnv bool, allowedOrigins []string, sentryClientOptions *sentry.ClientOptions, showErrors bool) []alice.Constructor {
	if isTestEnv {
		sentryClientOptions = sentry_mock.GetMockedClientOptions()
	}

	useSentry := sentryClientOptions != nil

	helmet := helmet.Default()

	handlers := Minimal(showErrors)
	handlers = append(handlers, helmet.Secure)
	handlers = append(handlers, Cors(allowedOrigins, useSentry))
	handlers = append(handlers, RateLimit)

	if useSentry {
		handlers = append(handlers, Sentry(isTestEnv, *sentryClientOptions))
	}

	return handlers
}
