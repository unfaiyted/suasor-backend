package middleware

import (
	logger "suasor/utils/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// In your middleware
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate request ID
		requestID := uuid.New().String()

		// Add request ID to response headers
		c.Header("X-Request-ID", requestID)

		// Create logger with request context
		ctx, reqLogger := logger.WithRequestID(c.Request.Context(), requestID)

		// Add more fields related to request
		reqLogger = reqLogger.With().
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Logger()

		// Store in context
		c.Request = c.Request.WithContext(logger.WithContext(ctx, reqLogger))

		// Log request start
		reqLogger.Info().Msg("Request started")

		c.Next()

		// Log request end
		reqLogger.Info().
			Int("status", c.Writer.Status()).
			Msg("Request completed")
	}
}
