package utils

import (
	"context"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type ctxKey struct{}

var loggerKey = ctxKey{}

// Initialize sets up the global logger
func Initialize() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

// FromContext extracts logger from context
func LoggerFromContext(ctx context.Context) zerolog.Logger {
	if ctx == nil {
		return log.Logger
	}
	if logger, ok := ctx.Value(loggerKey).(zerolog.Logger); ok {
		return logger
	}
	return log.Logger
}

// WithContext adds logger to context
func WithContext(ctx context.Context, logger zerolog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// WithRequestID adds request ID to logger and context
func WithRequestID(ctx context.Context, requestID string) (context.Context, zerolog.Logger) {
	logger := LoggerFromContext(ctx).With().Str("request_id", requestID).Logger()
	return WithContext(ctx, logger), logger
}
