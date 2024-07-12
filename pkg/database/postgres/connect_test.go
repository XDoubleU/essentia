package postgres_test

import (
	"testing"
	"time"

	"github.com/XDoubleU/essentia/pkg/database/postgres"
	"github.com/XDoubleU/essentia/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func TestConnectRetries(t *testing.T) {
	logger.SetLogger(logger.NullLogger)

	_, err := postgres.Connect("", 1, "1s", "1", 200*time.Millisecond, time.Second)
	assert.NotNil(t, err)
}
