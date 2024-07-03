package goroutine_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/XDoubleU/essentia/pkg/goroutine"
	"github.com/getsentry/sentry-go"
	"github.com/stretchr/testify/assert"
)

func TestSentryErrorHandler(t *testing.T) {
	name := "test"

	testFunc := func(ctx context.Context) error {
		transaction := sentry.TransactionFromContext(ctx)

		assert.Equal(t, fmt.Sprintf("GO ROUTINE %s", name), transaction.Name)
		assert.Equal(t, "go.routine", transaction.Op)

		return errors.New("test error")
	}

	go func() {
		goroutine.SentryErrorHandler(
			name,
			testFunc,
		)
	}()
}
