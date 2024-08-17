package logging

import (
	"log/slog"
	"os"
)

const (
	// EnvLocal local environment
	EnvLocal = "local"
	// EnvProd prod environment
	EnvProd = "prod"
)

// SetupLogger getting new logger
func SetupLogger(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case EnvLocal:
		logger = slog.New(
			slog.NewTextHandler(
				os.Stdout,
				&slog.HandlerOptions{Level: slog.LevelDebug},
			),
		)
	case EnvProd:
		logger = slog.New(
			slog.NewJSONHandler(
				os.Stdout,
				&slog.HandlerOptions{Level: slog.LevelInfo},
			),
		)
	}

	return logger
}
