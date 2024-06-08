package context_tools

import (
	"context"
	"net/http"
)

type ContextKey string

func SetContextValue(r *http.Request, key ContextKey, value any) {
	r = r.WithContext(context.WithValue(r.Context(), key, value))
}

func GetContextValue[T any](r *http.Request, key ContextKey) T {
	return r.Context().Value(key).(T)
}
