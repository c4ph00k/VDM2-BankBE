package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v9"

	"VDM2-BankBE/internal/service"
	"VDM2-BankBE/internal/util"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	authService service.AuthService
	validator   *validator.Validate
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validator:   validator.New(),
	}
}

// SignUpRequest represents a sign-up request
type SignUpRequest struct {
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required,min=8"`
	Username   string `json:"username" validate:"required,min=3"`
	FirstName  string `json:"first_name" validate:"required"`
	LastName   string `json:"last_name" validate:"required"`
	FiscalCode string `json:"fiscal_code" validate:"required"`
}

// LoginRequest represents a login request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// AuthResponse represents an authentication response
type AuthResponse struct {
	Token     string `json:"token"`
	ExpiresIn int    `json:"expires_in"`
}

// SignUp handles user registration
// @Summary Register a new user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param user body SignUpRequest true "User details"
// @Success 201 {object} model.User
// @Failure 400 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /auth/signup [post]
func (h *AuthHandler) SignUp(c *gin.Context) {
	// Parse and validate request
	var req SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.ErrorResponse{
			Error: util.NewBadRequestError("invalid request body"),
		})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, util.ErrorResponse{
			Error: util.NewBadRequestError(err.Error()),
		})
		return
	}

	// Call service
	user, err := h.authService.SignUp(
		c,
		req.Email,
		req.Username,
		req.FirstName,
		req.LastName,
		req.FiscalCode,
		req.Password,
	)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	// Return response
	c.JSON(http.StatusCreated, user)
}

// Login handles user login
// @Summary Login a user
// @Description Authenticate user and return a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "Login credentials"
// @Success 200 {object} AuthResponse
// @Failure 400 {object} util.ErrorResponse
// @Failure 401 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	// Parse and validate request
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.ErrorResponse{
			Error: util.NewBadRequestError("invalid request body"),
		})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, util.ErrorResponse{
			Error: util.NewBadRequestError(err.Error()),
		})
		return
	}

	// Call service
	token, err := h.authService.Login(c, req.Email, req.Password)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	// Return response
	response := AuthResponse{
		Token:     token,
		ExpiresIn: 3600, // 1 hour
	}

	c.JSON(http.StatusOK, response)
}

// GoogleAuth initiates Google OAuth flow
// @Summary Start Google OAuth flow
// @Description Redirect user to Google for authentication
// @Tags auth
// @Produce json
// @Success 307 {string} string "Redirects to Google"
// @Failure 500 {object} util.ErrorResponse
// @Router /auth/google [get]
func (h *AuthHandler) GoogleAuth(c *gin.Context) {
	// Get authentication URL
	url, state, err := h.authService.GoogleAuth(c)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	// Store state in cookie for validation in callback
	c.SetCookie("oauth_state", state, 3600, "/", "", false, true)

	// Redirect to Google
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// GoogleCallback handles the Google OAuth callback
// @Summary Handle Google OAuth callback
// @Description Process OAuth callback and authenticate user
// @Tags auth
// @Produce html
// @Param code query string true "OAuth code"
// @Param state query string true "CSRF state"
// @Success 200 {string} string "HTML with token"
// @Failure 400 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /auth/google/callback [get]
func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	// Get code and state from query
	code := c.Query("code")
	state := c.Query("state")

	if code == "" || state == "" {
		c.JSON(http.StatusBadRequest, util.ErrorResponse{
			Error: util.NewBadRequestError("missing code or state"),
		})
		return
	}

	// Handle callback
	token, err := h.authService.GoogleCallback(c, code, state)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	// Return HTML with token (in a real app, redirect to frontend with token)
	html := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>Authentication Successful</title>
		<script>
			// Store token in localStorage
			localStorage.setItem("auth_token", "` + token + `");
			// Redirect to home page
			window.location.href = "/";
		</script>
	</head>
	<body>
		<h1>Authentication Successful</h1>
		<p>You will be redirected shortly...</p>
	</body>
	</html>
	`

	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, html)
}
