package test_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	httptools "github.com/XDoubleU/essentia/pkg/communication/http"
	"github.com/XDoubleU/essentia/pkg/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	tReq.SetBody(reqData)

	rs := tReq.Do(t)

	var rsData map[string]string
	err := httptools.ReadJSON(rs.Body, &rsData)
	require.Nil(t, err)

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
	tReq.SetBody(reqData)
	tReq.SetTestServer(ts)

	rs := tReq.Do(t)

	var rsData map[string]string
	err := httptools.ReadJSON(rs.Body, &rsData)
	require.Nil(t, err)

	assert.Equal(t, http.StatusOK, rs.StatusCode)
	assert.Equal(t, reqData, rsData)
}

func TestRequestTesterNoTestServerOrHandler(t *testing.T) {
	tReq := test.CreateRequestTester(nil, http.MethodGet, "")

	assert.PanicsWithValue(
		t,
		"handler nor test server has been set",
		func() { tReq.Do(t) },
	)
}
