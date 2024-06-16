package test

import (
	"net/http"
	"testing"

	"github.com/XDoubleU/essentia/pkg/http_tools"
	"github.com/stretchr/testify/assert"
)

type ErrorMessage = map[string]interface{}

type MatrixTester struct {
	baseRequest       TestRequest
	errorMessageTests map[*TestRequest]ErrorMessage
	statusCodeTests   map[*TestRequest]int
}

func CreateMatrixTester(t *testing.T, baseRequest TestRequest) MatrixTester {
	return MatrixTester{
		baseRequest:       baseRequest,
		errorMessageTests: make(map[*TestRequest]ErrorMessage),
		statusCodeTests:   make(map[*TestRequest]int),
	}
}

func (mt *MatrixTester) AddTestCaseErrorMessage(reqData any, errorMessage ErrorMessage) {
	tReq := mt.baseRequest.Copy()
	tReq.SetReqData(reqData)
	mt.errorMessageTests[&tReq] = errorMessage
}

func (mt *MatrixTester) AddTestCaseStatusCode(query map[string]string, statusCode int) {
	tReq := mt.baseRequest.Copy()
	tReq.SetQuery(query)
	mt.statusCodeTests[&tReq] = statusCode
}

func (mt *MatrixTester) AddTestCaseCookieStatusCode(cookie *http.Cookie, statusCode int) {
	tReq := mt.baseRequest.Copy()

	if cookie != nil {
		tReq.AddCookie(cookie)
	}

	mt.statusCodeTests[&tReq] = statusCode
}

func (mt MatrixTester) Do(t *testing.T) {
	t.Helper()

	for tReq, errorMsg := range mt.errorMessageTests {
		var rsData http_tools.ErrorDto
		tReq.Do(t, &rsData)

		assert.Equal(t, errorMsg, rsData.Message)
	}

	for tReq, statusCode := range mt.statusCodeTests {
		rs := tReq.Do(t, nil)

		assert.Equal(t, statusCode, rs.StatusCode)
	}
}
