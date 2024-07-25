package sentrytools

import (
	"context"
	"errors"
	"log/slog"

	"github.com/getsentry/sentry-go"
)

type SentryLogHandler struct {
	attrs  []slog.Attr
	groups []string
}

func NewSentryLogHandler() slog.Handler {
	return &SentryLogHandler{
		attrs:  []slog.Attr{},
		groups: []string{},
	}
}

func (l *SentryLogHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= slog.LevelError
}

func (l *SentryLogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &SentryLogHandler{
		attrs: append(l.attrs, attrs...),
	}
}

func (l *SentryLogHandler) WithGroup(name string) slog.Handler {
	return &SentryLogHandler{
		attrs:  l.attrs,
		groups: append(l.groups, name),
	}
}

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
