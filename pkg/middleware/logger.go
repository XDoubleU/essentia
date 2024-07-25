package middleware

import (
	"log/slog"
	"net/http"
	"time"

	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/xdoubleu/essentia/internal/shared"
	"github.com/xdoubleu/essentia/pkg/contexttools"
)

// Logger is middleware used to add a logger to
// the context and log every request and their duration.
func Logger(logger *slog.Logger) shared.Middleware {
	return func(next http.Handler) http.Handler {
		return loggerHandler(logger, next)
	}
}

func loggerHandler(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := sentryhttp.NewWrapResponseWriter(w, r.ProtoMajor)
		t := time.Now()

		r = r.WithContext(contexttools.WithLogger(r.Context(), logger))

		next.ServeHTTP(rw, r)

		logger.Info(
			"processed request",
			slog.Int("status", rw.Status()),
			slog.String("endpoint", r.RequestURI),
			slog.Duration("duration", time.Since(t)),
		)
	})
}
