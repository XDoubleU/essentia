package main

import (
	"context"
	"net/http"
	"testing"
	"time"

	httptools "github.com/XDoubleU/essentia/pkg/communication/http"
	"github.com/XDoubleU/essentia/pkg/config"
	"github.com/XDoubleU/essentia/pkg/database/postgres"
	"github.com/XDoubleU/essentia/pkg/logging"
	"github.com/XDoubleU/essentia/pkg/test"
	"github.com/stretchr/testify/assert"
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

	tx := postgres.CreatePgxSyncTx(context.Background(), db)
	defer tx.Rollback(context.Background())

	app := NewApp(logger, cfg, tx)

	tReq := test.CreateRequestTester(app.Routes(), http.MethodGet, "/health")
	rs := tReq.Do(t)

	var rsData Health
	err = httptools.ReadJSON(rs.Body, &rsData)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, rs.StatusCode)
	assert.Equal(t, true, rsData.IsDatabaseActive)
}
