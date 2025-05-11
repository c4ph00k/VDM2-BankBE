package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"golang.org/x/crypto/bcrypt"

	"VDM2-BankBE/internal/config"
	"VDM2-BankBE/internal/model"
	"VDM2-BankBE/internal/repository"
	"VDM2-BankBE/internal/util"
	"VDM2-BankBE/pkg/cache"
	"VDM2-BankBE/pkg/oauth"
)

// JWTClaims represents the claims in the JWT
type JWTClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// DefaultAuthService implements AuthService
type DefaultAuthService struct {
	userRepo       repository.UserRepository
	accountRepo    repository.AccountRepository
	oauthTokenRepo repository.OAuthTokenRepository
	redisClient    *cache.RedisClient
	googleOAuth    *oauth.GoogleOAuthClient
	config         *config.Config
}

// NewAuthService creates a new authentication service
func NewAuthService(
	userRepo repository.UserRepository,
	accountRepo repository.AccountRepository,
	oauthTokenRepo repository.OAuthTokenRepository,
	redisClient *cache.RedisClient,
	googleOAuth *oauth.GoogleOAuthClient,
	config *config.Config,
) AuthService {
	return &DefaultAuthService{
		userRepo:       userRepo,
		accountRepo:    accountRepo,
		oauthTokenRepo: oauthTokenRepo,
		redisClient:    redisClient,
		googleOAuth:    googleOAuth,
		config:         config,
	}
}

// SignUp registers a new user
func (s *DefaultAuthService) SignUp(
	ctx context.Context,
	email, username, firstName, lastName, fiscalCode, password string,
) (*model.User, error) {
	// Check if email is already taken
	_, err := s.userRepo.GetByEmail(ctx, email)
	if err == nil {
		return nil, util.NewBadRequestError("email already in use")
	} else if _, ok := err.(*util.APIError); !ok {
		return nil, errors.Wrap(err, "failed to check email")
	}

	// Check if username is already taken
	_, err = s.userRepo.GetByUsername(ctx, username)
	if err == nil {
		return nil, util.NewBadRequestError("username already in use")
	} else if _, ok := err.(*util.APIError); !ok {
		return nil, errors.Wrap(err, "failed to check username")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Wrap(err, "failed to hash password")
	}

	// Create user
	user := &model.User{
		ID:           uuid.New(),
		Email:        email,
		Username:     username,
		FirstName:    firstName,
		LastName:     lastName,
		FiscalCode:   fiscalCode,
		PasswordHash: string(hashedPassword),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Save user to DB
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, errors.Wrap(err, "failed to create user")
	}

	// Create an account for the user
	account := &model.Account{
		UserID:    user.ID,
		Balance:   decimal.Zero,
		Currency:  "EUR",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := s.accountRepo.Create(ctx, account); err != nil {
		return nil, errors.Wrap(err, "failed to create account")
	}

	return user, nil
}

// Login authenticates a user and returns a JWT token
func (s *DefaultAuthService) Login(ctx context.Context, email, password string) (string, error) {
	// Find user by email
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if _, ok := err.(*util.APIError); ok {
			return "", util.NewUnauthorizedError("invalid email or password")
		}
		return "", errors.Wrap(err, "failed to get user by email")
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", util.NewUnauthorizedError("invalid email or password")
	}

	// Generate JWT token
	token, err := s.generateJWT(user.ID.String())
	if err != nil {
		return "", errors.Wrap(err, "failed to generate JWT")
	}

	return token, nil
}

// GoogleAuth starts the Google OAuth flow
func (s *DefaultAuthService) GoogleAuth(ctx context.Context) (string, string, error) {
	// Generate a random state for CSRF protection
	stateBytes := make([]byte, 16)
	if _, err := rand.Read(stateBytes); err != nil {
		return "", "", errors.Wrap(err, "failed to generate random state")
	}
	state := hex.EncodeToString(stateBytes)

	// Store the state in Redis for later verification
	err := s.redisClient.SetOAuthState(ctx, state, "")
	if err != nil {
		return "", "", errors.Wrap(err, "failed to store OAuth state")
	}

	// Get the authorization URL
	authURL := s.googleOAuth.GetAuthURL(state)

	return authURL, state, nil
}

// GoogleCallback handles the Google OAuth callback
func (s *DefaultAuthService) GoogleCallback(ctx context.Context, code, state string) (string, error) {
	// Verify the state to prevent CSRF
	_, err := s.redisClient.GetOAuthState(ctx, state)
	if err != nil {
		return "", util.NewBadRequestError("invalid or expired OAuth state")
	}

	// Exchange the code for tokens
	token, err := s.googleOAuth.Exchange(ctx, code)
	if err != nil {
		return "", errors.Wrap(err, "failed to exchange OAuth code")
	}

	// Get user info from Google
	userInfo, err := s.googleOAuth.GetUserInfo(ctx, token)
	if err != nil {
		return "", errors.Wrap(err, "failed to get Google user info")
	}

	// Look for a user with this email
	user, err := s.userRepo.GetByEmail(ctx, userInfo.Email)
	if err != nil {
		// If the user doesn't exist, create one
		if _, ok := err.(*util.APIError); ok {
			// Split name into first and last name
			names := strings.Split(userInfo.Name, " ")
			firstName := userInfo.GivenName
			if firstName == "" && len(names) > 0 {
				firstName = names[0]
			}
			lastName := userInfo.FamilyName
			if lastName == "" && len(names) > 1 {
				lastName = names[len(names)-1]
			}

			// Generate a username based on email
			username := strings.Split(userInfo.Email, "@")[0]

			// Generate a password (user won't actually use this)
			passwordBytes := make([]byte, 32)
			if _, err := rand.Read(passwordBytes); err != nil {
				return "", errors.Wrap(err, "failed to generate random password")
			}
			password := hex.EncodeToString(passwordBytes)

			// Create the user
			user, err = s.SignUp(ctx, userInfo.Email, username, firstName, lastName, userInfo.ID, password)
			if err != nil {
				return "", errors.Wrap(err, "failed to create user from Google account")
			}
		} else {
			return "", errors.Wrap(err, "failed to check for existing user")
		}
	}

	// Store or update OAuth tokens
	oauthToken := &model.OAuthToken{
		UserID:       user.ID,
		Provider:     "google",
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresAt:    token.Expiry,
	}

	// Check if token already exists
	existingToken, err := s.oauthTokenRepo.GetByUserID(ctx, user.ID)
	if err != nil {
		if _, ok := err.(*util.APIError); ok {
			// Create new token
			if err := s.oauthTokenRepo.Create(ctx, oauthToken); err != nil {
				return "", errors.Wrap(err, "failed to store OAuth token")
			}
		} else {
			return "", errors.Wrap(err, "failed to check existing OAuth token")
		}
	} else {
		// Update existing token
		existingToken.AccessToken = token.AccessToken
		existingToken.RefreshToken = token.RefreshToken
		existingToken.ExpiresAt = token.Expiry
		if err := s.oauthTokenRepo.Update(ctx, existingToken); err != nil {
			return "", errors.Wrap(err, "failed to update OAuth token")
		}
	}

	// Generate JWT token
	jwtToken, err := s.generateJWT(user.ID.String())
	if err != nil {
		return "", errors.Wrap(err, "failed to generate JWT")
	}

	return jwtToken, nil
}

// VerifyToken verifies a JWT token and returns the user
func (s *DefaultAuthService) VerifyToken(ctx context.Context, tokenString string) (*model.User, error) {
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Make sure the token method is what we expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, util.NewUnauthorizedError("unexpected token signing method")
		}
		return []byte(s.config.JWT.Secret), nil
	})

	if err != nil {
		return nil, util.NewUnauthorizedError("invalid token")
	}

	// Get the claims
	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, util.NewUnauthorizedError("invalid token")
	}

	// Check if the token is expired
	expirationTime, err := claims.GetExpirationTime()
	if err != nil || expirationTime == nil || expirationTime.Before(time.Now()) {
		return nil, util.NewUnauthorizedError("token expired")
	}

	// Convert user ID from string to UUID
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, util.NewUnauthorizedError("invalid user ID in token")
	}

	// Get the user from the database
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if _, ok := err.(*util.APIError); ok {
			return nil, util.NewUnauthorizedError("user not found")
		}
		return nil, errors.Wrap(err, "failed to get user from token")
	}

	return user, nil
}

// generateJWT generates a JWT token for a user
func (s *DefaultAuthService) generateJWT(userID string) (string, error) {
	// Determine token expiry
	expiresAt := time.Now().Add(s.config.JWT.Expiry)

	// Create the claims
	claims := &JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "VDM2-Bank",
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token
	tokenString, err := token.SignedString([]byte(s.config.JWT.Secret))
	if err != nil {
		return "", errors.Wrap(err, "failed to sign JWT")
	}

	return tokenString, nil
}
