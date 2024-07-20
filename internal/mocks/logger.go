package mocks

import (
	"bytes"
	"log"
)

type MockedLogger struct {
	buf bytes.Buffer
}

func NewMockedLogger() MockedLogger {
	return MockedLogger{
		buf: *bytes.NewBuffer([]byte{}),
	}
}

func (l *MockedLogger) Logger() *log.Logger {
	return log.New(&l.buf, "", log.LstdFlags)
}

func (l MockedLogger) CapturedLogs() string {
	return l.buf.String()
}
