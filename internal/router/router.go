package router

import (
	"net/http"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"VDM2-BankBE/internal/api"
	"VDM2-BankBE/internal/generated"
	"VDM2-BankBE/internal/handler"
	"VDM2-BankBE/internal/middleware"
	"VDM2-BankBE/internal/util"
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

	// Register Swagger documentation (wildcard route)
	// NOTE: Swagger is intentionally kept outside the generated server because OpenAPI cannot represent `/swagger/*any`.
	api.RegisterSwaggerRoutes(r.engine)

	// Build the generated-server adapter that delegates to existing handlers.
	server := api.NewServer(r.authHandler, r.accountHandler, r.movementHandler, r.transferHandler)

	// Register OpenAPI-generated routes with per-operation middlewares.
	// These middlewares run AFTER the generated wrapper sets operation security markers.
	generated.RegisterHandlersWithOptions(r.engine, server, generated.GinServerOptions{
		Middlewares: []generated.MiddlewareFunc{
			// Enforce JWT/PASETO only where the OpenAPI contract requires it.
			r.authMiddleware.AuthenticateIfRequiredFunc(),
			// Apply per-user rate limiting when a user is present.
			r.rateLimitMiddleware.LimitByUserFunc(),
		},
		// Preserve the existing error envelope format for parameter binding errors.
		ErrorHandler: func(c *gin.Context, err error, statusCode int) {
			c.AbortWithStatusJSON(statusCode, util.ErrorResponse{
				Error: util.NewAPIError(statusCode, err.Error()),
			})
		},
	})

	return r.engine
}
