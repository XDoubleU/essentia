package http_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	httptools "github.com/XDoubleU/essentia/pkg/communication/http"
	"github.com/stretchr/testify/assert"
)

func TestStatus(t *testing.T) {
	res := httptest.NewRecorder()
	rw := httptools.NewResponseWriter(res)
	rw.WriteHeader(http.StatusOK)
	assert.Equal(t, http.StatusOK, rw.Status())
}
