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

// GetEnvStr extracts a string environment variable.
func GetEnvStr(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	return value
}

// GetEnvStrArray extracts a string array environment variable. The values are seperated by ','.
func GetEnvStrArray(key string, defaultValue []string) []string {
	strVal := GetEnvStr(key, "")
	if len(strVal) == 0 {
		return defaultValue
	}

	return strings.Split(strVal, ",")
}

// GetEnvInt extracts an integer environment variable.
func GetEnvInt(key string, defaultValue int) int {
	strVal := GetEnvStr(key, "")
	if len(strVal) == 0 {
		return defaultValue
	}

	intVal, err := strconv.Atoi(strVal)
	if err != nil {
		panic(fmt.Sprintf(errorMessage, key, strVal, "int"))
	}

	return intVal
}

// GetEnvFloat extracts a float environment variable.
func GetEnvFloat(key string, defaultValue float64) float64 {
	strVal := GetEnvStr(key, "")
	if len(strVal) == 0 {
		return defaultValue
	}

	floatVal, err := strconv.ParseFloat(strVal, 64)
	if err != nil {
		panic(fmt.Sprintf(errorMessage, key, strVal, "float64"))
	}

	return floatVal
}

// GetEnvBool extracts a boolean environment variable.
func GetEnvBool(key string, defaultValue bool) bool {
	strVal := GetEnvStr(key, "")
	if len(strVal) == 0 {
		return defaultValue
	}

	boolVal, err := strconv.ParseBool(strVal)
	if err != nil {
		panic(fmt.Sprintf(errorMessage, key, strVal, "bool"))
	}

	return boolVal
}
