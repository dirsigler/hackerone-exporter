package logger

import (
	"log/slog"
	"os"
)

// Setup configures and sets the default logger with the specified level.
func Setup(level slog.Level) {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)
}
