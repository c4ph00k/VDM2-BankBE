package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"VDM2-BankBE/internal/generated"
	"VDM2-BankBE/internal/handler"
)

// Server delegates generated OpenAPI handlers to the existing handwritten handlers.
// This is the bridge that makes "contract = reality" enforceable at runtime.
type Server struct {
	Auth     *handler.AuthHandler
	Account  *handler.AccountHandler
	Movement *handler.MovementHandler
	Transfer *handler.TransferHandler
}

var _ generated.ServerInterface = (*Server)(nil)

func NewServer(
	auth *handler.AuthHandler,
	account *handler.AccountHandler,
	movement *handler.MovementHandler,
	transfer *handler.TransferHandler,
) *Server {
	return &Server{
		Auth:     auth,
		Account:  account,
		Movement: movement,
		Transfer: transfer,
	}
}

func (s *Server) AuthSignUp(c *gin.Context)                  { s.Auth.SignUp(c) }
func (s *Server) AuthLogin(c *gin.Context)                   { s.Auth.Login(c) }
func (s *Server) AuthGoogle(c *gin.Context)                  { s.Auth.GoogleAuth(c) }
func (s *Server) AuthGoogleCallback(c *gin.Context, _ generated.AuthGoogleCallbackParams) {
	// Existing handler reads query params directly.
	s.Auth.GoogleCallback(c)
}

func (s *Server) AccountsGetBalance(c *gin.Context) { s.Account.Balance(c) }

func (s *Server) AccountsListMovements(c *gin.Context, _ generated.AccountsListMovementsParams) {
	// Existing handler reads query params directly.
	s.Movement.List(c)
}

func (s *Server) AccountsCreateMovement(c *gin.Context) { s.Movement.Create(c) }

func (s *Server) TransfersList(c *gin.Context, _ generated.TransfersListParams) {
	// Existing handler reads query params directly.
	s.Transfer.List(c)
}

func (s *Server) TransfersCreate(c *gin.Context) { s.Transfer.Transfer(c) }

func (s *Server) HealthCheck(c *gin.Context) { c.Status(http.StatusOK) }

func (s *Server) Metrics(c *gin.Context) {
	// Match existing behavior: serve Prometheus metrics handler.
	gin.WrapH(promhttp.Handler())(c)
}

// Swagger UI is intentionally kept OUTSIDE the generated server because Gin swagger uses a wildcard route (`/swagger/*any`)
// which cannot be represented in OpenAPI paths. It is still registered by the router at runtime.
func RegisterSwaggerRoutes(router gin.IRouter) {
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

