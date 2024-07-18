package main

import (
	"io"
	"log"
	"net/http"
	"testing"

	"github.com/XDoubleU/essentia/pkg/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealth(t *testing.T) {
	app := NewApp(log.New(io.Discard, "", 0))
	app.config.Env = TestEnv

	routes, err := app.Routes()
	require.Nil(t, err)

	tReq := test.CreateRequestTester(*routes, http.MethodGet, "/health")
	rs := tReq.Do(t, nil)

	assert.Equal(t, http.StatusOK, rs.StatusCode)
}
