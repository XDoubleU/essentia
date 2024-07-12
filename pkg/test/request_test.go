package test_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/XDoubleU/essentia/pkg/httptools"
	"github.com/XDoubleU/essentia/pkg/test"
	"github.com/stretchr/testify/assert"
)

func testHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("cookiename")
	if err != nil || cookie.Value != "value" {
		httptools.ServerErrorResponse(w, r, err)
		return
	}

	var data map[string]string
	err = httptools.ReadJSON(r.Body, &data)
	if err != nil {
		httptools.ServerErrorResponse(w, r, err)
		return
	}

	err = httptools.WriteJSON(w, http.StatusOK, data, nil)
	if err != nil {
		httptools.ServerErrorResponse(w, r, err)
	}
}

func TestRequestTester(t *testing.T) {
	reqData := map[string]string{
		"test": "data",
	}

	tReq := test.CreateRequestTester(
		http.HandlerFunc(testHandler),
		http.MethodGet,
		"/test/%d",
		1,
	)
	tReq.AddCookie(&http.Cookie{Name: "cookiename", Value: "value"})
	tReq.SetReqData(reqData)

	var rsData map[string]string
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusOK, rs.StatusCode)
	assert.Equal(t, reqData, rsData)
}

func TestRequestTesterTestServer(t *testing.T) {
	reqData := map[string]string{
		"test": "data",
	}

	ts := httptest.NewServer(http.HandlerFunc(testHandler))

	tReq := test.CreateRequestTester(nil, http.MethodGet, "/test/%d", 1)
	tReq.AddCookie(&http.Cookie{Name: "cookiename", Value: "value"})
	tReq.SetReqData(reqData)
	tReq.SetTestServer(ts)

	var rsData map[string]string
	rs := tReq.Do(t, &rsData)

	assert.Equal(t, http.StatusOK, rs.StatusCode)
	assert.Equal(t, reqData, rsData)
}

func TestRequestTesterNoTestServerOrHandler(t *testing.T) {
	tReq := test.CreateRequestTester(nil, http.MethodGet, "")

	assert.PanicsWithValue(
		t,
		"handler nor test server has been set",
		func() { tReq.Do(t, nil) },
	)
}
