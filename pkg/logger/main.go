package logger

import (
	"io"
	"log"
	"os"
)

var logger *log.Logger

var NullLogger *log.Logger = log.New(io.Discard, "", 0)

func GetLogger() *log.Logger {
	if logger == nil {
		logger = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	}

	return logger
}

func SetLogger(newLogger *log.Logger) {
	logger = newLogger
}