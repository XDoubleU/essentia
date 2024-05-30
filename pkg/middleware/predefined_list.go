package middleware

import (
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/goddtriffin/helmet"
	"github.com/justinas/alice"
)

type middleware = func(next http.Handler) http.Handler

func Minimal() []alice.Constructor {
	return []alice.Constructor{
		Logger,
		Recover,
	}
}

func Default(throttle bool, allowedOrigins []string, sentryClientOptions *sentry.ClientOptions) []alice.Constructor {
	useSentry := sentryClientOptions != nil

	helmet := helmet.Default()

	handlers := Minimal()
	handlers = append(handlers, helmet.Secure)
	handlers = append(handlers, Cors(allowedOrigins, useSentry))

	if throttle {
		handlers = append(handlers, RateLimit)
	}

	if useSentry {
		handlers = append(handlers, Sentry(*sentryClientOptions))
	}

	return handlers
}
