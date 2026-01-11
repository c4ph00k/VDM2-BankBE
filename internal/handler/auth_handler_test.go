package handler_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"VDM2-BankBE/internal/config"
	"VDM2-BankBE/internal/handler"
	"VDM2-BankBE/internal/middleware"
	"VDM2-BankBE/internal/model"
	servicemocks "VDM2-BankBE/internal/service/mocks"
	"VDM2-BankBE/internal/testutil"
	"VDM2-BankBE/internal/util"
)

func TestAuth_Login(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		requestBody    any
		buildMocks     func(ctrl *gomock.Controller) *servicemocks.MockAuthService
		expectedStatus int
		assertResponse func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name:        "success",
			requestBody: map[string]any{"email": "a@example.com", "password": "pass"},
			buildMocks: func(ctrl *gomock.Controller) *servicemocks.MockAuthService {
				m := servicemocks.NewMockAuthService(ctrl)
				m.EXPECT().
					Login(gomock.Any(), "a@example.com", "pass").
					Return("token-123", nil)
				return m
			},
			expectedStatus: http.StatusOK,
			assertResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				testutil.AssertHTTPStatus(t, rec, http.StatusOK)
				var got handler.AuthResponse
				got = testutil.DecodeJSONResponse[handler.AuthResponse](t, rec)
				if got.Token != "token-123" {
					t.Fatalf("unexpected token: got=%q want=%q", got.Token, "token-123")
				}
				if got.ExpiresIn != 3600 {
					t.Fatalf("unexpected expires_in: got=%d want=%d", got.ExpiresIn, 3600)
				}
			},
		},
		{
			name:        "bad request body",
			requestBody: nil,
			buildMocks: func(ctrl *gomock.Controller) *servicemocks.MockAuthService {
				return servicemocks.NewMockAuthService(ctrl)
			},
			expectedStatus: http.StatusBadRequest,
			assertResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				testutil.AssertHTTPError(t, rec, http.StatusBadRequest, "invalid request body")
			},
		},
		{
			name:        "validation error",
			requestBody: map[string]any{"email": "not-an-email", "password": "pass"},
			buildMocks: func(ctrl *gomock.Controller) *servicemocks.MockAuthService {
				return servicemocks.NewMockAuthService(ctrl)
			},
			expectedStatus: http.StatusBadRequest,
			assertResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				testutil.AssertHTTPError(t, rec, http.StatusBadRequest, http.StatusBadRequest)
			},
		},
		{
			name:        "unauthorized mapped to 401",
			requestBody: map[string]any{"email": "a@example.com", "password": "wrong"},
			buildMocks: func(ctrl *gomock.Controller) *servicemocks.MockAuthService {
				m := servicemocks.NewMockAuthService(ctrl)
				m.EXPECT().
					Login(gomock.Any(), "a@example.com", "wrong").
					Return("", util.NewUnauthorizedError("invalid email or password"))
				return m
			},
			expectedStatus: http.StatusUnauthorized,
			assertResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				testutil.AssertHTTPError(t, rec, http.StatusUnauthorized, "invalid email or password")
			},
		},
		{
			name:        "generic error mapped to 500",
			requestBody: map[string]any{"email": "a@example.com", "password": "pass"},
			buildMocks: func(ctrl *gomock.Controller) *servicemocks.MockAuthService {
				m := servicemocks.NewMockAuthService(ctrl)
				m.EXPECT().
					Login(gomock.Any(), "a@example.com", "pass").
					Return("", errors.New("boom"))
				return m
			},
			expectedStatus: http.StatusInternalServerError,
			assertResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				testutil.AssertHTTPError(t, rec, http.StatusInternalServerError, "internal server error")
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			authSvc := tc.buildMocks(ctrl)
			r := newTestRouter(t, ctrl, authSvc, nil, nil, nil)

			rec := httptest.NewRecorder()
			req := testutil.NewJSONRequest(http.MethodPost, "/api/v1/auth/login", tc.requestBody, nil)
			r.ServeHTTP(rec, req)

			tc.assertResponse(t, rec)
		})
	}
}

func TestAuth_SignUp(t *testing.T) {
	t.Parallel()

	fixedTime := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	userID := uuid.MustParse("00000000-0000-0000-0000-000000000002")
	returnedUser := &model.User{
		ID:         userID,
		Email:      "a@example.com",
		Username:   "alice",
		FirstName:  "Alice",
		LastName:   "A",
		FiscalCode: "FC1",
		CreatedAt:  fixedTime,
		UpdatedAt:  fixedTime,
	}

	tests := []struct {
		name           string
		requestBody    any
		buildMocks     func(ctrl *gomock.Controller) *servicemocks.MockAuthService
		expectedStatus int
		assertResponse func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name: "success",
			requestBody: map[string]any{
				"email":       "a@example.com",
				"password":    "SecurePass123!",
				"username":    "alice",
				"first_name":  "Alice",
				"last_name":   "A",
				"fiscal_code": "FC1",
			},
			buildMocks: func(ctrl *gomock.Controller) *servicemocks.MockAuthService {
				m := servicemocks.NewMockAuthService(ctrl)
				m.EXPECT().
					SignUp(gomock.Any(), "a@example.com", "alice", "Alice", "A", "FC1", "SecurePass123!").
					Return(returnedUser, nil)
				return m
			},
			expectedStatus: http.StatusCreated,
			assertResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				testutil.AssertHTTPStatus(t, rec, http.StatusCreated)
				got := testutil.DecodeJSONResponse[model.User](t, rec)
				if got.ID != userID {
					t.Fatalf("unexpected user id: got=%s want=%s", got.ID.String(), userID.String())
				}
				if got.Email != "a@example.com" {
					t.Fatalf("unexpected email: got=%q", got.Email)
				}
			},
		},
		{
			name:        "bad request body",
			requestBody: nil,
			buildMocks: func(ctrl *gomock.Controller) *servicemocks.MockAuthService {
				return servicemocks.NewMockAuthService(ctrl)
			},
			expectedStatus: http.StatusBadRequest,
			assertResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				testutil.AssertHTTPError(t, rec, http.StatusBadRequest, "invalid request body")
			},
		},
		{
			name: "validation error",
			requestBody: map[string]any{
				"email":       "not-an-email",
				"password":    "short",
				"username":    "al",
				"first_name":  "",
				"last_name":   "",
				"fiscal_code": "",
			},
			buildMocks: func(ctrl *gomock.Controller) *servicemocks.MockAuthService {
				return servicemocks.NewMockAuthService(ctrl)
			},
			expectedStatus: http.StatusBadRequest,
			assertResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				testutil.AssertHTTPError(t, rec, http.StatusBadRequest, http.StatusBadRequest)
			},
		},
		{
			name: "service error mapped via util.HandleError",
			requestBody: map[string]any{
				"email":       "a@example.com",
				"password":    "SecurePass123!",
				"username":    "alice",
				"first_name":  "Alice",
				"last_name":   "A",
				"fiscal_code": "FC1",
			},
			buildMocks: func(ctrl *gomock.Controller) *servicemocks.MockAuthService {
				m := servicemocks.NewMockAuthService(ctrl)
				m.EXPECT().
					SignUp(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, util.NewBadRequestError("email already in use"))
				return m
			},
			expectedStatus: http.StatusBadRequest,
			assertResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				testutil.AssertHTTPError(t, rec, http.StatusBadRequest, "email already in use")
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			authSvc := tc.buildMocks(ctrl)
			r := newTestRouter(t, ctrl, authSvc, nil, nil, nil)

			rec := httptest.NewRecorder()
			req := testutil.NewJSONRequest(http.MethodPost, "/api/v1/auth/signup", tc.requestBody, nil)
			r.ServeHTTP(rec, req)

			tc.assertResponse(t, rec)
		})
	}
}

func newTestRouter(
	t *testing.T,
	ctrl *gomock.Controller,
	authSvc *servicemocks.MockAuthService,
	accountSvc *servicemocks.MockAccountService,
	movementSvc *servicemocks.MockMovementService,
	transferSvc *servicemocks.MockTransferService,
) http.Handler {
	t.Helper()

	// Default no-op mocks for unused services in a given test.
	if authSvc == nil {
		authSvc = servicemocks.NewMockAuthService(ctrl)
	}
	if accountSvc == nil {
		accountSvc = servicemocks.NewMockAccountService(ctrl)
	}
	if movementSvc == nil {
		movementSvc = servicemocks.NewMockMovementService(ctrl)
	}
	if transferSvc == nil {
		transferSvc = servicemocks.NewMockTransferService(ctrl)
	}

	authHandler := handler.NewAuthHandler(authSvc)
	accountHandler := handler.NewAccountHandler(accountSvc)
	movementHandler := handler.NewMovementHandler(movementSvc, accountSvc)
	transferHandler := handler.NewTransferHandler(transferSvc, accountSvc)

	authMw := middleware.NewAuthMiddleware(authSvc, zap.NewNop())
	rlCfg := &config.RateLimitConfig{Enabled: false}
	rlMw := middleware.NewRateLimitMiddleware(nil, rlCfg, zap.NewNop())

	return testutil.SetupGinRouter(t, testutil.RouterDeps{
		AuthHandler:          authHandler,
		AccountHandler:       accountHandler,
		MovementHandler:      movementHandler,
		TransferHandler:      transferHandler,
		AuthMiddleware:       authMw,
		RateLimitMiddleware:  rlMw,
	})
}

