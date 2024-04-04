package parser

import "time"

func RequiredString() Validation[string] {
	return func(value string) bool {
		return value != ""
	}
}

func MinimumIntegerValue(minimumValue int) Validation[int] {
	return func(value int) bool {
		return value > minimumValue
	}
}

func ValidTimeZone() Validation[any] {
	return func(value any) bool {
		stringValue := value.(string)
		_, err := time.LoadLocation(stringValue)
		return stringValue != "" && err == nil
	}
}
