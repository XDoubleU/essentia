package parser

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const (
	URLParam   string = "URL"
	QueryParam        = "query"
)

func ParseURLParam[T any](r *http.Request, name string, parserFunc ParserFunc[T]) (T, error) {
	return parseURLParam(r, name, parserFunc, true)
}

func ParseURLParamWithError[T any](r *http.Request, name string, parserFunc ParserFunc[T]) (T, error) {
	return parseURLParam(r, name, parserFunc, false)
}

func parseURLParam[T any](r *http.Request, name string, parserFunc ParserFunc[T], obfuscateError bool) (T, error) {
	params := httprouter.ParamsFromContext(r.Context())
	param := params.ByName(name)
	return parseParam(name, URLParam, param, parserFunc, obfuscateError)
}

func ParseQueryParam[T any](r *http.Request, name string, parserFunc ParserFunc[T]) (T, error) {
	return parseQueryParam(r, name, parserFunc, true)
}

func ParseQueryParamWithError[T any](r *http.Request, name string, parserFunc ParserFunc[T]) (T, error) {
	return parseQueryParam(r, name, parserFunc, false)
}

func parseQueryParam[T any](r *http.Request, name string, parserFunc ParserFunc[T], obfuscateError bool) (T, error) {
	param := r.URL.Query().Get(name)
	return parseParam(name, QueryParam, param, parserFunc, obfuscateError)
}

func parseParam[T any](name string, paramType string, param string, parserFunc ParserFunc[T], obfuscateError bool) (T, error) {
	result, err := parserFunc(param)
	if !obfuscateError {
		return result, err
	}

	return result, fmt.Errorf("invalid %s param '%s'", paramType, name)
}
