package middleware

import (
	"net/http"
	"time"

	"github.com/XDoubleU/essentia/pkg/httptools"
	"github.com/XDoubleU/essentia/pkg/logger"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := httptools.NewResponseWriter(w)
		t := time.Now()

		next.ServeHTTP(rw, r)

		logger.GetLogger().Printf(
			"[%d] %s in %v",
			rw.StatusCode(),
			r.RequestURI,
			time.Since(t),
		)
	})
}
