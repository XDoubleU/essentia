package main

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xdoubleu/essentia/pkg/config"
	"github.com/xdoubleu/essentia/pkg/logging"
	"github.com/xdoubleu/essentia/pkg/test"
)

func TestHealth(t *testing.T) {
	app := NewApp(logging.NewNopLogger())
	app.config.Env = config.TestEnv

	routes, err := app.Routes()
	require.Nil(t, err)

	tReq := test.CreateRequestTester(*routes, http.MethodGet, "/health")
	rs := tReq.Do(t, nil)

	assert.Equal(t, http.StatusOK, rs.StatusCode)
}
