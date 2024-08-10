package main

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xdoubleu/essentia/pkg/config"
	"github.com/xdoubleu/essentia/pkg/logging"
	"github.com/xdoubleu/essentia/pkg/test"
)

func TestHealth(t *testing.T) {
	app := NewApp(logging.NewNopLogger())
	app.config.Env = config.TestEnv

	tReq := test.CreateRequestTester(app.Routes(), http.MethodGet, "/health")
	rs := tReq.Do(t)

	assert.Equal(t, http.StatusOK, rs.StatusCode)
}
