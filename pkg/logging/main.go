// Package logging contains helpers for logging.
package logging

import (
	"bytes"
	"io"
	"log/slog"
)

// ErrAttr provides a [slog.Attr] to reduce a bit of boilerplate when logging errors.
// Credits go to https://github.com/golang/go/issues/59364#issuecomment-1493237877.
func ErrAttr(err error) slog.Attr {
	return slog.Any("error", err)
}

// NewNopLogger provides a NopLogger which uses [io.Discard] to write logs to.
func NewNopLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

// NewBufLogHandler provides a [slog.TextHandler]
// which logs to the provided [bytes.Buffer].
func NewBufLogHandler(buf *bytes.Buffer, opts *slog.HandlerOptions) *slog.TextHandler {
	return slog.NewTextHandler(buf, opts)
}
