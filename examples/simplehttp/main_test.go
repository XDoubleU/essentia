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
	cfg := NewConfig()
	cfg.Env = config.TestEnv

	app := NewApp(logging.NewNopLogger(), cfg)

	tReq := test.CreateRequestTester(app.Routes(), http.MethodGet, "/health")
	rs := tReq.Do(t)

	assert.Equal(t, http.StatusOK, rs.StatusCode)
}
