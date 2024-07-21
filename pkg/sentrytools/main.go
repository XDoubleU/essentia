// Package sentrytools contains all sorts of tools for using Sentry.
package sentrytools

import (
	"context"

	"github.com/getsentry/sentry-go"
)

// SendErrorToSentry can be used to send any error to Sentry.
func SendErrorToSentry(ctx context.Context, err error) {
	if hub := sentry.GetHubFromContext(ctx); hub != nil {
		hub.WithScope(func(scope *sentry.Scope) {
			scope.SetLevel(sentry.LevelError)
			hub.CaptureException(err)
		})
	}
}
