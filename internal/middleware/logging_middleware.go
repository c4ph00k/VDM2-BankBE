package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// LoggingMiddleware provides request logging
type LoggingMiddleware struct {
	logger *zap.Logger
}

// NewLoggingMiddleware creates a new logging middleware
func NewLoggingMiddleware(logger *zap.Logger) *LoggingMiddleware {
	return &LoggingMiddleware{
		logger: logger,
	}
}

// LogRequest returns a Gin middleware that logs request details
func (m *LoggingMiddleware) LogRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate a request ID
		requestID := uuid.New().String()

		// Add request ID to context
		c.Set("request_id", requestID)

		// Log request start
		startTime := time.Now()
		m.logger.Info("request started",
			zap.String("request_id", requestID),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("remote_addr", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
		)

		// Process request
		c.Next()

		// Log request completion
		duration := time.Since(startTime)
		m.logger.Info("request completed",
			zap.String("request_id", requestID),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("duration_ms", duration),
		)
	}
}
