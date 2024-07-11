package httptools_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/XDoubleU/essentia/pkg/httptools"
	"github.com/stretchr/testify/assert"
)

func TestServerFunc(t *testing.T) {
	mux := http.NewServeMux()

	srv := &http.Server{
		Addr:         "localhost:8000",
		Handler:      mux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		time.Sleep(time.Second)
		err := srv.Shutdown(context.Background())
		assert.Nil(t, err)
	}()

	err := httptools.Serve(srv, "test")

	assert.Nil(t, err)
}
