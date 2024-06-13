package parser

import (
	"strconv"

	"github.com/google/uuid"
)

type ParserFunc[T any] func(value string) (T, error)
type IntType interface {
	int | int64
}

func ParseUUID(value string) (string, error) {
	uuidVal, err := uuid.Parse(value)
	if err != nil {
		return "", err
	}

	return uuidVal.String(), nil
}

func ParseInt(value string) (int, error) {
	return parseInt[int](value, 0)
}

func ParseInt64(value string) (int64, error) {
	return parseInt[int64](value, 64)
}

func parseInt[T IntType](value string, bitSize int) (T, error) {
	result, err := strconv.ParseInt(value, 10, bitSize)
	return T(result), err
}
