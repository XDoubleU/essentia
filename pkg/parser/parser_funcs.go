package parser

import (
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type ParserFunc[T any] func(paramType string, paramName string, value string) (T, error)
type IntType interface {
	int | int64
}

func ParseUUID(paramType string, paramName string, value string) (string, error) {
	uuidVal, err := uuid.Parse(value)
	if err != nil {
		return "", err
	}

	return uuidVal.String(), nil
}

func ParseIntFunc(isPositive bool, isZero bool) ParserFunc[int] {
	return func(paramType string, paramName string, value string) (int, error) {
		return parseInt[int](isPositive, isZero, paramType, paramName, value, 0)
	}
}

func ParseInt64Func(isPositive bool, isZero bool) ParserFunc[int64] {
	return func(paramType string, paramName string, value string) (int64, error) {
		return parseInt[int64](isPositive, isZero, paramType, paramName, value, 64)
	}
}

func parseInt[T IntType](isPositive bool, isZero bool, paramType string, paramName string, value string, bitSize int) (T, error) {
	result, err := strconv.ParseInt(value, 10, bitSize)

	if err != nil {
		return *new(T), fmt.Errorf("invalid %s param '%s' with value '%s', should be an integer", paramType, paramName, value)
	}

	if isPositive && result < 0 {
		return 0, fmt.Errorf("invalid %s param '%s' with value '%s', can't be less than '0'", paramType, paramName, value)
	}

	if !isZero && result == 0 {
		return 0, fmt.Errorf("invalid %s param '%s' with value '%s', can't be '0'", paramType, paramName, value)
	}

	return T(result), nil
}

func ParseDateFunc(layout string) ParserFunc[time.Time] {
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
