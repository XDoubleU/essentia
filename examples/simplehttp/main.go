package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/XDoubleU/essentia/pkg/httptools"
)

type application struct {
	logger *log.Logger
	config Config
}

func NewApp(logger *log.Logger) application {
	return application{
		logger: logger,
		config: NewConfig(),
	}
}

func main() {
	logger := log.Default()

	app := NewApp(logger)

	routes, err := app.Routes()
	if err != nil {
		logger.Fatal(err)
		return
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.Port),
		Handler:      *routes,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,  //nolint:gomnd //no magic number
		WriteTimeout: 10 * time.Second, //nolint:gomnd //no magic number
	}

	err = httptools.Serve(logger, srv, app.config.Env)
	if err != nil {
		logger.Fatal(err)
	}
}
