package logger

import (
	"context"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type ctxKey struct{}

var loggerKey = ctxKey{}

// Initialize sets up the global logger with default level (Debug)
func Initialize() {
	InitializeWithLevel(zerolog.DebugLevel)
}

// InitializeWithLevel sets up the global logger with specified level
func InitializeWithLevel(level zerolog.Level) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(level)
	// Enable caller information (file and line)
	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: "15:04:05",
	}
	log.Logger = zerolog.New(consoleWriter).
		With().
		Timestamp().
		Caller().
		Logger()
}

// SetLogLevel changes the global log level
func SetLogLevel(level zerolog.Level) {
	zerolog.SetGlobalLevel(level)
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

func WithJobID(ctx context.Context, jobID uint64) (context.Context, zerolog.Logger) {
	logger := LoggerFromContext(ctx).With().Uint64("job_id", jobID).Logger()
	return WithContext(ctx, logger), logger
}

// WithCaller adds file and line information to the logger
func WithCaller() zerolog.Logger {
	return log.Logger.With().Caller().Logger()
}
