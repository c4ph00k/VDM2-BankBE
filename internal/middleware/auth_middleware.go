package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"VDM2-BankBE/internal/service"
	"VDM2-BankBE/internal/util"
)

// AuthMiddleware provides JWT authentication middleware for Gin
type AuthMiddleware struct {
	authService service.AuthService
	logger      *zap.Logger
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(authService service.AuthService, logger *zap.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
		logger:      logger,
	}
}

// Authenticate verifies the JWT token and adds the user to the context
func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from Authorization header
		token := extractToken(c.Request)
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, util.ErrorResponse{
				Error: util.NewUnauthorizedError("missing authorization token"),
			})
			return
		}

		// Verify token and get user
		user, err := m.authService.VerifyToken(c, token)
		if err != nil {
			m.logger.Error("failed to verify token", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusUnauthorized, util.ErrorResponse{
				Error: util.NewUnauthorizedError("invalid or expired token"),
			})
			return
		}

		// Add user to context
		c.Set("user", user)
		c.Next()
	}
}

// extractToken extracts JWT token from Authorization header
func extractToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	// Check for Bearer prefix
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}
