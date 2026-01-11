package testutil

import (
	"testing"

	"github.com/gin-gonic/gin"

	"VDM2-BankBE/internal/api"
	"VDM2-BankBE/internal/generated"
	"VDM2-BankBE/internal/handler"
	"VDM2-BankBE/internal/middleware"
	"VDM2-BankBE/internal/util"
)

type RouterDeps struct {
	AuthHandler     *handler.AuthHandler
	AccountHandler  *handler.AccountHandler
	MovementHandler *handler.MovementHandler
	TransferHandler *handler.TransferHandler

	AuthMiddleware      *middleware.AuthMiddleware
	RateLimitMiddleware *middleware.RateLimitMiddleware
}

func SetupGinRouter(t *testing.T, deps RouterDeps) *gin.Engine {
	t.Helper()

	gin.SetMode(gin.TestMode)
	r := gin.New()

	// Mirror production error envelope for OpenAPI parameter binding errors.
	errorHandler := func(c *gin.Context, err error, statusCode int) {
		c.AbortWithStatusJSON(statusCode, util.ErrorResponse{
			Error: util.NewAPIError(statusCode, err.Error()),
		})
	}

	server := api.NewServer(deps.AuthHandler, deps.AccountHandler, deps.MovementHandler, deps.TransferHandler)

	var mws []generated.MiddlewareFunc
	if deps.AuthMiddleware != nil {
		mws = append(mws, deps.AuthMiddleware.AuthenticateIfRequiredFunc())
	}
	if deps.RateLimitMiddleware != nil {
		mws = append(mws, deps.RateLimitMiddleware.LimitByUserFunc())
	}

	generated.RegisterHandlersWithOptions(r, server, generated.GinServerOptions{
		Middlewares:  mws,
		ErrorHandler: errorHandler,
	})

	return r
}

