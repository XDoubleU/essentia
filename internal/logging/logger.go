package logging

import (
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
)

type Logger interface {
	Trace(...any)
	Debug(...any)
	Info(...any)
	Notice(...any)
	Warn(...any)
	Error(...any)
	Fatal(...any)
	Tracef(string, ...any)
	Debugf(string, ...any)
	Infof(string, ...any)
	Noticef(string, ...any)
	Warnf(string, ...any)
	Errorf(string, ...any)
	Fatalf(string, ...any)
}

func NewLogger() Logger {
	//TODO: loglevel magic
	logger := slog.NewWithHandlers(handler.NewConsoleHandler(slog.AllLevels))
	logger.ChannelName = "essentia"
	return logger
}
