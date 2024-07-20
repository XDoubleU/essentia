package tools

import (
	"errors"
	"strconv"
)

// AnyToString converts any value to a string.
func AnyToString(value any) (string, error) {
	switch value := value.(type) {
	case string:
		return value, nil
	case int:
		return strconv.Itoa(value), nil
	case int64:
		return strconv.FormatInt(value, 10), nil
	default:
		return "", errors.New("undefined type")
	}
}
