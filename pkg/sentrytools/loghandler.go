package sentrytools

import (
	"context"
	"errors"
	"log/slog"

	"github.com/getsentry/sentry-go"
)

// SentryLogHandler is used for capturing logs and sending these to Sentry.
type SentryLogHandler struct {
	attrs  []slog.Attr
	groups []string
}

// NewSentryLogHandler returns a new [SentryLogHandler].
func NewSentryLogHandler() slog.Handler {
	return &SentryLogHandler{
		attrs:  []slog.Attr{},
		groups: []string{},
	}
}

// Enabled checks if logs are enabled in
// a [SentryLogHandler] for a certain [slog.Level].
func (l *SentryLogHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= slog.LevelError
}

// WithAttrs adds [[]slog.Attr] to a [SentryLogHandler].
func (l *SentryLogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &SentryLogHandler{
		attrs:  append(l.attrs, attrs...),
		groups: l.groups,
	}
}

// WithGroup adds a group to a [SentryLogHandler].
func (l *SentryLogHandler) WithGroup(name string) slog.Handler {
	return &SentryLogHandler{
		attrs:  l.attrs,
		groups: append(l.groups, name),
	}
}

// Handle handles a [slog.Record] by a [SentryLogHandler].
func (l *SentryLogHandler) Handle(ctx context.Context, record slog.Record) error {
	sendErrorToSentry(ctx, errors.New(record.Message))
	return nil
}

func sendErrorToSentry(ctx context.Context, err error) {
	if hub := sentry.GetHubFromContext(ctx); hub != nil {
		hub.WithScope(func(scope *sentry.Scope) {
			scope.SetLevel(sentry.LevelError)
			hub.CaptureException(err)
		})
	}
}
