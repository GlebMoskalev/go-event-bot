package logger

import (
	"fmt"
	"log/slog"
	"os"
)

const (
	appEnvLocal = "local"
	appEnvDev   = "dev"
	appEnvProd  = "prod"
)

func NewLogger(appEnv string) *slog.Logger {
	var log *slog.Logger
	switch appEnv {
	case appEnvLocal:
		log = slog.New(newPrettySlog(slog.LevelDebug))
	case appEnvDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case appEnvProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn}))
	default:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
		log.Warn(fmt.Sprintf("unknown APP_ENV value '%s', falling back to default configuration", appEnv))
	}
	return log
}

func SetupLogger(baseLogger *slog.Logger, layer, operation string, extraFields ...any) *slog.Logger {
	log := baseLogger.With("layer", layer, "operation", operation)
	if len(extraFields) > 0 {
		log = log.With(extraFields...)
	}
	return log
}
