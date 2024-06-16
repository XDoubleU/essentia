package config_test

import (
	"strconv"
	"testing"

	"github.com/XDoubleU/essentia/pkg/config"
	"github.com/stretchr/testify/assert"
)

const existingKey, nonExistingKey = "key", "non_key"

func TestGetEnvStr(t *testing.T) {
	expected, def := "string", ""

	t.Setenv(existingKey, expected)

	exists := config.GetEnvStr(existingKey, def)
	notExists := config.GetEnvStr(nonExistingKey, def)

	assert.Equal(t, exists, expected)
	assert.Equal(t, notExists, def)
}

func TestGetEnvInt(t *testing.T) {
	expected, def := 14, 0

	t.Setenv(existingKey, strconv.Itoa(expected))

	exists := config.GetEnvInt(existingKey, def)
	notExists := config.GetEnvInt(nonExistingKey, def)

	assert.Equal(t, exists, expected)
	assert.Equal(t, notExists, def)
}

func TestGetEnvBool(t *testing.T) {
	expected, def := true, false

	t.Setenv(existingKey, strconv.FormatBool(expected))

	exists := config.GetEnvBool(existingKey, def)
	notExists := config.GetEnvBool(nonExistingKey, def)

	assert.Equal(t, exists, expected)
	assert.Equal(t, notExists, def)
}
