package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/XDoubleU/essentia/internal/helpers"
)

func TestGetPaged[T any](t *testing.T, ts *httptest.Server, path string, rsData T) {
	TestGeneric(t, ts, http.MethodGet, path, rsData)
}

func TestGetSingle[T any](t *testing.T, ts *httptest.Server, path string, rsData T) {
	TestGeneric(t, ts, http.MethodGet, path, rsData)
}

func TestCreate[T any](t *testing.T, ts *httptest.Server, path string, rsData T) {
	TestGeneric(t, ts, http.MethodPost, path, rsData)
}

func TestUpdate[T any](t *testing.T, ts *httptest.Server, path string, rsData T) {
	TestGeneric(t, ts, http.MethodPatch, path, rsData)
}

func TestDelete[T any](t *testing.T, ts *httptest.Server, path string, rsData T) {
	TestGeneric(t, ts, http.MethodDelete, path, rsData)
}

func TestGeneric[T any](t *testing.T, ts *httptest.Server, method string, path string, rsData T) {
	req, err := http.NewRequest(
		method,
		fmt.Sprintf("%s/%s", ts.URL, path),
		nil,
	)

	if err != nil {
		t.Errorf("error when creating request: %v", err)
		t.FailNow()
		return
	}

	rs, err := ts.Client().Do(req)
	if err != nil {
		t.Errorf("error when making request: %v", err)
		t.FailNow()
		return
	}

	err = helpers.ReadJSON(rs.Body, &rsData, false)
	if err != nil {
		t.Errorf("error when parsing response: %v", err)
		t.FailNow()
		return
	}
}
