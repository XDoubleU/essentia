package query_values

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

func ReadUUIDArray(r *http.Request, name string) ([]string, error) {
	param := r.URL.Query().Get(name)

	values := strings.Split(param, ",")

	var results []string

	for _, value := range values {
		result, err := uuid.Parse(value)
		if err != nil {
			return nil, err
		}
		results = append(results, result.String())
	}

	return results, nil
}

func ReadStr(r *http.Request, name string, defaultValue string) string {
	param := r.URL.Query().Get(name)

	if param == "" {
		return defaultValue
	}

	return param
}

func ReadInt(
	r *http.Request,
	name string,
	defaultValue int64,
) (int64, error) {
	param := r.URL.Query().Get(name)

	if param == "" {
		return defaultValue, nil
	}

	value, err := strconv.ParseInt(param, 10, 64)
	if err != nil || value < 1 {
		return 0, fmt.Errorf("invalid %s query param", name)
	}

	return value, nil
}

func ReadDate(
	r *http.Request,
	name string,
	defaultValue *time.Time,
) (*time.Time, error) {
	param := r.URL.Query().Get(name)

	if param == "" {
		return defaultValue, nil
	}

	// TODO
	value, err := time.Parse("", param)
	if err != nil {
		return nil, fmt.Errorf(
			"invalid %s query param, need format yyyy-MM-dd",
			name,
		)
	}

	return &value, nil
}
