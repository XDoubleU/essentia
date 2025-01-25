package test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	httptools "github.com/XDoubleU/essentia/pkg/communication/http"
)

type ContentType = int

const (
	JsonContentType ContentType = iota
	FormContentType             = iota
)

// A RequestTester is used to test a certain HTTP request.
type RequestTester struct {
	handler     http.Handler
	ts          *httptest.Server
	contentType ContentType
	method      string
	path        string
	data        any
	query       url.Values
	cookies     []*http.Cookie
}

// CreateRequestTester creates a new [RequestTester].
func CreateRequestTester(
	handler http.Handler,
	contentType ContentType,
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
		contentType,
		method,
		path,
		nil,
		make(url.Values),
		[]*http.Cookie{},
	}
}

// SetTestServer sets the test server of a [RequestTester].
// This allows to reuse existing test servers.
func (tReq *RequestTester) SetTestServer(ts *httptest.Server) {
	tReq.ts = ts
}

// SetData sets the request data of a [RequestTester].
func (tReq *RequestTester) SetData(data any) {
	tReq.data = data
}

// SetQuery sets the query of a [RequestTester].
func (tReq *RequestTester) SetQuery(query url.Values) {
	tReq.query = query
}

// AddCookie adds a cookie to a [RequestTester].
// Can be used multiple times for adding several cookies.
func (tReq *RequestTester) AddCookie(cookie *http.Cookie) {
	tReq.cookies = append(tReq.cookies, cookie)
}

// Do executes a [RequestTester] returning the response of a request
// and providing the returned data to rsData.
func (tReq RequestTester) Do(t *testing.T) *http.Response {
	t.Helper()

	var contentType string
	var bodyReader io.Reader
	var err error

	if tReq.ts == nil && tReq.handler != nil {
		tReq.ts = httptest.NewServer(tReq.handler)
		defer tReq.ts.Close()
	}

	if tReq.ts == nil {
		panic("handler nor test server has been set")
	}

	if tReq.data != nil {
		switch tReq.contentType {
		case JsonContentType:
			contentType = "application/json"

			var body []byte
			body, err = json.Marshal(tReq.data)
			bodyReader = bytes.NewReader(body)
		case FormContentType:
			contentType = "application/x-www-form-urlencoded"

			var query url.Values
			query, err = httptools.WriteForm(tReq.data)
			bodyReader = strings.NewReader(query.Encode())
		}

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
		bodyReader,
	)

	if err != nil {
		t.Errorf("error when creating request: %v", err)
		t.FailNow()
		return nil
	}

	if len(tReq.query) > 0 {
		query := req.URL.Query()

		for key, values := range tReq.query {
			for _, value := range values {
				query.Add(key, value)
			}
		}

		req.URL.RawQuery = query.Encode()
	}

	for _, cookie := range tReq.cookies {
		req.AddCookie(cookie)
	}

	req.Header.Set("Content-Type", contentType)

	rs, err := tReq.ts.Client().Do(req)
	if err != nil {
		t.Errorf("error when making request: %v", err)
		t.FailNow()
		return nil
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
		data:    tReq.data,
		query:   tReq.query,
		cookies: tReq.cookies,
	}
}
