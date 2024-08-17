package sentry_test

import (
	"bytes"
	"errors"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xdoubleu/essentia/pkg/config"
	"github.com/xdoubleu/essentia/pkg/logging"
	"github.com/xdoubleu/essentia/pkg/sentry"
)

func TestLogHandlerDev(t *testing.T) {
	var buf bytes.Buffer

	logger := slog.New(
		sentry.NewLogHandler(
			config.DevEnv,
			//nolint:exhaustruct //other fields are optional
			slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug}),
		),
	)

	logger.Error("test", logging.ErrAttr(errors.New("testerror")))

	assert.Contains(t, buf.String(), "level=ERROR msg=test error=testerror")
}
