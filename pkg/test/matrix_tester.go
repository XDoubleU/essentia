package test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	httptools "github.com/xdoubleu/essentia/pkg/communication/http"
)

// CaseResponse is used to compare to the actual response
// of a [RequestTester] when used by a [MatrixTester].
type CaseResponse struct {
	statusCode int
	cookies    []*http.Cookie
	body       *map[string]any
}

// NewCaseResponse returns a new [CaseResponse].
func NewCaseResponse(statusCode int) CaseResponse {
	//nolint:exhaustruct //other fields are optional
	return CaseResponse{
		statusCode: statusCode,
	}
}

// SetCookies sets the cookies expected in the response of a test case.
func (rs *CaseResponse) SetCookies(cookies []*http.Cookie) {
	rs.cookies = cookies
}

// SetBody sets the body expected in the response of a test case.
func (rs *CaseResponse) SetBody(body any) {
	var bodyMap map[string]any
	data, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(data, &bodyMap)
	if err != nil {
		panic(err)
	}

	rs.body = &bodyMap
}

// MatrixTester is used for executing matrix tests.
type MatrixTester struct {
	testCases map[*RequestTester]CaseResponse
}

// CreateMatrixTester creates a new [MatrixTester].
func CreateMatrixTester() MatrixTester {
	return MatrixTester{
		testCases: make(map[*RequestTester]CaseResponse),
	}
}

// AddTestCase adds a test case which consists of
// a [RequestTester] and a [CaseResponse].
// When executing [Do] the [RequestTester] will be executed
// and its response will be compared to the provided [CaseResponse].
func (mt *MatrixTester) AddTestCase(
	tReq RequestTester,
	tRes CaseResponse) {
	mt.testCases[&tReq] = tRes
}

// Do executes a [MatrixTester].
func (mt MatrixTester) Do(t *testing.T) {
	t.Helper()

	for tReq, tRes := range mt.testCases {
		var rs *http.Response
		var rsData map[string]any

		if tRes.body == nil {
			rs = tReq.Do(t, nil, nil)
			defer rs.Body.Close()
		} else {
			//todo dynamic content type adapter
			rs = tReq.Do(t, &rsData, httptools.ReadJSON)
			defer rs.Body.Close()
		}

		assert.Equal(t, tRes.statusCode, rs.StatusCode)

		if tRes.cookies != nil {
			for _, cookie := range tRes.cookies {
				found := false

				for _, acCookie := range rs.Cookies() {
					found = cookie.String() == acCookie.String()
				}

				assert.True(t, found)
			}
		}

		if rsData != nil {
			assert.Equal(t, *tRes.body, rsData)
		}
	}
}
