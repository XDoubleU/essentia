package logger_test

import (
	"testing"

	"github.com/XDoubleU/essentia/internal/mocks"
	"github.com/XDoubleU/essentia/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func TestGetLoggerSingleton(t *testing.T) {
	logger1 := logger.GetLogger()
	logger2 := logger.GetLogger()

	assert.Equal(t, logger1, logger2)
}

func TestSetGetLogger(t *testing.T) {
	logger.SetLogger(logger.NullLogger)

	assert.Equal(t, logger.NullLogger, logger.GetLogger())
}

func TestOutput(t *testing.T) {
	mockedLogger := mocks.MockedLogger{}
	logger.SetLogger(mockedLogger.GetLogger())

	logger.GetLogger().Print("test")

	assert.Contains(t, mockedLogger.Buffer.String(), "test")
}
