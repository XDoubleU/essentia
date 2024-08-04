package sentry

import (
	"context"
	"errors"
	"log/slog"

	"github.com/getsentry/sentry-go"
)

// LogHandler is used for capturing logs and sending these to Sentry.
type LogHandler struct {
	attrs  []slog.Attr
	groups []string
}

// NewLogHandler returns a new [SentryLogHandler].
func NewLogHandler() slog.Handler {
	return &LogHandler{
		attrs:  []slog.Attr{},
		groups: []string{},
	}
}

// Enabled checks if logs are enabled in
// a [LogHandler] for a certain [slog.Level].
func (l *LogHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= slog.LevelError
}

// WithAttrs adds [[]slog.Attr] to a [SentryLogHandler].
func (l *LogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &LogHandler{
		attrs:  append(l.attrs, attrs...),
		groups: l.groups,
	}
}

// WithGroup adds a group to a [SentryLogHandler].
func (l *LogHandler) WithGroup(name string) slog.Handler {
	return &LogHandler{
		attrs:  l.attrs,
		groups: append(l.groups, name),
	}
}

// Handle handles a [slog.Record] by a [SentryLogHandler].
func (l *LogHandler) Handle(ctx context.Context, record slog.Record) error {
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
