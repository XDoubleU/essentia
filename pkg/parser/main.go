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

func ParseURLParam[T any](r *http.Request, paramName string, parserFunc ParserFunc[T]) (value T, prettyErr error, originalErr error) {
	return parseURLParam(r, paramName, false, *new(T), parserFunc)
}

func ParseRequiredQueryParam[T any](r *http.Request, paramName string, parserFunc ParserFunc[T]) (value T, prettyErr error, originalErr error) {
	return parseQueryParam(r, paramName, true, *new(T), parserFunc)
}

func ParseQueryParam[T any](r *http.Request, paramName string, defaultValue T, parserFunc ParserFunc[T]) (value T, prettyErr error, originalErr error) {
	return parseQueryParam(r, paramName, false, defaultValue, parserFunc)
}

func ParseRequiredArrayQueryParam[T any](r *http.Request, paramName string, parserFunc ParserFunc[T]) (value []T, prettyErr error, originalErr error) {
	return parseArrayQueryParam[T](r, paramName, true, []T{}, parserFunc)
}

func ParseArrayQueryParam[T any](r *http.Request, paramName string, defaultValue []T, parserFunc ParserFunc[T]) (value []T, prettyErr error, originalErr error) {
	return parseArrayQueryParam[T](r, paramName, false, defaultValue, parserFunc)
}

func parseURLParam[T any](r *http.Request, paramName string, required bool, defaultValue T, parserFunc ParserFunc[T]) (value T, prettyErr error, originalErr error) {
	params := httprouter.ParamsFromContext(r.Context())
	param := params.ByName(paramName)
	return parseParam(paramName, URLParam, param, required, defaultValue, parserFunc)
}

func parseQueryParam[T any](r *http.Request, paramName string, required bool, defaultValue T, parserFunc ParserFunc[T]) (value T, prettyErr error, originalErr error) {
	param := r.URL.Query().Get(paramName)
	return parseParam(paramName, QueryParam, param, required, defaultValue, parserFunc)
}

func parseArrayQueryParam[T any](r *http.Request, paramName string, required bool, defaultValue []T, parserFunc ParserFunc[T]) (value []T, prettyErr error, originalErr error) {
	param := r.URL.Query().Get(paramName)
	values := strings.Split(param, ",")

	var results []T

	for i, value := range values {
		result, prettyErr, originalErr := parseParam(paramName, QueryParam, value, required, defaultValue[i], parserFunc)
		if prettyErr != nil {
			return []T{}, prettyErr, originalErr
		}
		results = append(results, result)
	}

	return results, nil, nil
}

func parseParam[T any](paramName string, paramType string, value string, required bool, defaultValue T, parserFunc ParserFunc[T]) (result T, prettyErr error, originalErr error) {
	missingRequiredErr := fmt.Errorf("missing %s param '%s'", paramType, paramName)
	invalidErr := fmt.Errorf("invalid %s param '%s' with value '%s'", paramType, paramName, value)

	if parserFunc == nil {
		return any(value).(T), nil, nil
	}

	if value == "" {
		if !required {
			return defaultValue, nil, nil
		}

		return *new(T), missingRequiredErr, missingRequiredErr
	}

	result, err := parserFunc(paramType, paramName, value)
	return result, invalidErr, err
}
