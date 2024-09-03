package logger

import (
	"log/slog"
)

type logger struct {
	logger *slog.Logger
}

func (l logger) Debug(msg string, args ...interface{}) {
	l.logger.Debug(msg, args...)
}

func (l logger) Info(msg string, args ...interface{}) {
	l.logger.Info(msg, args...)
}

func (l logger) Warn(msg string, args ...interface{}) {
	l.logger.Warn(msg, args...)
}

func (l logger) Error(msg string, args ...interface{}) {
	l.logger.Error(msg, args...)
}
