package middleware

import (
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RecoveryMiddleware provides panic recovery with zap logging
type RecoveryMiddleware struct {
	logger *zap.Logger
}

// NewRecoveryMiddleware creates a new recovery middleware
func NewRecoveryMiddleware(logger *zap.Logger) *RecoveryMiddleware {
	return &RecoveryMiddleware{
		logger: logger,
	}
}

// Recover returns a gin middleware that recovers from panics and logs using zap
// This is a wrapper around the ginzap.RecoveryWithZap middleware
func (m *RecoveryMiddleware) Recover() gin.HandlerFunc {
	return ginzap.RecoveryWithZap(m.logger, true)
}
