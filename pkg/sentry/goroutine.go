package sentry

import (
	"context"
	"fmt"
	"net/http"

	"github.com/getsentry/sentry-go"
)

// GoRoutineErrorHandler makes sure a
// go routine and its errors are captured by Sentry.
func GoRoutineErrorHandler(
	ctx context.Context,
	name string,
	f func(ctx context.Context) error,
) {
	name = fmt.Sprintf("GO ROUTINE %s", name)

	hub := sentry.CurrentHub().Clone()
	ctx = sentry.SetHubOnContext(ctx, hub)

	options := []sentry.SpanOption{
		sentry.WithOpName("go.routine"),
	}

	transaction := sentry.StartTransaction(ctx, name, options...)
	transaction.Status = sentry.HTTPtoSpanStatus(http.StatusOK)
	defer transaction.Finish()

	err := f(transaction.Context())

	if err != nil {
		transaction.Status = sentry.HTTPtoSpanStatus(http.StatusInternalServerError)

		hub.WithScope(func(scope *sentry.Scope) {
			scope.SetLevel(sentry.LevelError)
			hub.CaptureException(err)
		})
	}
}
