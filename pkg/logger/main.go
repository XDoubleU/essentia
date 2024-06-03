package logger

import (
	"log"
	"os"
)

var logger *log.Logger

func GetLogger() *log.Logger {
	if logger == nil {
		logger = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	}

	return logger
}

func SetLogger(newLogger *log.Logger) {
	logger = newLogger
}
