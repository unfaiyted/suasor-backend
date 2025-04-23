package database

import (
	"context"
	"os"

	"github.com/rs/zerolog"
)

type loggerCtxKey struct{}

var testLoggerKey = loggerCtxKey{}

// NewTestLogger creates a new test logger
func NewTestLogger() zerolog.Logger {
	// For tests that need to see output
	output := zerolog.ConsoleWriter{Out: os.Stdout, NoColor: true}
	
	// For tests that should run quietly, use this in production tests:
	// output = zerolog.ConsoleWriter{Out: io.Discard, NoColor: true}
	
	return zerolog.New(output).With().Timestamp().Logger()
}

// ContextWithTestLogger adds logger to context
func ContextWithTestLogger(ctx context.Context, logger zerolog.Logger) context.Context {
	return context.WithValue(ctx, testLoggerKey, logger)
}

// TestLoggerFromContext extracts logger from context
func TestLoggerFromContext(ctx context.Context) zerolog.Logger {
	if ctx == nil {
		return zerolog.Nop()
	}
	if logger, ok := ctx.Value(testLoggerKey).(zerolog.Logger); ok {
		return logger
	}
	return zerolog.Nop()
}