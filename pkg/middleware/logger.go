package middleware

import (
	"log"
	"net/http"
	"time"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()

		next.ServeHTTP(w, r)

		log.Printf(
			"[%d] %s in %v",
			0,
			//TODO w.StatusCode(),
			r.RequestURI,
			time.Since(t),
		)
	})
}
