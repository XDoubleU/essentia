package test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	httptools "github.com/xdoubleu/essentia/pkg/communication/http"
)

// A RequestTester is used to test a certain HTTP request.
type RequestTester struct {
	handler http.Handler
	ts      *httptest.Server
	method  string
	path    string
	body    any
	query   map[string]string
	cookies []*http.Cookie
}

// CreateRequestTester creates a new [RequestTester].
func CreateRequestTester(
	handler http.Handler,
	method, path string,
	pathValues ...any,
) RequestTester {
	if len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}

	if len(pathValues) > 0 {
		path = fmt.Sprintf(path, pathValues...)
	}

	return RequestTester{
		handler,
		nil,
		method,
		path,
		nil,
		make(map[string]string),
		[]*http.Cookie{},
	}
}

// SetTestServer sets the test server of a [RequestTester].
// This allows to reuse existing test servers.
func (tReq *RequestTester) SetTestServer(ts *httptest.Server) {
	tReq.ts = ts
}

// SetBody sets the request body of a [RequestTester].
func (tReq *RequestTester) SetBody(body any) {
	tReq.body = body
}

// SetQuery sets the query of a [RequestTester].
func (tReq *RequestTester) SetQuery(query map[string]string) {
	tReq.query = query
}

// AddCookie adds a cookie to a [RequestTester].
// Can be used multiple times for adding several cookies.
func (tReq *RequestTester) AddCookie(cookie *http.Cookie) {
	tReq.cookies = append(tReq.cookies, cookie)
}

// Do executes a [RequestTester] returning the response of the request
// and providing the returned data to rsData.
func (tReq RequestTester) Do(t *testing.T, rsData any) *http.Response {
	t.Helper()

	var body []byte
	var err error

	if tReq.ts == nil && tReq.handler != nil {
		tReq.ts = httptest.NewServer(tReq.handler)
		defer tReq.ts.Close()
	}

	if tReq.ts == nil {
		panic("handler nor test server has been set")
	}

	if tReq.body != nil {
		body, err = json.Marshal(tReq.body)
		if err != nil {
			t.Errorf("error when marshalling body: %v", err)
			t.FailNow()
			return nil
		}
	}

	req, err := http.NewRequestWithContext(
		context.Background(),
		tReq.method,
		fmt.Sprintf("%s/%s", tReq.ts.URL, tReq.path),
		bytes.NewReader(body),
	)

	if err != nil {
		t.Errorf("error when creating request: %v", err)
		t.FailNow()
		return nil
	}

	if len(tReq.query) > 0 {
		query := req.URL.Query()

		for key, value := range tReq.query {
			query.Add(key, value)
		}

		req.URL.RawQuery = query.Encode()
	}

	for _, cookie := range tReq.cookies {
		req.AddCookie(cookie)
	}

	rs, err := tReq.ts.Client().Do(req)
	if err != nil {
		t.Errorf("error when making request: %v", err)
		t.FailNow()
		return nil
	}

	if rsData != nil {
		err = httptools.ReadJSON(rs.Body, &rsData)
		if err != nil {
			t.Errorf("error when parsing response: %v", err)
			t.FailNow()
			return nil
		}
	}

	return rs
}

// Copy creates a copy of a [RequestTester] in order to easily test similar requests.
func (tReq RequestTester) Copy() RequestTester {
	return RequestTester{
		handler: tReq.handler,
		ts:      tReq.ts,
		method:  tReq.method,
		path:    tReq.path,
		body:    tReq.body,
		query:   tReq.query,
		cookies: tReq.cookies,
	}
}
