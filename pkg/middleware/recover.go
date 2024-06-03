package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/XDoubleU/essentia/pkg/logger"
)

func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				w.WriteHeader(http.StatusInternalServerError)
				logger.GetLogger().Printf("PANIC: %s\nstacktrace: %s\n", err, string(debug.Stack()))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
