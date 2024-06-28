package middleware

import (
	"net/http"
	"time"

	"github.com/XDoubleU/essentia/internal/mocks"
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

func Default(
	isTestEnv bool,
	allowedOrigins []string,
	sentryClientOptions *sentry.ClientOptions,
	showErrors bool,
) []alice.Constructor {
	if isTestEnv {
		sentryClientOptions = mocks.GetMockedSentryClientOptions()
	}

	useSentry := sentryClientOptions != nil

	helmet := helmet.Default()

	handlers := Minimal(showErrors)
	handlers = append(handlers, helmet.Secure)
	handlers = append(handlers, Cors(allowedOrigins, useSentry))
	//nolint:mnd//no magic number
	handlers = append(handlers, RateLimit(10, 30, time.Minute, 3*time.Minute))

	if useSentry {
		handlers = append(handlers, Sentry(isTestEnv, *sentryClientOptions))
	}

	return handlers
}
