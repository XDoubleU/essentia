package sentry_test

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"testing"

	"github.com/XDoubleU/essentia/pkg/logging"
	sentrytools "github.com/XDoubleU/essentia/pkg/sentry"
	"github.com/getsentry/sentry-go"
	"github.com/stretchr/testify/assert"
)

func TestSentryErrorHandler(t *testing.T) {
	name := "test"

	logger := logging.NewNopLogger()

	testFunc := func(ctx context.Context, logger *slog.Logger) error {
		transaction := sentry.TransactionFromContext(ctx)

		logger.Debug("started execution")

		assert.Equal(t, fmt.Sprintf("GO ROUTINE %s", name), transaction.Name)
		assert.Equal(t, "go.routine", transaction.Op)

		return errors.New("test error")
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		sentrytools.GoRoutineWrapper(
			context.Background(),
			logger,
			name,
			testFunc,
		)
		wg.Done()
	}()

	wg.Wait()
}
