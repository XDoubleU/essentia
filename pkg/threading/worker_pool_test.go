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

func doWork(_ context.Context, _ *slog.Logger) {
	time.Sleep(300 * time.Millisecond)
}

func TestBasicWorkerPool(t *testing.T) {
	workerpool := threading.NewWorkerPool(logging.NewNopLogger(), 1, 2)

	workerpool.EnqueueWork(doWork)
	workerpool.EnqueueWork(doWork)

	workerpool.WaitUntilDone()
	assert.True(t, workerpool.Active())

	workerpool.Stop()
	assert.False(t, workerpool.Active())
}
