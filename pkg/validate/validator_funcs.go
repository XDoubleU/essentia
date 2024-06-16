package validate

import (
	"fmt"
	"time"

	"github.com/XDoubleU/essentia/internal/shared"
)

type ValidatorFunc[T any] func(value T) (bool, string)

func IsNotEmpty(value string) (bool, string) {
	return value != "", "must be provided"
}

func IsGreaterThanFunc[T shared.IntType](value1 T) ValidatorFunc[T] {
	return func(value2 T) (bool, string) {
		return value2 > value1, fmt.Sprintf("must be greater than %d", value1)
	}
}

func IsGreaterThanOrEqualFunc[T shared.IntType](value1 T) ValidatorFunc[T] {
	return func(value2 T) (bool, string) {
		return value2 >= value1,
			fmt.Sprintf(
				"must be greater than or equal to %d",
				value1,
			)
	}
}

func IsLesserThanFunc[T shared.IntType](value1 T) ValidatorFunc[T] {
	return func(value2 T) (bool, string) {
		return value2 < value1, fmt.Sprintf("must be lesser than %d", value1)
	}
}

func IsLesserThanOrEqualFunc[T shared.IntType](value1 T) ValidatorFunc[T] {
	return func(value2 T) (bool, string) {
		return value2 <= value1,
			fmt.Sprintf(
				"must be lesser than or equal to %d",
				value1,
			)
	}
}

func IsValidTimeZone(value string) (bool, string) {
	_, err := time.LoadLocation(value)
	return err == nil, "must be a valid IANA value"
}
