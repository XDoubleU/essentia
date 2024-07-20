// Package config provides functions which can be used to
// extract environment variables and parse them to the right type.
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const errorMessage = "can't convert env var '%s' with value '%s' to %s"

// EnvStr extracts a string environment variable.
func EnvStr(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	return value
}

// GetEnvStrArray extracts a string
// array environment variable. The values should be seperated by ','.
func GetEnvStrArray(key string, defaultValue []string) []string {
	strVal := GetEnvStr(key, "")
	if len(strVal) == 0 {
		return defaultValue
	}

	return strings.Split(strVal, ",")
}

// EnvStrArray extracts a string
// array environment variable. The values should be seperated by ','.
func EnvStrArray(key string, defaultValue []string) []string {
	strVal := EnvStr(key, "")
	if len(strVal) == 0 {
		return defaultValue
	}

	return strings.Split(strVal, ",")
}

// EnvInt extracts an integer environment variable.
func EnvInt(key string, defaultValue int) int {
	strVal := EnvStr(key, "")
	if len(strVal) == 0 {
		return defaultValue
	}

	intVal, err := strconv.Atoi(strVal)
	if err != nil {
		panic(fmt.Sprintf(errorMessage, key, strVal, "int"))
	}

	return intVal
}

// EnvFloat extracts a float environment variable.
func EnvFloat(key string, defaultValue float64) float64 {
	strVal := EnvStr(key, "")
	if len(strVal) == 0 {
		return defaultValue
	}

	floatVal, err := strconv.ParseFloat(strVal, 64)
	if err != nil {
		panic(fmt.Sprintf(errorMessage, key, strVal, "float64"))
	}

	return floatVal
}

// EnvBool extracts a boolean environment variable.
func EnvBool(key string, defaultValue bool) bool {
	strVal := EnvStr(key, "")
	if len(strVal) == 0 {
		return defaultValue
	}

	boolVal, err := strconv.ParseBool(strVal)
	if err != nil {
		panic(fmt.Sprintf(errorMessage, key, strVal, "bool"))
	}

	return boolVal
}
