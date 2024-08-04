package validate

import (
	"fmt"
	"time"

	"github.com/xdoubleu/essentia/internal/shared"
)

// ValidatorFunc is the expected format used for validating data using [Check].
type ValidatorFunc[T any] func(value T) (bool, string)

// IsNotEmpty checks that the provided value is not empty.
func IsNotEmpty(value string) (bool, string) {
	return value != "", "must be provided"
}

// IsGreaterThanFunc checks if the provided value2 > value1.
func IsGreaterThanFunc[T shared.IntType](value1 T) ValidatorFunc[T] {
	return func(value2 T) (bool, string) {
		return value2 > value1, fmt.Sprintf("must be greater than %d", value1)
	}
}

// IsGreaterThanOrEqualFunc checks if the provided value2 >= value1.
func IsGreaterThanOrEqualFunc[T shared.IntType](value1 T) ValidatorFunc[T] {
	return func(value2 T) (bool, string) {
		return value2 >= value1,
			fmt.Sprintf(
				"must be greater than or equal to %d",
				value1,
			)
	}
}

// IsLesserThanFunc checks if the provided value2 < value1.
func IsLesserThanFunc[T shared.IntType](value1 T) ValidatorFunc[T] {
	return func(value2 T) (bool, string) {
		return value2 < value1, fmt.Sprintf("must be lesser than %d", value1)
	}
}

// IsLesserThanOrEqualFunc checks if the provided value2 <= value1.
func IsLesserThanOrEqualFunc[T shared.IntType](value1 T) ValidatorFunc[T] {
	return func(value2 T) (bool, string) {
		return value2 <= value1,
			fmt.Sprintf(
				"must be lesser than or equal to %d",
				value1,
			)
	}
}

// IsValidTimeZone checks if the provided value is a valid IANA timezone.
func IsValidTimeZone(value string) (bool, string) {
	_, err := time.LoadLocation(value)
	return err == nil, "must be a valid IANA value"
}
