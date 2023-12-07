package config

import (
	"strconv"
	"testing"

	"github.com/XDoubleU/essentia/pkg/test"
)

var existing_key, non_existing_key = "key", "non_key"

func TestGetEnvStr(t *testing.T) {
	expected, def := "string", ""

	t.Setenv(existing_key, expected)

	exists := GetEnvStr(existing_key, def)
	not_exists := GetEnvStr(non_existing_key, def)

	test.Equal(t, exists, expected)
	test.Equal(t, not_exists, def)
}

func TestGetEnvInt(t *testing.T) {
	expected, def := 14, 0

	t.Setenv(existing_key, strconv.Itoa(expected))

	exists := GetEnvInt(existing_key, def)
	not_exists := GetEnvInt(non_existing_key, def)

	test.Equal(t, exists, expected)
	test.Equal(t, not_exists, def)
}

func TestGetEnvBool(t *testing.T) {
	expected, def := true, false

	t.Setenv(existing_key, strconv.FormatBool(expected))

	exists := GetEnvBool(existing_key, def)
	not_exists := GetEnvBool(non_existing_key, def)

	test.Equal(t, exists, expected)
	test.Equal(t, not_exists, def)
}
