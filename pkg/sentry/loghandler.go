package sentry

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/XDoubleU/essentia/pkg/config"
	"github.com/getsentry/sentry-go"
)

// LogHandler is used for capturing logs and sending these to Sentry.
type LogHandler struct {
	level   slog.Level
	handler slog.Handler
	attrs   []slog.Attr
	groups  []string
}

// NewLogHandler returns a new [SentryLogHandler].
func NewLogHandler(env string, handler slog.Handler) slog.Handler {
	level := slog.LevelInfo

	if env == config.DevEnv {
		level = slog.LevelDebug
	}

	return &LogHandler{
		handler: handler,
		attrs:   []slog.Attr{},
		groups:  []string{},
		level:   level,
	}
}

// Enabled checks if logs are enabled in
// a [LogHandler] for a certain [slog.Level].
func (l *LogHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= l.level
}

// WithAttrs adds [[]slog.Attr] to a [SentryLogHandler].
func (l *LogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	//nolint:exhaustruct //other fields are optional
	return &LogHandler{
		attrs:  append(l.attrs, attrs...),
		groups: l.groups,
	}
}

// WithGroup adds a group to a [SentryLogHandler].
func (l *LogHandler) WithGroup(name string) slog.Handler {
	//nolint:exhaustruct //other fields are optional
	return &LogHandler{
		attrs:  l.attrs,
		groups: append(l.groups, name),
	}
}

// Handle handles a [slog.Record] by a [SentryLogHandler].
func (l *LogHandler) Handle(ctx context.Context, record slog.Record) error {
	if record.Level == slog.LevelError {
		sendErrorToSentry(ctx, recordToError(record))
	}

	return l.handler.Handle(ctx, record)
}

func sendErrorToSentry(ctx context.Context, err error) {
	if hub := sentry.GetHubFromContext(ctx); hub != nil {
		hub.WithScope(func(scope *sentry.Scope) {
			scope.SetLevel(sentry.LevelError)
			hub.CaptureException(err)
		})
	}
}

func recordToError(record slog.Record) error {
	err := record.Message

	record.Attrs(func(a slog.Attr) bool {
		err += fmt.Sprintf(" %s=%s", a.Key, a.Value)
		return true
	})

	return errors.New(err)
}
