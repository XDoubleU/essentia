package postgres_test

import (
	"testing"
	"time"

	"github.com/XDoubleU/essentia/internal/mocks"
	"github.com/XDoubleU/essentia/pkg/database/postgres"
	"github.com/stretchr/testify/assert"
)

func TestConnectRetries(t *testing.T) {
	mockedLogger := mocks.NewMockedLogger()

	_, err := postgres.Connect(
		mockedLogger.GetLogger(),
		"",
		1,
		"1s",
		"1",
		200*time.Millisecond,
		time.Second,
	)

	assert.NotNil(t, err)
	assert.Contains(t, mockedLogger.GetCapturedLogs(), "retrying in")
}
