package httptools_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/XDoubleU/essentia/pkg/httptools"
	"github.com/stretchr/testify/assert"
)

func TestStatusCode(t *testing.T) {
	res := httptest.NewRecorder()
	rw := httptools.NewResponseWriter(res)

	rw.WriteHeader(http.StatusOK)

	assert.Equal(t, http.StatusOK, rw.StatusCode())
}
