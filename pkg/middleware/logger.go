package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/XDoubleU/essentia/pkg/contexttools"
	"github.com/XDoubleU/essentia/pkg/httptools"
)

// Logger is middleware used to add a logger to
// the context and log every request and their duration.
func Logger(logger *log.Logger) middleware {
	return func(next http.Handler) http.Handler {
		return loggerHandler(logger, next)
	}
}

func loggerHandler(logger *log.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := httptools.NewResponseWriter(w)
		t := time.Now()

		r = contexttools.SetLogger(r, logger)

		next.ServeHTTP(rw, r)

		logger.Printf(
			"[%d] %s in %v",
			rw.StatusCode(),
			r.RequestURI,
			time.Since(t),
		)
	})
}
