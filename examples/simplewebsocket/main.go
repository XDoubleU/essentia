package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	httptools "github.com/xdoubleu/essentia/pkg/communication/http"
	"github.com/xdoubleu/essentia/pkg/logging"
	sentrytools "github.com/xdoubleu/essentia/pkg/sentry"
)

type application struct {
	logger *slog.Logger
	config Config
}

func NewApp(logger *slog.Logger, config Config) application {
	return application{
		logger: logger,
		config: config,
	}
}

func main() {
	cfg := NewConfig()

	logger := slog.New(sentrytools.NewLogHandler(cfg.Env, slog.NewTextHandler(os.Stdout, nil)))

	app := NewApp(logger, cfg)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.Port),
		Handler:      app.Routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,  //nolint:gomnd //no magic number
		WriteTimeout: 10 * time.Second, //nolint:gomnd //no magic number
	}

	err := httptools.Serve(logger, srv, app.config.Env)
	if err != nil {
		logger.Error("failed to serve server", logging.ErrAttr(err))
	}
}
