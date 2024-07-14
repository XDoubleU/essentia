package middleware

import (
	"net/http"

	"github.com/XDoubleU/essentia/pkg/contexttools"
)

// ShowErrors is middleware used to show errors.
// When used errors handled by [httptools.ServerErrorResponse] will be shown.
// Otherwise these will be hidden.
func ShowErrors() middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = contexttools.SetShowErrors(r)
			next.ServeHTTP(w, r)
		})
	}
}
