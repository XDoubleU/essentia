package test

import (
	"net/http"
	"testing"

	"github.com/XDoubleU/essentia/pkg/http_tools"
	"github.com/stretchr/testify/assert"
)

type ValidatorTester struct {
	invalidTestRequests   []TestRequest
	expectedErrorMessages []map[string]interface{}
}

func CreateValidatorTester(t *testing.T) ValidatorTester {
	return ValidatorTester{
		invalidTestRequests:   []TestRequest{},
		expectedErrorMessages: []map[string]interface{}{},
	}
}

func (vt *ValidatorTester) AddTestCase(invalidTestRequest TestRequest, expectedErrorMessage map[string]interface{}) {
	vt.invalidTestRequests = append(vt.invalidTestRequests, invalidTestRequest)
	vt.expectedErrorMessages = append(vt.expectedErrorMessages, expectedErrorMessage)
}

func (vt ValidatorTester) Do(t *testing.T) {
	t.Helper()

	for i := 0; i < len(vt.invalidTestRequests); i++ {
		var rsData http_tools.ErrorDto
		rs := vt.invalidTestRequests[i].Do(t, &rsData)

		assert.Equal(t, http.StatusUnprocessableEntity, rs.StatusCode)
		assert.Equal(t, vt.expectedErrorMessages[i], rsData.Message)
	}
}
