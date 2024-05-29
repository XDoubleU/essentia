package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/XDoubleU/essentia/internal/helpers"
)

type TestRequest struct {
	t       *testing.T
	ts      *httptest.Server
	method  string
	path    string
	reqData any
	cookies []*http.Cookie
}

func CreateTestRequest(t *testing.T, ts *httptest.Server, method, path string) TestRequest {
	return TestRequest{
		t,
		ts,
		method,
		path,
		nil,
		[]*http.Cookie{},
	}
}

func (tReq *TestRequest) SetReqData(reqData any) {
	tReq.reqData = reqData
}

func (tReq *TestRequest) AddCookie(cookie *http.Cookie) {
	tReq.cookies = append(tReq.cookies, cookie)
}

func (tReq *TestRequest) Do(rsData any) *http.Response {
	var body []byte
	var err error

	if tReq.reqData != nil {
		body, err = json.Marshal(tReq.reqData)
		if err != nil {
			tReq.t.Errorf("error when marshalling reqData: %v", err)
			tReq.t.FailNow()
			return nil
		}
	}

	req, err := http.NewRequest(
		tReq.method,
		fmt.Sprintf("%s/%s", tReq.ts.URL, tReq.path),
		bytes.NewReader(body),
	)

	if err != nil {
		tReq.t.Errorf("error when creating request: %v", err)
		tReq.t.FailNow()
		return nil
	}

	for _, cookie := range tReq.cookies {
		req.AddCookie(cookie)
	}

	rs, err := tReq.ts.Client().Do(req)
	if err != nil {
		tReq.t.Errorf("error when making request: %v", err)
		tReq.t.FailNow()
		return nil
	}

	if rsData != nil {
		err = helpers.ReadJSON(rs.Body, &rsData, false)
		if err != nil {
			tReq.t.Errorf("error when parsing response: %v", err)
			tReq.t.FailNow()
			return nil
		}
	}

	return rs
}
