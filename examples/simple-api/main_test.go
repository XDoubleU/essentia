package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/XDoubleU/essentia/pkg/test"
)

func NewTestDataRepository() *DataRepository {
	dataRepo := NewDataRepository()

	dataRepo.data["first"] = Data{
		name: "First",
	}

	return &dataRepo
}

func TestGeneric(t *testing.T) {
	ts := httptest.NewTLSServer(setupRouter(NewTestDataRepository()))
	defer ts.Close()

	var rsData map[string]string
	test.TestGeneric(t, ts, http.MethodGet, "generic", &rsData)

	test.Equal(t, rsData["message"], "ok")
}

func TestGetSingle(t *testing.T) {
	ts := httptest.NewTLSServer(setupRouter(NewTestDataRepository()))
	defer ts.Close()

	var rsData Data
	test.TestGetSingle(t, ts, "single/first", &rsData)

	test.Equal(t, rsData.name, "First")
}

/*
func TestGetPaged(t *testing.T) {
	ts := httptest.NewTLSServer(App())

}


func TestCreate() {
	ts := httptest.NewTLSServer(App())
}

func TestUpdate() {
	ts := httptest.NewTLSServer(App())
}

func TestDelete() {
	ts := httptest.NewTLSServer(App())
}
*/
