package main

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	httptools "github.com/xdoubleu/essentia/pkg/communication/http"
	"github.com/xdoubleu/essentia/pkg/config"
	"github.com/xdoubleu/essentia/pkg/database/postgres"
	"github.com/xdoubleu/essentia/pkg/logging"
	"github.com/xdoubleu/essentia/pkg/test"
)

func TestHealth(t *testing.T) {
	cfg := NewConfig()
	cfg.Env = config.TestEnv

	logger := logging.NewNopLogger()

	db, err := postgres.Connect(
		logger,
		cfg.DBDsn,
		25, //nolint:mnd //no magic number
		"15m",
		30,             //nolint:mnd //no magic number
		30*time.Second, //nolint:mnd //no magic number
		5*time.Minute,  //nolint:mnd //no magic number
	)
	if err != nil {
		panic(err)
	}

	//todo use synctx

	app := NewApp(logger, cfg, db)

	tReq := test.CreateRequestTester(app.Routes(), http.MethodGet, "/health")
	rs := tReq.Do(t)

	var rsData Health
	err = httptools.ReadJSON(rs.Body, &rsData)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, rs.StatusCode)
	assert.Equal(t, true, rsData.IsDatabaseActive)
}
