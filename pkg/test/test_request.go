package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/XDoubleU/essentia/pkg/http_tools"
)

type TestRequest struct {
	handler http.Handler
	ts      *httptest.Server
	method  string
	path    string
	reqData any
	query   map[string]string
	cookies []*http.Cookie
}

func CreateTestRequest(handler http.Handler, method, path string, pathValues ...any) TestRequest {
	if len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}

	if len(pathValues) > 0 {
		path = fmt.Sprintf(path, pathValues...)
	}

	return TestRequest{
		handler,
		nil,
		method,
		path,
		nil,
		make(map[string]string),
		[]*http.Cookie{},
	}
}

func (tReq *TestRequest) SetTestServer(ts *httptest.Server) {
	tReq.ts = ts
}

func (tReq *TestRequest) SetReqData(reqData any) {
	tReq.reqData = reqData
}

func (tReq *TestRequest) SetQuery(query map[string]string) {
	tReq.query = query
}

func (tReq *TestRequest) AddCookie(cookie *http.Cookie) {
	tReq.cookies = append(tReq.cookies, cookie)
}

func (tReq TestRequest) Do(t *testing.T, rsData any) *http.Response {
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

	req, err := http.NewRequest(
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
		err = http_tools.ReadJSON(rs.Body, &rsData)
		if err != nil {
			t.Errorf("error when parsing response: %v", err)
			t.FailNow()
			return nil
		}
	}

	return rs
}

func (tReq TestRequest) Copy() TestRequest {
	return TestRequest{
		handler: tReq.handler,
		ts:      tReq.ts,
		method:  tReq.method,
		path:    tReq.path,
		reqData: tReq.reqData,
		query:   tReq.query,
		cookies: tReq.cookies,
	}
}
