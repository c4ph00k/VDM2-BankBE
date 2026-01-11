package handler_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"VDM2-BankBE/internal/model"
	servicemocks "VDM2-BankBE/internal/service/mocks"
	"VDM2-BankBE/internal/testutil"
	"VDM2-BankBE/internal/util"
)

func TestAccounts_Balance(t *testing.T) {
	t.Parallel()

	token := "header.payload.sig"
	userID := uuid.MustParse("00000000-0000-0000-0000-000000000010")
	accountID := uuid.MustParse("00000000-0000-0000-0000-000000000011")
	user := &model.User{ID: userID}
	account := &model.Account{ID: accountID, UserID: userID, Currency: "EUR"}

	tests := []struct {
		name           string
		setupAuth      func(headers map[string]string)
		buildMocks     func(ctrl *gomock.Controller) (*servicemocks.MockAuthService, *servicemocks.MockAccountService)
		expectedStatus int
		assertResponse func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name: "success",
			setupAuth: func(headers map[string]string) {
				headers["Authorization"] = "Bearer " + token
			},
			buildMocks: func(ctrl *gomock.Controller) (*servicemocks.MockAuthService, *servicemocks.MockAccountService) {
				authSvc := servicemocks.NewMockAuthService(ctrl)
				accountSvc := servicemocks.NewMockAccountService(ctrl)

				authSvc.EXPECT().VerifyToken(gomock.Any(), token).Return(user, nil)
				accountSvc.EXPECT().GetByUserID(gomock.Any(), userID).Return(account, nil)
				accountSvc.EXPECT().GetBalance(gomock.Any(), accountID).Return(mustDecimal(t, "10.50"), nil)

				return authSvc, accountSvc
			},
			expectedStatus: http.StatusOK,
			assertResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				testutil.AssertHTTPStatus(t, rec, http.StatusOK)
			},
		},
		{
			name:      "missing token returns 401",
			setupAuth: func(headers map[string]string) {},
			buildMocks: func(ctrl *gomock.Controller) (*servicemocks.MockAuthService, *servicemocks.MockAccountService) {
				return servicemocks.NewMockAuthService(ctrl), servicemocks.NewMockAccountService(ctrl)
			},
			expectedStatus: http.StatusUnauthorized,
			assertResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				testutil.AssertHTTPError(t, rec, http.StatusUnauthorized, "missing authorization token")
			},
		},
		{
			name: "wrong scheme returns 401",
			setupAuth: func(headers map[string]string) {
				headers["Authorization"] = "Basic abc"
			},
			buildMocks: func(ctrl *gomock.Controller) (*servicemocks.MockAuthService, *servicemocks.MockAccountService) {
				return servicemocks.NewMockAuthService(ctrl), servicemocks.NewMockAccountService(ctrl)
			},
			expectedStatus: http.StatusUnauthorized,
			assertResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				testutil.AssertHTTPError(t, rec, http.StatusUnauthorized, "missing authorization token")
			},
		},
		{
			name: "invalid token returns 401",
			setupAuth: func(headers map[string]string) {
				headers["Authorization"] = "Bearer " + token
			},
			buildMocks: func(ctrl *gomock.Controller) (*servicemocks.MockAuthService, *servicemocks.MockAccountService) {
				authSvc := servicemocks.NewMockAuthService(ctrl)
				accountSvc := servicemocks.NewMockAccountService(ctrl)
				authSvc.EXPECT().VerifyToken(gomock.Any(), token).Return(nil, util.NewUnauthorizedError("invalid token"))
				return authSvc, accountSvc
			},
			expectedStatus: http.StatusUnauthorized,
			assertResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				testutil.AssertHTTPError(t, rec, http.StatusUnauthorized, "invalid or expired token")
			},
		},
		{
			name: "not found maps to 404",
			setupAuth: func(headers map[string]string) {
				headers["Authorization"] = "Bearer " + token
			},
			buildMocks: func(ctrl *gomock.Controller) (*servicemocks.MockAuthService, *servicemocks.MockAccountService) {
				authSvc := servicemocks.NewMockAuthService(ctrl)
				accountSvc := servicemocks.NewMockAccountService(ctrl)
				authSvc.EXPECT().VerifyToken(gomock.Any(), token).Return(user, nil)
				accountSvc.EXPECT().GetByUserID(gomock.Any(), userID).Return(nil, util.NewNotFoundError("account not found"))
				return authSvc, accountSvc
			},
			expectedStatus: http.StatusNotFound,
			assertResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				testutil.AssertHTTPError(t, rec, http.StatusNotFound, "account not found")
			},
		},
		{
			name: "generic error maps to 500",
			setupAuth: func(headers map[string]string) {
				headers["Authorization"] = "Bearer " + token
			},
			buildMocks: func(ctrl *gomock.Controller) (*servicemocks.MockAuthService, *servicemocks.MockAccountService) {
				authSvc := servicemocks.NewMockAuthService(ctrl)
				accountSvc := servicemocks.NewMockAccountService(ctrl)
				authSvc.EXPECT().VerifyToken(gomock.Any(), token).Return(user, nil)
				accountSvc.EXPECT().GetByUserID(gomock.Any(), userID).Return(nil, errors.New("db down"))
				return authSvc, accountSvc
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

			authSvc, accountSvc := tc.buildMocks(ctrl)
			r := newTestRouter(t, ctrl, authSvc, accountSvc, nil, nil)

			headers := map[string]string{}
			tc.setupAuth(headers)
			req := testutil.NewJSONRequest(http.MethodGet, "/api/v1/accounts/balance", nil, headers)
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			tc.assertResponse(t, rec)
		})
	}
}

func TestAccounts_ListMovements(t *testing.T) {
	t.Parallel()

	token := "header.payload.sig"
	userID := uuid.MustParse("00000000-0000-0000-0000-000000000020")
	accountID := uuid.MustParse("00000000-0000-0000-0000-000000000021")
	user := &model.User{ID: userID}
	account := &model.Account{ID: accountID, UserID: userID, Currency: "EUR"}

	tests := []struct {
		name           string
		path           string
		setupAuth      func(headers map[string]string)
		buildMocks     func(ctrl *gomock.Controller) (*servicemocks.MockAuthService, *servicemocks.MockAccountService, *servicemocks.MockMovementService)
		expectedStatus int
		assertResponse func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name: "binding error on non-int query param uses router ErrorHandler envelope",
			path: "/api/v1/accounts/movements?page=abc&limit=10",
			setupAuth: func(headers map[string]string) {
				// binding fails before auth middleware runs
			},
			buildMocks: func(ctrl *gomock.Controller) (*servicemocks.MockAuthService, *servicemocks.MockAccountService, *servicemocks.MockMovementService) {
				return servicemocks.NewMockAuthService(ctrl), servicemocks.NewMockAccountService(ctrl), servicemocks.NewMockMovementService(ctrl)
			},
			expectedStatus: http.StatusBadRequest,
			assertResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				testutil.AssertHTTPError(t, rec, http.StatusBadRequest, http.StatusBadRequest)
			},
		},
		{
			name: "success with clamped pagination (page<=0 -> 1, limit>100 -> 100)",
			path: "/api/v1/accounts/movements?page=0&limit=200",
			setupAuth: func(headers map[string]string) {
				headers["Authorization"] = "Bearer " + token
			},
			buildMocks: func(ctrl *gomock.Controller) (*servicemocks.MockAuthService, *servicemocks.MockAccountService, *servicemocks.MockMovementService) {
				authSvc := servicemocks.NewMockAuthService(ctrl)
				accountSvc := servicemocks.NewMockAccountService(ctrl)
				movementSvc := servicemocks.NewMockMovementService(ctrl)

				authSvc.EXPECT().VerifyToken(gomock.Any(), token).Return(user, nil)
				accountSvc.EXPECT().GetByUserID(gomock.Any(), userID).Return(account, nil)
				movementSvc.EXPECT().GetByAccountID(gomock.Any(), accountID, 1, 100).Return(&util.PaginatedResponse{}, nil)

				return authSvc, accountSvc, movementSvc
			},
			expectedStatus: http.StatusOK,
			assertResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				testutil.AssertHTTPStatus(t, rec, http.StatusOK)
			},
		},
		{
			name:      "missing token returns 401",
			path:      "/api/v1/accounts/movements?page=1&limit=10",
			setupAuth: func(headers map[string]string) {},
			buildMocks: func(ctrl *gomock.Controller) (*servicemocks.MockAuthService, *servicemocks.MockAccountService, *servicemocks.MockMovementService) {
				return servicemocks.NewMockAuthService(ctrl), servicemocks.NewMockAccountService(ctrl), servicemocks.NewMockMovementService(ctrl)
			},
			expectedStatus: http.StatusUnauthorized,
			assertResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				testutil.AssertHTTPError(t, rec, http.StatusUnauthorized, "missing authorization token")
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			authSvc, accountSvc, movementSvc := tc.buildMocks(ctrl)
			r := newTestRouter(t, ctrl, authSvc, accountSvc, movementSvc, nil)

			headers := map[string]string{}
			tc.setupAuth(headers)
			req := testutil.NewJSONRequest(http.MethodGet, tc.path, nil, headers)
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			tc.assertResponse(t, rec)
		})
	}
}

func TestAccounts_CreateMovement(t *testing.T) {
	t.Parallel()

	token := "header.payload.sig"
	userID := uuid.MustParse("00000000-0000-0000-0000-000000000030")
	accountID := uuid.MustParse("00000000-0000-0000-0000-000000000031")
	user := &model.User{ID: userID}
	account := &model.Account{ID: accountID, UserID: userID, Currency: "EUR", Balance: mustDecimal(t, "100.00")}

	tests := []struct {
		name           string
		setupAuth      func(headers map[string]string)
		requestBody    any
		buildMocks     func(ctrl *gomock.Controller) (*servicemocks.MockAuthService, *servicemocks.MockAccountService, *servicemocks.MockMovementService)
		expectedStatus int
		assertResponse func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name: "success",
			setupAuth: func(headers map[string]string) {
				headers["Authorization"] = "Bearer " + token
			},
			requestBody: map[string]any{"amount": "10.50", "type": "credit", "description": "deposit"},
			buildMocks: func(ctrl *gomock.Controller) (*servicemocks.MockAuthService, *servicemocks.MockAccountService, *servicemocks.MockMovementService) {
				authSvc := servicemocks.NewMockAuthService(ctrl)
				accountSvc := servicemocks.NewMockAccountService(ctrl)
				movementSvc := servicemocks.NewMockMovementService(ctrl)

				amount := mustDecimal(t, "10.50")
				movement := &model.Movement{ID: 1, AccountID: accountID, Amount: amount, Type: "credit", Description: "deposit"}

				authSvc.EXPECT().VerifyToken(gomock.Any(), token).Return(user, nil)
				accountSvc.EXPECT().GetByUserID(gomock.Any(), userID).Return(account, nil)
				movementSvc.EXPECT().Create(gomock.Any(), accountID, amount, "credit", "deposit").Return(movement, nil)

				return authSvc, accountSvc, movementSvc
			},
			expectedStatus: http.StatusCreated,
			assertResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				testutil.AssertHTTPStatus(t, rec, http.StatusCreated)
			},
		},
		{
			name: "invalid request body",
			setupAuth: func(headers map[string]string) {
				headers["Authorization"] = "Bearer " + token
			},
			requestBody: nil,
			buildMocks: func(ctrl *gomock.Controller) (*servicemocks.MockAuthService, *servicemocks.MockAccountService, *servicemocks.MockMovementService) {
				authSvc := servicemocks.NewMockAuthService(ctrl)
				accountSvc := servicemocks.NewMockAccountService(ctrl)
				movementSvc := servicemocks.NewMockMovementService(ctrl)
				authSvc.EXPECT().VerifyToken(gomock.Any(), token).Return(user, nil)
				return authSvc, accountSvc, movementSvc
			},
			expectedStatus: http.StatusBadRequest,
			assertResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				testutil.AssertHTTPError(t, rec, http.StatusBadRequest, "invalid request body")
			},
		},
		{
			name: "invalid amount",
			setupAuth: func(headers map[string]string) {
				headers["Authorization"] = "Bearer " + token
			},
			requestBody: map[string]any{"amount": "abc", "type": "credit"},
			buildMocks: func(ctrl *gomock.Controller) (*servicemocks.MockAuthService, *servicemocks.MockAccountService, *servicemocks.MockMovementService) {
				authSvc := servicemocks.NewMockAuthService(ctrl)
				accountSvc := servicemocks.NewMockAccountService(ctrl)
				movementSvc := servicemocks.NewMockMovementService(ctrl)
				authSvc.EXPECT().VerifyToken(gomock.Any(), token).Return(user, nil)
				return authSvc, accountSvc, movementSvc
			},
			expectedStatus: http.StatusBadRequest,
			assertResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				testutil.AssertHTTPError(t, rec, http.StatusBadRequest, "invalid amount")
			},
		},
		{
			name:      "missing token returns 401",
			setupAuth: func(headers map[string]string) {},
			requestBody: map[string]any{
				"amount": "10.00",
				"type":   "credit",
			},
			buildMocks: func(ctrl *gomock.Controller) (*servicemocks.MockAuthService, *servicemocks.MockAccountService, *servicemocks.MockMovementService) {
				return servicemocks.NewMockAuthService(ctrl), servicemocks.NewMockAccountService(ctrl), servicemocks.NewMockMovementService(ctrl)
			},
			expectedStatus: http.StatusUnauthorized,
			assertResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				testutil.AssertHTTPError(t, rec, http.StatusUnauthorized, "missing authorization token")
			},
		},
		{
			name: "service generic error maps to 500",
			setupAuth: func(headers map[string]string) {
				headers["Authorization"] = "Bearer " + token
			},
			requestBody: map[string]any{"amount": "10.00", "type": "credit"},
			buildMocks: func(ctrl *gomock.Controller) (*servicemocks.MockAuthService, *servicemocks.MockAccountService, *servicemocks.MockMovementService) {
				authSvc := servicemocks.NewMockAuthService(ctrl)
				accountSvc := servicemocks.NewMockAccountService(ctrl)
				movementSvc := servicemocks.NewMockMovementService(ctrl)

				amount := mustDecimal(t, "10.00")

				authSvc.EXPECT().VerifyToken(gomock.Any(), token).Return(user, nil)
				accountSvc.EXPECT().GetByUserID(gomock.Any(), userID).Return(account, nil)
				movementSvc.EXPECT().Create(gomock.Any(), accountID, amount, "credit", "").Return(nil, errors.New("boom"))

				return authSvc, accountSvc, movementSvc
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

			authSvc, accountSvc, movementSvc := tc.buildMocks(ctrl)
			r := newTestRouter(t, ctrl, authSvc, accountSvc, movementSvc, nil)

			headers := map[string]string{}
			tc.setupAuth(headers)
			req := testutil.NewJSONRequest(http.MethodPost, "/api/v1/accounts/movements", tc.requestBody, headers)
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			tc.assertResponse(t, rec)
		})
	}
}

func mustDecimal(t *testing.T, s string) decimal.Decimal {
	t.Helper()
	d, err := decimal.NewFromString(s)
	if err != nil {
		t.Fatalf("bad decimal %q: %v", s, err)
	}
	return d
}

