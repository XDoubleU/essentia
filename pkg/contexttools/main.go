package contexttools

import (
	"context"
	"net/http"
)

type ContextKey string

func SetContextValue(r *http.Request, key ContextKey, value any) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), key, value))
}

func GetContextValue[T any](r *http.Request, key ContextKey) *T {
	val := r.Context().Value(key)
	if val == nil {
		return nil
	}

	castedValue := val.(T)
	return &castedValue
}
