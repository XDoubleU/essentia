package contexttools_test

import (
	"context"
	"io"
	"log"
	"testing"

	"github.com/XDoubleU/essentia/pkg/contexttools"
	"github.com/stretchr/testify/assert"
)

const testContextKey = contexttools.ContextKey("test")

func TestGetContextValue(t *testing.T) {
	ctx := context.WithValue(
		context.Background(),
		testContextKey,
		true,
	)

	value := contexttools.GetContextValue[bool](ctx, testContextKey)

	assert.Equal(t, true, *value)
}

func TestGetContextValueNotPresent(t *testing.T) {
	value := contexttools.GetContextValue[bool](context.Background(), testContextKey)

	assert.Nil(t, value)
}

func TestGetContextValueIncorrectType(t *testing.T) {
	ctx := context.WithValue(
		context.Background(),
		testContextKey,
		10,
	)

	value := contexttools.GetContextValue[bool](ctx, testContextKey)

	assert.Nil(t, value)
}

func TestSetGetLogger(t *testing.T) {
	ctx := context.Background()

	logger := log.Default()
	ctx = contexttools.WithLogger(ctx, logger)

	value := contexttools.Logger(ctx)

	assert.Equal(t, logger, value)
}

func TestGetNullLogger(t *testing.T) {
	ctx := context.Background()

	value := contexttools.Logger(ctx)

	assert.Equal(t, io.Discard, value.Writer())
}
