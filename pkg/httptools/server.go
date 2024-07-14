package httptools

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Serve calls [http.Server.ListenAndServe] with some more fluff
// around it to handle unexpected shutdowns nicely.
func Serve(logger *log.Logger, srv *http.Server, environment string) error {
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		logger.Printf("shutting down server %s", s.String())

		ctx, cancel := context.WithTimeout(
			context.Background(),
			30*time.Second, //nolint:mnd // no magic number
		)
		defer cancel()

		//nolint:errcheck // not useful to capture for now
		srv.Shutdown(ctx)
	}()

	logger.Printf("starting %s server on %s", environment, srv.Addr)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	logger.Printf("stopped server %s", srv.Addr)

	return nil
}
