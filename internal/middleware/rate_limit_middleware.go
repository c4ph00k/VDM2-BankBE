package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"VDM2-BankBE/internal/config"
	"VDM2-BankBE/internal/model"
	"VDM2-BankBE/internal/util"
	"VDM2-BankBE/pkg/cache"
)

// RateLimitMiddleware limits request rates per user for Gin
type RateLimitMiddleware struct {
	redisClient *cache.RedisClient
	config      *config.RateLimitConfig
	logger      *zap.Logger
}

// NewRateLimitMiddleware creates a new rate limit middleware
func NewRateLimitMiddleware(
	redisClient *cache.RedisClient,
	config *config.RateLimitConfig,
	logger *zap.Logger,
) *RateLimitMiddleware {
	return &RateLimitMiddleware{
		redisClient: redisClient,
		config:      config,
		logger:      logger,
	}
}

// LimitByUser limits requests per user ID
func (m *RateLimitMiddleware) LimitByUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip if rate limiting is disabled
		if !m.config.Enabled {
			c.Next()
			return
		}

		// Get user from context
		user, exists := c.Get("user")
		if !exists {
			// If no user in context, let the request through
			// The auth middleware will handle unauthorized requests
			c.Next()
			return
		}

		userModel, ok := user.(*model.User)
		if !ok {
			m.logger.Error("failed to get user from context")
			c.Next()
			return
		}

		// Get the route for more granular limiting
		route := c.Request.Method + ":" + c.FullPath()

		// Increment and check rate limit
		count, err := m.redisClient.IncrRateLimit(c, userModel.ID.String(), route, m.config.Duration)
		if err != nil {
			m.logger.Error("failed to check rate limit", zap.Error(err))
			// Continue serving the request on Redis error
			c.Next()
			return
		}

		// Check if over the limit
		if count > int64(m.config.Requests) {
			m.logger.Warn("rate limit exceeded",
				zap.String("user_id", userModel.ID.String()),
				zap.String("route", route),
				zap.Int64("count", count),
				zap.Int("limit", m.config.Requests),
			)

			// Return rate limit exceeded error
			c.Header("Retry-After", "60")
			c.AbortWithStatusJSON(http.StatusTooManyRequests, util.ErrorResponse{
				Error: util.NewAPIError(http.StatusTooManyRequests, "rate limit exceeded"),
			})
			return
		}

		// Continue with the request
		c.Next()
	}
}

// LimitByIP limits requests per IP address
func (m *RateLimitMiddleware) LimitByIP() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip if rate limiting is disabled
		if !m.config.Enabled {
			c.Next()
			return
		}

		// Get client IP
		ip := c.ClientIP()

		// Get the route for more granular limiting
		route := c.Request.Method + ":" + c.FullPath()

		// Use IP address as identifier
		identifier := "ip:" + ip + ":" + route

		// Increment and check rate limit
		count, err := m.redisClient.IncrRateLimit(c, identifier, route, m.config.Duration)
		if err != nil {
			m.logger.Error("failed to check IP rate limit", zap.Error(err))
			// Continue serving the request on Redis error
			c.Next()
			return
		}

		// Check if over the limit
		// For IP-based limiting, we often want a higher limit than for authenticated requests
		ipLimit := m.config.Requests * 2
		if count > int64(ipLimit) {
			m.logger.Warn("IP rate limit exceeded",
				zap.String("ip", ip),
				zap.String("route", route),
				zap.Int64("count", count),
				zap.Int("limit", ipLimit),
			)

			// Return rate limit exceeded error
			c.Header("Retry-After", "60")
			c.AbortWithStatusJSON(http.StatusTooManyRequests, util.ErrorResponse{
				Error: util.NewAPIError(http.StatusTooManyRequests, "rate limit exceeded"),
			})
			return
		}

		// Continue with the request
		c.Next()
	}
}
