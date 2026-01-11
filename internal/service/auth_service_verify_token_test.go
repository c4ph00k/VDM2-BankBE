package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"

	"VDM2-BankBE/internal/config"
	"VDM2-BankBE/internal/model"
	repmocks "VDM2-BankBE/internal/repository/mocks"
	"VDM2-BankBE/internal/service"
	"VDM2-BankBE/internal/testutil"
	"VDM2-BankBE/internal/util"
)

func TestAuthService_VerifyToken_JWTAndPASETO(t *testing.T) {
	t.Parallel()

	jwtSecret := "test-jwt-secret"
	pasetoSecret := "test-paseto-secret"
	userID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	user := &model.User{ID: userID}

	tests := []struct {
		name          string
		cfg           *config.Config
		token         func(t *testing.T) string
		buildMocks    func(ctrl *gomock.Controller) *repmocks.MockUserRepository
		assertOutcome func(t *testing.T, got *model.User, err error)
	}{
		{
			name: "valid JWT returns user",
			cfg: &config.Config{
				JWT:    config.JWTConfig{Secret: jwtSecret, Expiry: time.Hour},
				PASETO: config.PASETOConfig{Secret: pasetoSecret},
			},
			token: func(t *testing.T) string {
				return testutil.ValidJWT(t, jwtSecret, userID.String())
			},
			buildMocks: func(ctrl *gomock.Controller) *repmocks.MockUserRepository {
				m := repmocks.NewMockUserRepository(ctrl)
				m.EXPECT().GetByID(gomock.Any(), userID).Return(user, nil)
				return m
			},
			assertOutcome: func(t *testing.T, got *model.User, err error) {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if got == nil || got.ID != userID {
					t.Fatalf("unexpected user: %+v", got)
				}
			},
		},
		{
			// Current behavior: jwt.ParseWithClaims validates exp and returns an error for expired tokens,
			// which the service maps to "invalid token" before reaching the explicit exp check.
			name: "expired JWT returns unauthorized invalid token",
			cfg: &config.Config{
				JWT:    config.JWTConfig{Secret: jwtSecret, Expiry: time.Hour},
				PASETO: config.PASETOConfig{Secret: pasetoSecret},
			},
			token: func(t *testing.T) string {
				return testutil.ExpiredJWT(t, jwtSecret, userID.String())
			},
			buildMocks: func(ctrl *gomock.Controller) *repmocks.MockUserRepository {
				return repmocks.NewMockUserRepository(ctrl)
			},
			assertOutcome: func(t *testing.T, got *model.User, err error) {
				if got != nil {
					t.Fatalf("expected nil user, got %+v", got)
				}
				apiErr, ok := err.(*util.APIError)
				if !ok || apiErr.Code != 401 || apiErr.Message != "invalid token" {
					t.Fatalf("unexpected error: %#v", err)
				}
			},
		},
		{
			name: "invalid JWT returns unauthorized invalid token",
			cfg: &config.Config{
				JWT:    config.JWTConfig{Secret: jwtSecret, Expiry: time.Hour},
				PASETO: config.PASETOConfig{Secret: pasetoSecret},
			},
			token: func(t *testing.T) string {
				return "not-a-jwt"
			},
			buildMocks: func(ctrl *gomock.Controller) *repmocks.MockUserRepository {
				return repmocks.NewMockUserRepository(ctrl)
			},
			assertOutcome: func(t *testing.T, got *model.User, err error) {
				if got != nil {
					t.Fatalf("expected nil user, got %+v", got)
				}
				apiErr, ok := err.(*util.APIError)
				if !ok || apiErr.Code != 401 || apiErr.Message != "invalid token" {
					t.Fatalf("unexpected error: %#v", err)
				}
			},
		},
		{
			name: "JWT with invalid user ID returns unauthorized invalid user ID in token",
			cfg: &config.Config{
				JWT:    config.JWTConfig{Secret: jwtSecret, Expiry: time.Hour},
				PASETO: config.PASETOConfig{Secret: pasetoSecret},
			},
			token: func(t *testing.T) string {
				claims := jwt.MapClaims{
					"user_id": "not-a-uuid",
					"exp":     time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				s, err := token.SignedString([]byte(jwtSecret))
				if err != nil {
					t.Fatalf("failed to sign JWT: %v", err)
				}
				return s
			},
			buildMocks: func(ctrl *gomock.Controller) *repmocks.MockUserRepository {
				return repmocks.NewMockUserRepository(ctrl)
			},
			assertOutcome: func(t *testing.T, got *model.User, err error) {
				if got != nil {
					t.Fatalf("expected nil user, got %+v", got)
				}
				apiErr, ok := err.(*util.APIError)
				if !ok || apiErr.Code != 401 || apiErr.Message != "invalid user ID in token" {
					t.Fatalf("unexpected error: %#v", err)
				}
			},
		},
		{
			name: "valid PASETO returns user",
			cfg: &config.Config{
				JWT:    config.JWTConfig{Secret: jwtSecret, Expiry: time.Hour},
				PASETO: config.PASETOConfig{Secret: pasetoSecret},
			},
			token: func(t *testing.T) string {
				return testutil.ValidPASETO(t, pasetoSecret, userID.String())
			},
			buildMocks: func(ctrl *gomock.Controller) *repmocks.MockUserRepository {
				m := repmocks.NewMockUserRepository(ctrl)
				m.EXPECT().GetByID(gomock.Any(), userID).Return(user, nil)
				return m
			},
			assertOutcome: func(t *testing.T, got *model.User, err error) {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if got == nil || got.ID != userID {
					t.Fatalf("unexpected user: %+v", got)
				}
			},
		},
		{
			name: "expired PASETO returns unauthorized token expired",
			cfg: &config.Config{
				JWT:    config.JWTConfig{Secret: jwtSecret, Expiry: time.Hour},
				PASETO: config.PASETOConfig{Secret: pasetoSecret},
			},
			token: func(t *testing.T) string {
				return testutil.ExpiredPASETO(t, pasetoSecret, userID.String())
			},
			buildMocks: func(ctrl *gomock.Controller) *repmocks.MockUserRepository {
				return repmocks.NewMockUserRepository(ctrl)
			},
			assertOutcome: func(t *testing.T, got *model.User, err error) {
				if got != nil {
					t.Fatalf("expected nil user, got %+v", got)
				}
				apiErr, ok := err.(*util.APIError)
				if !ok || apiErr.Code != 401 || apiErr.Message != "token expired" {
					t.Fatalf("unexpected error: %#v", err)
				}
			},
		},
		{
			name: "PASETO key falls back to JWT secret when PASETO secret is empty",
			cfg: &config.Config{
				JWT:    config.JWTConfig{Secret: jwtSecret, Expiry: time.Hour},
				PASETO: config.PASETOConfig{Secret: ""},
			},
			token: func(t *testing.T) string {
				// AuthService derives key from JWT secret when PASETO secret is empty.
				return testutil.ValidPASETO(t, jwtSecret, userID.String())
			},
			buildMocks: func(ctrl *gomock.Controller) *repmocks.MockUserRepository {
				m := repmocks.NewMockUserRepository(ctrl)
				m.EXPECT().GetByID(gomock.Any(), userID).Return(user, nil)
				return m
			},
			assertOutcome: func(t *testing.T, got *model.User, err error) {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if got == nil || got.ID != userID {
					t.Fatalf("unexpected user: %+v", got)
				}
			},
		},
		{
			name: "unsupported PASETO version returns unauthorized",
			cfg: &config.Config{
				JWT:    config.JWTConfig{Secret: jwtSecret, Expiry: time.Hour},
				PASETO: config.PASETOConfig{Secret: pasetoSecret},
			},
			token: func(t *testing.T) string {
				return "v4.local.this-is-not-supported"
			},
			buildMocks: func(ctrl *gomock.Controller) *repmocks.MockUserRepository {
				return repmocks.NewMockUserRepository(ctrl)
			},
			assertOutcome: func(t *testing.T, got *model.User, err error) {
				if got != nil {
					t.Fatalf("expected nil user, got %+v", got)
				}
				apiErr, ok := err.(*util.APIError)
				if !ok || apiErr.Code != 401 || apiErr.Message != "unsupported paseto version" {
					t.Fatalf("unexpected error: %#v", err)
				}
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userRepo := tc.buildMocks(ctrl)
			svc := service.NewAuthService(
				userRepo,
				nil,
				nil,
				nil,
				nil,
				tc.cfg,
			)

			got, err := svc.VerifyToken(context.Background(), tc.token(t))
			tc.assertOutcome(t, got, err)
		})
	}
}

