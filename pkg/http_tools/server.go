package http_tools

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Serve(port int, handler http.Handler, environment string) error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      handler,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,  //nolint:gomnd //no magic number
		WriteTimeout: 10 * time.Second, //nolint:gomnd //no magic number
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		GetLogger().Printf("shutting down server %s", s.String())

		ctx, cancel := context.WithTimeout(
			context.Background(),
			30*time.Second, //nolint:gomnd // no magic number
		)
		defer cancel()

		shutdownError <- srv.Shutdown(ctx)
	}()

	GetLogger().Printf("starting %s server on %s", environment, srv.Addr)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	GetLogger().Printf("stopped server %s", srv.Addr)

	return nil
}
