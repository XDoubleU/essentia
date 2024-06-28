package mocks

import (
	"bytes"
	"log"
)

type MockedLogger struct {
	Buffer bytes.Buffer
}

func (l *MockedLogger) GetLogger() *log.Logger {
	return log.New(&l.Buffer, "", log.LstdFlags)
}
