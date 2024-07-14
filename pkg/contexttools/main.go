// Package contexttools provides functions which can be used to
// set/get values to/from [context.Context].
package contexttools

import (
	"context"
	"io"
	"log"
	"net/http"
)

// ContextKey is the type used for specifying context keys.
type ContextKey string

// SetContextValue sets a value by key on the context.
func SetContextValue(r *http.Request, key ContextKey, value any) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), key, value))
}

// GetContextValue returns a value by key from the context.
func GetContextValue[T any](r *http.Request, key ContextKey) *T {
	val := r.Context().Value(key)
	if val == nil {
		return nil
	}

	castedValue, ok := val.(T)
	if !ok {
		return nil
	}

	return &castedValue
}

// SetLogger sets the logger on the context.
func SetLogger(r *http.Request, logger *log.Logger) *http.Request {
	return SetContextValue(r, loggerContextKey, logger)
}

// GetLogger returns the logger stored in the context or a NullLogger.
func GetLogger(r *http.Request) *log.Logger {
	logger := GetContextValue[log.Logger](r, loggerContextKey)

	if logger == nil {
		return log.New(io.Discard, "", 0)
	}

	return logger
}

// SetShowErrors enables showing errors
// of [httptools.ServerErrorResponse].
func SetShowErrors(r *http.Request) *http.Request {
	return SetContextValue(r, showErrorsContextKey, true)
}

// GetShowErrors returns if errors should be shown
// in [httptools.ServerErrorResponse.].
func GetShowErrors(r *http.Request) bool {
	showErrors := GetContextValue[bool](r, showErrorsContextKey)

	if showErrors == nil {
		return false
	}

	return *showErrors
}
