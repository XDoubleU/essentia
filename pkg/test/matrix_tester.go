package test

import (
	"net/http"
	"testing"

	"github.com/XDoubleU/essentia/pkg/httptools"
	"github.com/stretchr/testify/assert"
)

// ErrorMessage is used to describe the expected error in [AddTestCaseErrorMessage].
type ErrorMessage = map[string]interface{}

// MatrixTester is used for executing matrix tests.
type MatrixTester struct {
	baseRequest       RequestTester
	errorMessageTests map[*RequestTester]ErrorMessage
	statusCodeTests   map[*RequestTester]int
}

// CreateMatrixTester creates a new [MatrixTester].
func CreateMatrixTester(baseRequest RequestTester) MatrixTester {
	return MatrixTester{
		baseRequest:       baseRequest,
		errorMessageTests: make(map[*RequestTester]ErrorMessage),
		statusCodeTests:   make(map[*RequestTester]int),
	}
}

// AddTestCaseErrorMessage adds a testcase where request data is
// supplied as input and a certain error message is expected as output.
func (mt *MatrixTester) AddTestCaseErrorMessage(
	reqData any,
	errorMessage ErrorMessage,
) {
	tReq := mt.baseRequest.Copy()
	tReq.SetReqData(reqData)
	mt.errorMessageTests[&tReq] = errorMessage
}

// AddTestCaseStatusCode adds a testcase where a query is
// supplied as input and a certain status code is expected as output.
func (mt *MatrixTester) AddTestCaseStatusCode(query map[string]string, statusCode int) {
	tReq := mt.baseRequest.Copy()
	tReq.SetQuery(query)
	mt.statusCodeTests[&tReq] = statusCode
}

// AddTestCaseCookieStatusCode adds a testcase where a cookie is
// supplied as input and a certain status code is expected as output.
func (mt *MatrixTester) AddTestCaseCookieStatusCode(
	cookie *http.Cookie,
	statusCode int,
) {
	tReq := mt.baseRequest.Copy()

	if cookie != nil {
		tReq.AddCookie(cookie)
	}

	mt.statusCodeTests[&tReq] = statusCode
}

// Do executes a [MatrixTester].
func (mt MatrixTester) Do(t *testing.T) {
	t.Helper()

	for tReq, errorMsg := range mt.errorMessageTests {
		var rsData httptools.ErrorDto
		rs := tReq.Do(t, &rsData)
		defer rs.Body.Close()

		assert.Equal(t, errorMsg, rsData.Message)
	}

	for tReq, statusCode := range mt.statusCodeTests {
		rs := tReq.Do(t, nil)
		defer rs.Body.Close()

		assert.Equal(t, statusCode, rs.StatusCode)
	}
}
