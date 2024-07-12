package test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/XDoubleU/essentia/pkg/httptools"
)

type RequestTester struct {
	handler http.Handler
	ts      *httptest.Server
	method  string
	path    string
	reqData any
	query   map[string]string
	cookies []*http.Cookie
}

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

func (tReq *RequestTester) SetTestServer(ts *httptest.Server) {
	tReq.ts = ts
}

func (tReq *RequestTester) SetReqData(reqData any) {
	tReq.reqData = reqData
}

func (tReq *RequestTester) SetQuery(query map[string]string) {
	tReq.query = query
}

func (tReq *RequestTester) AddCookie(cookie *http.Cookie) {
	tReq.cookies = append(tReq.cookies, cookie)
}

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

	if tReq.reqData != nil {
		body, err = json.Marshal(tReq.reqData)
		if err != nil {
			t.Errorf("error when marshalling reqData: %v", err)
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

func (tReq RequestTester) Copy() RequestTester {
	return RequestTester{
		handler: tReq.handler,
		ts:      tReq.ts,
		method:  tReq.method,
		path:    tReq.path,
		reqData: tReq.reqData,
		query:   tReq.query,
		cookies: tReq.cookies,
	}
}
