package test_test

import (
	"math"
	"net/http"
	"testing"

	httptools "github.com/XDoubleU/essentia/pkg/communication/http"
	"github.com/XDoubleU/essentia/pkg/errors"
	"github.com/XDoubleU/essentia/pkg/parse"
	"github.com/XDoubleU/essentia/pkg/test"
)

func paginatedEndpointHandler(w http.ResponseWriter, r *http.Request) {
	pageSize := 2
	data := []string{"1", "2", "3"}

	page, err := parse.RequiredQueryParam(r, "page", parse.Int(true, false))
	if err != nil {
		httptools.BadRequestResponse(w, r, errors.NewBadRequestError(err))
		return
	}

	start := int(math.Min(float64((page-1)*pageSize), float64(len(data)-1)))
	end := int(math.Min(float64(page*pageSize), float64(len(data)-1)))

	err = httptools.WriteJSON(w, http.StatusOK, data[start:end], nil)
	if err != nil {
		httptools.ServerErrorResponse(w, r, err)
	}
}

func TestPaginatedEndpointTester(t *testing.T) {
	tReq := test.CreateRequestTester(
		http.HandlerFunc(paginatedEndpointHandler),
		test.JSONContentType,
		http.MethodGet,
		"",
	)
	test.PaginatedEndpointTester(t, tReq, "page", 2)
}
