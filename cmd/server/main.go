package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"VDM2-BankBE/internal/generated"
)

// Server is a minimal oapi-codegen Gin server implementation.
//
// IMPORTANT:
// - This file is intentionally minimal wiring for an OpenAPI-first future.
// - Your current production server wiring is in `cmd/api/main.go` and uses the handwritten router.
// - After the baseline spec is extracted, the intended direction is to implement
//   the generated `generated.ServerInterface` by delegating to services/handlers.
type Server struct{}

// authMiddlewarePlaceholder documents where Bearer auth should be enforced.
// Current production enforcement is done via `internal/middleware/auth_middleware.go`.
func authMiddlewarePlaceholder() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: validate Authorization: Bearer <JWT-or-PASETO> and set user in context
		c.Next()
	}
}

func main() {
	r := gin.New()
	r.Use(gin.Recovery())

	// Example: plug in auth middleware if you switch to generated routing.
	_ = authMiddlewarePlaceholder

	// Register generated handlers. This will mount paths exactly as defined in OpenAPI.
	// TODO: Replace `&Server{}` with a real implementation.
	generated.RegisterHandlers(r, &Server{})

	_ = r.Run(":8080")
}

// NOTE: The following methods are stubs to keep the build reproducible.
// Implement them by delegating to your real services/handlers.

func (s *Server) AuthSignUp(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

func (s *Server) AuthLogin(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

func (s *Server) AuthGoogle(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

func (s *Server) AuthGoogleCallback(c *gin.Context, params generated.AuthGoogleCallbackParams) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

func (s *Server) AccountsGetBalance(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

func (s *Server) AccountsListMovements(c *gin.Context, params generated.AccountsListMovementsParams) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

func (s *Server) AccountsCreateMovement(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

func (s *Server) TransfersList(c *gin.Context, params generated.TransfersListParams) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

func (s *Server) TransfersCreate(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

func (s *Server) HealthCheck(c *gin.Context) {
	c.Status(http.StatusOK)
}

func (s *Server) Metrics(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

