package router

import (
	"net/http"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"

	"VDM2-BankBE/internal/handler"
	"VDM2-BankBE/internal/middleware"
)

// Router handles HTTP routing with Gin
type Router struct {
	engine              *gin.Engine
	authHandler         *handler.AuthHandler
	accountHandler      *handler.AccountHandler
	movementHandler     *handler.MovementHandler
	transferHandler     *handler.TransferHandler
	authMiddleware      *middleware.AuthMiddleware
	rateLimitMiddleware *middleware.RateLimitMiddleware
	logger              *zap.Logger
}

// NewRouter creates a new router
func NewRouter(
	authHandler *handler.AuthHandler,
	accountHandler *handler.AccountHandler,
	movementHandler *handler.MovementHandler,
	transferHandler *handler.TransferHandler,
	authMiddleware *middleware.AuthMiddleware,
	rateLimitMiddleware *middleware.RateLimitMiddleware,
	logger *zap.Logger,
) *Router {
	return &Router{
		engine:              gin.New(),
		authHandler:         authHandler,
		accountHandler:      accountHandler,
		movementHandler:     movementHandler,
		transferHandler:     transferHandler,
		authMiddleware:      authMiddleware,
		rateLimitMiddleware: rateLimitMiddleware,
		logger:              logger,
	}
}

// Setup sets up the routes with Gin
func (r *Router) Setup() http.Handler {
	// Use Gin's zap logger and recovery middleware
	r.engine.Use(ginzap.Ginzap(r.logger, time.RFC3339, true))
	r.engine.Use(middleware.NewLoggingMiddleware(r.logger).LogRequest())
	r.engine.Use(ginzap.RecoveryWithZap(r.logger, true))

	// Setup API versioning
	v1 := r.engine.Group("/api/v1")

	// Public routes
	auth := v1.Group("/auth")
	{
		auth.POST("/signup", r.authHandler.SignUp)
		auth.POST("/login", r.authHandler.Login)
		auth.GET("/google", r.authHandler.GoogleAuth)
		auth.GET("/google/callback", r.authHandler.GoogleCallback)
	}

	// Health check and metrics
	r.engine.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	r.engine.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Swagger documentation
	r.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Protected routes
	protected := v1.Group("")
	protected.Use(r.authMiddleware.Authenticate())
	protected.Use(r.rateLimitMiddleware.LimitByUser())

	// Account routes
	accounts := protected.Group("/accounts")
	{
		accounts.GET("/balance", r.accountHandler.Balance)
		accounts.GET("/movements", r.movementHandler.List)
		accounts.POST("/movements", r.movementHandler.Create)
	}

	// Transfer routes
	transfers := protected.Group("/transfers")
	{
		transfers.POST("", r.transferHandler.Transfer)
		transfers.GET("", r.transferHandler.List)
	}

	return r.engine
}
