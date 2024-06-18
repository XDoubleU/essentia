package contexttools_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/XDoubleU/essentia/pkg/contexttools"
	"github.com/stretchr/testify/assert"
)

func TestSetContextValue(t *testing.T) {
	r, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "", nil)

	r = contexttools.SetContextValue(r, contexttools.ShowErrorsContextKey, true)

	value := r.Context().Value(contexttools.ShowErrorsContextKey).(bool)

	assert.Equal(t, true, value)
}

func TestGetContextValue(t *testing.T) {
	ctx := context.WithValue(context.Background(), contexttools.ShowErrorsContextKey, true)
	r, _ := http.NewRequestWithContext(ctx, http.MethodGet, "", nil)

	value := contexttools.GetContextValue[bool](r, contexttools.ShowErrorsContextKey)

	assert.Equal(t, true, value)
}
