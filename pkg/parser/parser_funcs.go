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

func ParseInt(paramType string, paramName string, value string) (int, error) {
	return parseInt[int](value, 0)
}

func ParseInt64(paramType string, paramName string, value string) (int64, error) {
	return parseInt[int64](value, 64)
}

func parseInt[T IntType](value string, bitSize int) (T, error) {
	result, err := strconv.ParseInt(value, 10, bitSize)
	return T(result), err
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
