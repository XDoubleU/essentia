// Package parse contains helper functions for parsing
// different kinds of URL and Query parameters.
package parse

import (
	"fmt"
	"net/http"
	"strings"
)

const (
	urlParamType   string = "URL"
	queryParamType string = "query"
)

// URLParam is used to parse a required parameter provided to a URL.
func URLParam[T any](
	r *http.Request,
	paramName string,
	parserFunc ParserFunc[T],
) (T, error) {
	return parseURLParam(r, paramName, true, *new(T), parserFunc)
}

// RequiredQueryParam is used to parse a required parameter provided to a query.
func RequiredQueryParam[T any](
	r *http.Request,
	paramName string,
	parserFunc ParserFunc[T],
) (T, error) {
	return parseQueryParam(r, paramName, true, *new(T), parserFunc)
}

// QueryParam is used to parse an optional parameter provided to a query.
func QueryParam[T any](
	r *http.Request,
	paramName string,
	defaultValue T,
	parserFunc ParserFunc[T],
) (T, error) {
	return parseQueryParam(r, paramName, false, defaultValue, parserFunc)
}

// RequiredArrayQueryParam is used to parse
// a required array parameter provided to a query.
// The format that should be used here is: ?paramName=1,2,3&...
func RequiredArrayQueryParam[T any](
	r *http.Request,
	paramName string,
	parserFunc ParserFunc[T],
) ([]T, error) {
	return parseArrayQueryParam(r, paramName, true, []T{}, parserFunc)
}

// ArrayQueryParam is used to parse an optional array parameter provided to a query.
// The format that should be used here is: ?paramName=1,2,3&...
func ArrayQueryParam[T any](
	r *http.Request,
	paramName string,
	defaultValue []T,
	parserFunc ParserFunc[T],
) ([]T, error) {
	return parseArrayQueryParam(r, paramName, false, defaultValue, parserFunc)
}

func parseURLParam[T any](
	r *http.Request,
	paramName string,
	required bool,
	defaultValue T,
	parserFunc ParserFunc[T],
) (T, error) {
	param := r.PathValue(paramName)
	return parseParam(
		paramName,
		urlParamType,
		param,
		required,
		defaultValue,
		parserFunc,
	)
}

func parseQueryParam[T any](
	r *http.Request,
	paramName string,
	required bool,
	defaultValue T,
	parserFunc ParserFunc[T],
) (T, error) {
	param := r.URL.Query().Get(paramName)
	return parseParam(
		paramName,
		queryParamType,
		param,
		required,
		defaultValue,
		parserFunc,
	)
}

func parseArrayQueryParam[T any](
	r *http.Request,
	paramName string,
	required bool,
	defaultValue []T,
	parserFunc ParserFunc[T],
) ([]T, error) {
	param := r.URL.Query().Get(paramName)
	values := strings.Split(param, ",")

	if param == "" {
		if !required {
			return defaultValue, nil
		}

		return []T{}, fmt.Errorf("missing %s param '%s'", queryParamType, paramName)
	}

	var results []T

	for _, value := range values {
		result, err := parseParam(
			paramName,
			queryParamType,
			value,
			true,
			*new(T),
			parserFunc,
		)
		if err != nil {
			return []T{}, err
		}
		results = append(results, result)
	}

	return results, nil
}

func parseParam[T any](
	paramName string,
	paramType string,
	value string,
	required bool,
	defaultValue T,
	parserFunc ParserFunc[T],
) (T, error) {
	if value == "" {
		if !required {
			return defaultValue, nil
		}

		return *new(T), fmt.Errorf("missing %s param '%s'", paramType, paramName)
	}

	if parserFunc == nil {
		castedValue, ok := any(value).(T)
		if !ok {
			return *new(T), fmt.Errorf("can't cast value to provided type T")
		}

		return castedValue, nil
	}

	result, err := parserFunc(paramType, paramName, value)
	if err != nil {
		return result, err
	}

	return result, nil
}
