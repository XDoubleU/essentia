package contexttools_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/XDoubleU/essentia/pkg/contexttools"
	"github.com/stretchr/testify/assert"
)

const testContextKey = contexttools.ContextKey("test")

func TestSetContextValue(t *testing.T) {
	r, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "", nil)

	r = contexttools.SetContextValue(r, testContextKey, true)

	value, _ := r.Context().Value(testContextKey).(bool)

	assert.Equal(t, true, value)
}

func TestGetContextValue(t *testing.T) {
	ctx := context.WithValue(
		context.Background(),
		testContextKey,
		true,
	)
	r, _ := http.NewRequestWithContext(ctx, http.MethodGet, "", nil)

	value := contexttools.GetContextValue[bool](r, testContextKey)

	assert.Equal(t, true, *value)
}

func TestGetContextValueNotPresent(t *testing.T) {
	r, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "", nil)

	value := contexttools.GetContextValue[bool](r, testContextKey)

	assert.Nil(t, value)
}

func TestGetContextValueIncorrectType(t *testing.T) {
	ctx := context.WithValue(
		context.Background(),
		testContextKey,
		10,
	)
	r, _ := http.NewRequestWithContext(ctx, http.MethodGet, "", nil)

	value := contexttools.GetContextValue[bool](r, testContextKey)

	assert.Nil(t, value)
}
