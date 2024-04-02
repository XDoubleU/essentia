package input

import "time"

type Validation func(value any) bool

func RequiredString() Validation {
	return func(value any) bool {
		return value.(string) != ""
	}
}

func MinimumIntegerValue(minimumValue int) Validation {
	return func(value any) bool {
		return value.(int) > minimumValue
	}
}

func ValidTimeZone() Validation {
	return func(value any) bool {
		stringValue := value.(string)
		_, err := time.LoadLocation(stringValue)
		return stringValue != "" && err == nil
	}
}
