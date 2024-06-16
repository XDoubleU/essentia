package parser

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

const (
	URLParam   string = "URL"
	QueryParam string = "query"
)

func ParseURLParam[T any](r *http.Request, paramName string, parserFunc ParserFunc[T]) (T, error) {
	return parseURLParam(r, paramName, false, *new(T), parserFunc)
}

func ParseRequiredQueryParam[T any](r *http.Request, paramName string, parserFunc ParserFunc[T]) (T, error) {
	return parseQueryParam(r, paramName, true, *new(T), parserFunc)
}

func ParseQueryParam[T any](r *http.Request, paramName string, defaultValue T, parserFunc ParserFunc[T]) (T, error) {
	return parseQueryParam(r, paramName, false, defaultValue, parserFunc)
}

func ParseRequiredArrayQueryParam[T any](r *http.Request, paramName string, parserFunc ParserFunc[T]) ([]T, error) {
	return parseArrayQueryParam(r, paramName, true, []T{}, parserFunc)
}

func ParseArrayQueryParam[T any](r *http.Request, paramName string, defaultValue []T, parserFunc ParserFunc[T]) ([]T, error) {
	return parseArrayQueryParam(r, paramName, false, defaultValue, parserFunc)
}

func parseURLParam[T any](r *http.Request, paramName string, required bool, defaultValue T, parserFunc ParserFunc[T]) (T, error) {
	params := httprouter.ParamsFromContext(r.Context())
	param := params.ByName(paramName)
	return parseParam(paramName, URLParam, param, required, defaultValue, parserFunc)
}

func parseQueryParam[T any](r *http.Request, paramName string, required bool, defaultValue T, parserFunc ParserFunc[T]) (T, error) {
	param := r.URL.Query().Get(paramName)
	return parseParam(paramName, QueryParam, param, required, defaultValue, parserFunc)
}

func parseArrayQueryParam[T any](r *http.Request, paramName string, required bool, defaultValue []T, parserFunc ParserFunc[T]) ([]T, error) {
	param := r.URL.Query().Get(paramName)
	values := strings.Split(param, ",")

	if len(values) == 0 {
		if !required {
			return defaultValue, nil
		}

		return []T{}, fmt.Errorf("missing %s param '%s'", QueryParam, paramName)
	}

	var results []T

	for _, value := range values {
		result, err := parseParam(paramName, QueryParam, value, true, *new(T), parserFunc)
		if err != nil {
			return []T{}, err
		}
		results = append(results, result)
	}

	return results, nil
}

func parseParam[T any](paramName string, paramType string, value string, required bool, defaultValue T, parserFunc ParserFunc[T]) (T, error) {
	if value == "" {
		if !required {
			return defaultValue, nil
		}

		return *new(T), fmt.Errorf("missing %s param '%s'", paramType, paramName)
	}

	if parserFunc == nil {
		return any(value).(T), nil
	}

	result, err := parserFunc(paramType, paramName, value)
	if err != nil {
		return result, err
	}

	return result, nil
}
