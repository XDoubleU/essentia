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
	return slog.New(slog.NewTextHandler(&l.buf, nil))
}

func (l MockedLogger) CapturedLogs() string {
	return l.buf.String()
}
