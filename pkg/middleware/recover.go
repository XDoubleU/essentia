package middleware

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/xdoubleu/essentia/internal/shared"
)

// Recover is middleware used to recover from a panic.
func Recover(logger *log.Logger) shared.Middleware {
	return func(next http.Handler) http.Handler {
		return recoverHandler(logger, next)
	}
}

func recoverHandler(logger *log.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("connection", "close")
				w.WriteHeader(http.StatusInternalServerError)
				logger.
					Printf("PANIC: %s\nstacktrace: %s\n", err, string(debug.Stack()))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
