package logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
)

type Logger struct {
	logger *slog.Logger
}

func New(level string, output io.Writer, fileName string) *Logger {
	// Открываем файл для записи логов. Если файл существует, он будет открыт в режиме добавления.
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o666)
	if err != nil {
		panic(fmt.Errorf("could not open log file: %w", err))
	}

	// Используем MultiWriter для записи логов и в файл, и в предоставленный output (например, os.Stdout).
	writer := io.MultiWriter(file, output)

	// Устанавливаем уровень логирования в зависимости от входного параметра level.
	var opts slog.HandlerOptions
	switch level {
	case "debug":
		opts.Level = slog.LevelDebug
	case "info":
		opts.Level = slog.LevelInfo
	case "warn":
		opts.Level = slog.LevelWarn
	case "error":
		opts.Level = slog.LevelError
	default:
		opts.Level = slog.LevelInfo // уровень по умолчанию.
	}
	logger := slog.New(slog.NewTextHandler(writer, &opts))
	return &Logger{logger: logger}
}

func (l Logger) Debug(msg string) {
	l.logger.Debug(msg)
}

func (l Logger) Info(msg string) {
	l.logger.Info(msg)
}

func (l Logger) Warn(msg string) {
	l.logger.Warn(msg)
}

func (l Logger) Error(msg string) {
	l.logger.Error(msg)
}
