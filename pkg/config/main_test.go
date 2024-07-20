package config_test

import (
	"strconv"
	"testing"

	"github.com/XDoubleU/essentia/pkg/config"
	"github.com/stretchr/testify/assert"
)

const existingKey, nonExistingKey = "key", "non_key"

func TestEnvStr(t *testing.T) {
	expected, def := "string", ""

	t.Setenv(existingKey, expected)

	exists := config.EnvStr(existingKey, def)
	notExists := config.EnvStr(nonExistingKey, def)

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

func TestEnvStrArray(t *testing.T) {
	rawExpected := "string1,string2"
	expected, def := []string{"string1", "string2"}, []string{""}

	t.Setenv(existingKey, rawExpected)

	exists := config.EnvStrArray(existingKey, def)
	notExists := config.EnvStrArray(nonExistingKey, def)

	assert.Equal(t, exists, expected)
	assert.Equal(t, notExists, def)
}

func TestEnvInt(t *testing.T) {
	expected, def := 14, 0

	t.Setenv(existingKey, strconv.Itoa(expected))

	exists := config.EnvInt(existingKey, def)
	notExists := config.EnvInt(nonExistingKey, def)

	assert.Equal(t, exists, expected)
	assert.Equal(t, notExists, def)
}

func TestEnvIntWrong(t *testing.T) {
	expected, def := "string", 0

	t.Setenv(existingKey, expected)

	assert.PanicsWithValue(
		t,
		"can't convert env var 'key' with value 'string' to int",
		func() { config.EnvInt(existingKey, def) },
	)
}

func TestEnvFloat(t *testing.T) {
	expected, def := 14.0, 0.0

	t.Setenv(existingKey, strconv.FormatFloat(expected, 'f', -1, 64))

	exists := config.EnvFloat(existingKey, def)
	notExists := config.EnvFloat(nonExistingKey, def)

	assert.Equal(t, exists, expected)
	assert.Equal(t, notExists, def)
}

func TestEnvFloatWrong(t *testing.T) {
	expected, def := "string", 0.0

	t.Setenv(existingKey, expected)

	assert.PanicsWithValue(
		t,
		"can't convert env var 'key' with value 'string' to float64",
		func() { config.EnvFloat(existingKey, def) },
	)
}

func TestEnvBool(t *testing.T) {
	expected, def := true, false

	t.Setenv(existingKey, strconv.FormatBool(expected))

	exists := config.EnvBool(existingKey, def)
	notExists := config.EnvBool(nonExistingKey, def)

	assert.Equal(t, exists, expected)
	assert.Equal(t, notExists, def)
}

func TestEnvBoolWrong(t *testing.T) {
	expected, def := "string", false

	t.Setenv(existingKey, expected)

	assert.PanicsWithValue(
		t,
		"can't convert env var 'key' with value 'string' to bool",
		func() { config.EnvBool(existingKey, def) },
	)
}
