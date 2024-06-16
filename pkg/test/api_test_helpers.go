package test

import (
	"net/http"
	"strconv"
	"testing"
)

func TestPaginatedEndpoint(t *testing.T, baseRequest TestRequest, pageQueryParamName string, maxPage int) {
	t.Helper()

	mt := CreateMatrixTester(t, baseRequest)

	pagesAndStatusCodes := map[int]int{
		-1:          http.StatusBadRequest,
		0:           http.StatusBadRequest,
		1:           http.StatusOK,
		maxPage:     http.StatusOK,
		maxPage + 1: http.StatusOK,
	}

	for page, statusCode := range pagesAndStatusCodes {
		query := map[string]string{
			"page": strconv.Itoa(page),
		}
		mt.AddTestCaseStatusCode(query, statusCode)
	}

	mt.Do(t)
}
