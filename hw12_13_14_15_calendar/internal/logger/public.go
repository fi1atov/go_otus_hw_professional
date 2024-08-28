package logger

import (
	"io"
	"log/slog"
)

type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

func New(level string, output io.Writer) Logger {
	// Используем MultiWriter для записи логов и в файл, и в предоставленный output (например, os.Stdout).
	writer := io.MultiWriter(output)

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
	log := slog.New(slog.NewTextHandler(writer, &opts))
	return logger{logger: log}
}
