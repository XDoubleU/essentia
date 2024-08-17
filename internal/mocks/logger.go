package mocks

import (
	"bytes"
	"log/slog"
)

type MockedLogger struct {
	buf bytes.Buffer
}

func NewMockedLogger() MockedLogger {
	return MockedLogger{
		buf: *bytes.NewBuffer([]byte{}),
	}
}

func (l *MockedLogger) Logger() *slog.Logger {
	//nolint:exhaustruct //other fields are optional
	return slog.New(slog.NewTextHandler(&l.buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
}

func (l MockedLogger) CapturedLogs() string {
	return l.buf.String()
}
