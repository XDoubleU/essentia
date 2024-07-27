package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/xdoubleu/essentia/pkg/httptools"
	"github.com/xdoubleu/essentia/pkg/logging"
	"github.com/xdoubleu/essentia/pkg/sentrytools"
)

type application struct {
	logger *slog.Logger
	config Config
}

func NewApp(logger *slog.Logger) application {
	return application{
		logger: logger,
		config: NewConfig(),
	}
}

func main() {
	logger := slog.New(sentrytools.NewSentryLogHandler())
	app := NewApp(logger)

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
