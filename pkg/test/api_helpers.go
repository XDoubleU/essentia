package test

import (
	"net/http"
	"strconv"
	"testing"
)

// PaginatedEndpointTester uses a predefined configuration
// for a MatrixTester to test boundaries of a paginated HTTP endpoint.
func PaginatedEndpointTester(
	t *testing.T,
	baseRequest RequestTester,
	pageQueryParamName string,
	maxPage int,
) {
	t.Helper()

	mt := CreateMatrixTester(baseRequest)

	pagesAndStatusCodes := map[int]int{
		-1:          http.StatusBadRequest,
		0:           http.StatusBadRequest,
		1:           http.StatusOK,
		maxPage:     http.StatusOK,
		maxPage + 1: http.StatusOK,
	}

	for page, statusCode := range pagesAndStatusCodes {
		query := map[string]string{
			pageQueryParamName: strconv.Itoa(page),
		}
		mt.AddTestCaseStatusCode(query, statusCode)
	}

	mt.Do(t)
}
