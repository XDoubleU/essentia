package sentry_test

import (
	"bytes"
	"errors"
	"log/slog"
	"testing"

	"github.com/XDoubleU/essentia/pkg/config"
	"github.com/XDoubleU/essentia/pkg/logging"
	"github.com/XDoubleU/essentia/pkg/sentry"
	"github.com/stretchr/testify/assert"
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
