package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	httptools "github.com/XDoubleU/essentia/pkg/communication/http"
	"github.com/XDoubleU/essentia/pkg/database/postgres"
	"github.com/XDoubleU/essentia/pkg/logging"
	sentrytools "github.com/XDoubleU/essentia/pkg/sentry"
)

type application struct {
	logger *slog.Logger
	config Config
	db     postgres.DB
}

func NewApp(logger *slog.Logger, config Config, db postgres.DB) application {
	spandb := postgres.NewSpanDB(db)

	return application{
		logger: logger,
		config: config,
		db:     spandb,
	}
}

func main() {
	cfg := NewConfig()

	logger := slog.New(
		sentrytools.NewLogHandler(cfg.Env, slog.NewTextHandler(os.Stdout, nil)),
	)
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

	app := NewApp(logger, cfg, db)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.Port),
		Handler:      app.Routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,  //nolint:gomnd //no magic number
		WriteTimeout: 10 * time.Second, //nolint:gomnd //no magic number
	}

	err = httptools.Serve(logger, srv, app.config.Env)
	if err != nil {
		logger.Error("failed to serve server", logging.ErrAttr(err))
	}
}
