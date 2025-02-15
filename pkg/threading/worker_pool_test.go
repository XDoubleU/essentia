package threading_test

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/XDoubleU/essentia/pkg/logging"
	"github.com/XDoubleU/essentia/pkg/threading"
	"github.com/stretchr/testify/assert"
)

func doWork(_ context.Context, _ *slog.Logger) error {
	time.Sleep(1 * time.Second)
	return nil
}

func TestBasicWorkerPool(t *testing.T) {
	workerpool := threading.NewWorkerPool(logging.NewNopLogger(), 1, 2)

	workerpool.EnqueueWork(doWork)
	workerpool.EnqueueWork(doWork)

	workerpool.WaitUntilDone()
	assert.False(t, workerpool.IsDoingWork())
}
