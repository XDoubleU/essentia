package test_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/XDoubleU/essentia/pkg/httptools"
	"github.com/XDoubleU/essentia/pkg/parse"
	"github.com/XDoubleU/essentia/pkg/test"
)

func matrixTestHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("cookie")
	if err != nil && !errors.Is(err, http.ErrNoCookie) {
		httptools.ServerErrorResponse(w, r, err)
		return
	}

	if cookie != nil && cookie.Value == "value" {
		httptools.UnauthorizedResponse(w, r, "unauthorized")
		return
	}

	param, err := parse.QueryParam(r, "param", "default", nil)
	if err != nil {
		httptools.ServerErrorResponse(w, r, err)
		return
	}

	if param == "value" {
		httptools.ForbiddenResponse(w, r)
		return
	}

	var data map[string]string
	err = httptools.ReadJSON(r.Body, &data)
	if err != nil {
		httptools.ServerErrorResponse(w, r, err)
		return
	}

	if data["test"] == "error" {
		httptools.ErrorResponse(
			w,
			r,
			http.StatusBadRequest,
			map[string]any{"message": "test"},
		)
		return
	}

	err = httptools.WriteJSON(w, http.StatusOK, nil, nil)
	if err != nil {
		httptools.ServerErrorResponse(w, r, err)
	}
}

func TestMatrixTester(t *testing.T) {
	tReq := test.CreateRequestTester(
		http.HandlerFunc(matrixTestHandler),
		http.MethodGet,
		"",
	)
	mt := test.CreateMatrixTester(tReq)

	mt.AddTestCaseErrorMessage(map[string]string{
		"test": "error",
	}, map[string]any{"message": "test"})

	mt.AddTestCaseStatusCode(map[string]string{
		"param": "value",
	}, http.StatusForbidden)

	mt.AddTestCaseCookieStatusCode(
		&http.Cookie{Name: "cookie", Value: "value"},
		http.StatusUnauthorized,
	)

	mt.Do(t)
}
