package sentrytools_test

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/getsentry/sentry-go"
	"github.com/stretchr/testify/assert"
	"github.com/xdoubleu/essentia/pkg/sentrytools"
)

func TestSentryErrorHandler(t *testing.T) {
	name := "test"

	testFunc := func(ctx context.Context) error {
		transaction := sentry.TransactionFromContext(ctx)

		assert.Equal(t, fmt.Sprintf("GO ROUTINE %s", name), transaction.Name)
		assert.Equal(t, "go.routine", transaction.Op)

		return errors.New("test error")
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		sentrytools.GoRoutineErrorHandler(
			name,
			testFunc,
		)
		wg.Done()
	}()

	wg.Wait()
}
