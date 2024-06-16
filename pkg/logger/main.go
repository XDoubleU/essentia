package logger

import (
	"io"
	"log"
	"os"
)

//nolint:gochecknoglobals // on purpose
var logger *log.Logger

//nolint:gochecknoglobals // on purpose
var NullLogger = log.New(io.Discard, "", 0)

func GetLogger() *log.Logger {
	if logger == nil {
		logger = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	}

	return logger
}

func SetLogger(newLogger *log.Logger) {
	logger = newLogger
}
