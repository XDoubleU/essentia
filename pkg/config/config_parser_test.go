package config

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

var existingKey, nonExistingKey = "key", "non_key"

func TestGetEnvStr(t *testing.T) {
	expected, def := "string", ""

	t.Setenv(existingKey, expected)

	exists := GetEnvStr(existingKey, def)
	notExists := GetEnvStr(nonExistingKey, def)

	assert.Equal(t, exists, expected)
	assert.Equal(t, notExists, def)
}

func TestGetEnvInt(t *testing.T) {
	expected, def := 14, 0

	t.Setenv(existingKey, strconv.Itoa(expected))

	exists := GetEnvInt(existingKey, def)
	not_exists := GetEnvInt(nonExistingKey, def)

	assert.Equal(t, exists, expected)
	assert.Equal(t, not_exists, def)
}

func TestGetEnvBool(t *testing.T) {
	expected, def := true, false

	t.Setenv(existingKey, strconv.FormatBool(expected))

	exists := GetEnvBool(existingKey, def)
	not_exists := GetEnvBool(nonExistingKey, def)

	assert.Equal(t, exists, expected)
	assert.Equal(t, not_exists, def)
}
