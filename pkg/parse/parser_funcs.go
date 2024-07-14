package parse

import (
	"fmt"
	"strconv"
	"time"

	"github.com/XDoubleU/essentia/internal/shared"
	"github.com/google/uuid"
)

// ParserFunc is the expected format used for parsing data using any parsing function.
type ParserFunc[T any] func(paramType string, paramName string, value string) (T, error)

// UUID is used to parse a parameter as UUID value.
// Technically this only validates if a string is a UUID.
func UUID(paramType string, paramName string, value string) (string, error) {
	uuidVal, err := uuid.Parse(value)
	if err != nil {
		return "", fmt.Errorf(
			"invalid %s param '%s' with value '%s', should be a UUID",
			paramType,
			paramName,
			value,
		)
	}

	return uuidVal.String(), nil
}

// IntFunc parses a parameter as [int].
func IntFunc(isPositive bool, isZero bool) ParserFunc[int] {
	return func(paramType string, paramName string, value string) (int, error) {
		return parseInt[int](isPositive, isZero, paramType, paramName, value, 0)
	}
}

// Int64Func parses a parameter as [int64].
func Int64Func(isPositive bool, isZero bool) ParserFunc[int64] {
	return func(paramType string, paramName string, value string) (int64, error) {
		//nolint:mnd // no magic number
		return parseInt[int64](isPositive, isZero, paramType, paramName, value, 64)
	}
}

func parseInt[T shared.IntType](
	isPositive bool,
	isZero bool,
	paramType string,
	paramName string,
	value string,
	bitSize int,
) (T, error) {
	result, err := strconv.ParseInt(value, 10, bitSize)

	if err != nil {
		return *new(T), fmt.Errorf(
			"invalid %s param '%s' with value '%s', should be an integer",
			paramType,
			paramName,
			value,
		)
	}

	if isPositive && result < 0 {
		return 0, fmt.Errorf(
			"invalid %s param '%s' with value '%s', can't be less than '0'",
			paramType,
			paramName,
			value,
		)
	}

	if !isZero && result == 0 {
		return 0, fmt.Errorf(
			"invalid %s param '%s' with value '%s', can't be '0'",
			paramType,
			paramName,
			value,
		)
	}

	return T(result), nil
}

// DateFunc parses a parameter as a date.
// The parameter should match the required date layout.
func DateFunc(layout string) ParserFunc[time.Time] {
	return func(paramType string, paramName string, value string) (time.Time, error) {
		result, err := time.Parse(layout, value)
		if err != nil {
			return time.Time{}, fmt.Errorf(
				"invalid %s param '%s' with value '%s', need format '%s'",
				paramType,
				paramName,
				value,
				layout,
			)
		}

		return result, nil
	}
}
