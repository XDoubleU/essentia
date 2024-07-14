package tools

import (
	"errors"
	"strconv"
)

// AnyToString converts any value to a string.
func AnyToString(value any) (string, error) {
	var result string

	var strVal string
	var int64Val int64

	var ok bool

	if strVal, ok = value.(string); ok {
		result = strVal
	} else if int64Val, ok = value.(int64); ok {
		result = strconv.FormatInt(int64Val, 10)
	} else {
		return "", errors.New("undefined type")
	}

	return result, nil
}
