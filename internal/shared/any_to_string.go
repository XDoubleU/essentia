package shared

import (
	"errors"
	"fmt"
	"strconv"
)

func arrayToString[T any](array []T) (string, error) {
	output := ""

	for _, value := range array {
		strVal, err := AnyToString(value)
		if err != nil {
			return "", err
		}

		output += fmt.Sprintf("%s,", strVal)
	}

	output = output[:len(output)-1]
	return output, nil
}

// AnyToString converts any value to a string.
func AnyToString(value any) (string, error) {
	switch value := value.(type) {
	case string:
		return value, nil
	case int:
		return strconv.Itoa(value), nil
	case int64:
		return strconv.FormatInt(value, 10), nil
	case []string:
		return arrayToString(value)
	case []int:
		return arrayToString(value)
	case []int64:
		return arrayToString(value)
	default:
		return "", errors.New("undefined type")
	}
}
