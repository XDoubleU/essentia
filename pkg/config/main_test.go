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

func TestGetEnvStrArray(t *testing.T) {
	rawExpected := "string1,string2"
	expected, def := []string{"string1", "string2"}, []string{""}

	t.Setenv(existingKey, rawExpected)

	exists := config.GetEnvStrArray(existingKey, def)
	notExists := config.GetEnvStrArray(nonExistingKey, def)

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

func TestGetEnvIntWrong(t *testing.T) {
	expected, def := "string", 0

	t.Setenv(existingKey, expected)

	assert.PanicsWithValue(
		t,
		"can't convert env var 'key' with value 'string' to int",
		func() { config.GetEnvInt(existingKey, def) },
	)
}

func TestGetEnvFloat(t *testing.T) {
	expected, def := 14.0, 0.0

	t.Setenv(existingKey, strconv.FormatFloat(expected, 'f', -1, 64))

	exists := config.GetEnvFloat(existingKey, def)
	notExists := config.GetEnvFloat(nonExistingKey, def)

	assert.Equal(t, exists, expected)
	assert.Equal(t, notExists, def)
}

func TestGetEnvFloatWrong(t *testing.T) {
	expected, def := "string", 0.0

	t.Setenv(existingKey, expected)

	assert.PanicsWithValue(
		t,
		"can't convert env var 'key' with value 'string' to float64",
		func() { config.GetEnvFloat(existingKey, def) },
	)
}

func TestGetEnvBool(t *testing.T) {
	expected, def := true, false

	t.Setenv(existingKey, strconv.FormatBool(expected))

	exists := config.GetEnvBool(existingKey, def)
	notExists := config.GetEnvBool(nonExistingKey, def)

	assert.Equal(t, exists, expected)
	assert.Equal(t, notExists, def)
}

func TestGetEnvBoolWrong(t *testing.T) {
	expected, def := "string", false

	t.Setenv(existingKey, expected)

	assert.PanicsWithValue(
		t,
		"can't convert env var 'key' with value 'string' to bool",
		func() { config.GetEnvBool(existingKey, def) },
	)
}
