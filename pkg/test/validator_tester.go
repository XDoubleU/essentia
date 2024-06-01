package test

import (
	"net/http"
	"testing"

	"github.com/XDoubleU/essentia/pkg/http_tools"
	"github.com/stretchr/testify/assert"
)

type ValidatorTester struct {
	t                     *testing.T
	invalidTestRequests   []TestRequest
	expectedErrorMessages []map[string]interface{}
}

func CreateValidatorTester(t *testing.T) ValidatorTester {
	return ValidatorTester{
		t:                     t,
		invalidTestRequests:   []TestRequest{},
		expectedErrorMessages: []map[string]interface{}{},
	}
}

func (vt *ValidatorTester) AddTestCase(invalidTestRequest TestRequest, expectedErrorMessage map[string]interface{}) {
	vt.invalidTestRequests = append(vt.invalidTestRequests, invalidTestRequest)
	vt.expectedErrorMessages = append(vt.expectedErrorMessages, expectedErrorMessage)
}

func (vt ValidatorTester) Do() {
	vt.t.Helper()

	for i := 0; i < len(vt.invalidTestRequests); i++ {
		var rsData http_tools.ErrorDto
		rs := vt.invalidTestRequests[i].Do(&rsData)

		assert.Equal(vt.t, http.StatusUnprocessableEntity, rs.StatusCode)
		assert.Equal(vt.t, vt.expectedErrorMessages[i], rsData.Message)
	}
}
