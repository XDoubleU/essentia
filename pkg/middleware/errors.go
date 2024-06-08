package middleware

import (
	"net/http"

	"github.com/XDoubleU/essentia/pkg/context_tools"
)

func ErrorObfuscater(showErrors bool) middleware {
	return func(next http.Handler) http.Handler {
		return obfuscateErrors(showErrors, next)
	}
}

func obfuscateErrors(showErrors bool, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context_tools.SetContextValue(r, context_tools.ShowErrorsContextKey, showErrors)
		next.ServeHTTP(w, r)
	})
}
