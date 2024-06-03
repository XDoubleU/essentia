package middleware

import (
	"net/http"
	"time"

	"github.com/XDoubleU/essentia/pkg/http_tools"
	"github.com/XDoubleU/essentia/pkg/logger"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := http_tools.NewResponseWriter(w)
		t := time.Now()

		next.ServeHTTP(w, r)

		logger.GetLogger().Printf(
			"[%d] %s in %v",
			rw.StatusCode(),
			r.RequestURI,
			time.Since(t),
		)
	})
}
