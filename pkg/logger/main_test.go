package logger_test

import (
	"testing"

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
