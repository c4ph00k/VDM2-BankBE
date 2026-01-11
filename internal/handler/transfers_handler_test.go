package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"

	"VDM2-BankBE/internal/model"
	servicemocks "VDM2-BankBE/internal/service/mocks"
	"VDM2-BankBE/internal/testutil"
	"VDM2-BankBE/internal/util"
)

func TestTransfers_Create(t *testing.T) {
	t.Parallel()

	token := "header.payload.sig"
	userID := uuid.MustParse("00000000-0000-0000-0000-000000000040")
	fromAccountID := uuid.MustParse("00000000-0000-0000-0000-000000000041")
	toAccountID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")

	user := &model.User{ID: userID}
	fromAccount := &model.Account{ID: fromAccountID, UserID: userID, Currency: "EUR"}

	tests := []struct {
		name           string
		setupAuth      func(headers map[string]string)
		requestBody    any
		buildMocks     func(ctrl *gomock.Controller) (*servicemocks.MockAuthService, *servicemocks.MockAccountService, *servicemocks.MockTransferService)
		expectedStatus int
		assertResponse func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name: "success",
			setupAuth: func(headers map[string]string) {
				headers["Authorization"] = "Bearer " + token
			},
			requestBody: map[string]any{"to_account": toAccountID.String(), "amount": "25.00", "description": "test"},
			buildMocks: func(ctrl *gomock.Controller) (*servicemocks.MockAuthService, *servicemocks.MockAccountService, *servicemocks.MockTransferService) {
				authSvc := servicemocks.NewMockAuthService(ctrl)
				accountSvc := servicemocks.NewMockAccountService(ctrl)
				transferSvc := servicemocks.NewMockTransferService(ctrl)

				amount := mustDecimal(t, "25.00")
				transfer := &model.Transfer{ID: 1, FromAccount: fromAccountID, ToAccount: toAccountID, Amount: amount, Status: "completed"}

				authSvc.EXPECT().VerifyToken(gomock.Any(), token).Return(user, nil)
				accountSvc.EXPECT().GetByUserID(gomock.Any(), userID).Return(fromAccount, nil)
				transferSvc.EXPECT().Transfer(gomock.Any(), fromAccountID, toAccountID, amount, "test").Return(transfer, nil)

				return authSvc, accountSvc, transferSvc
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
			buildMocks: func(ctrl *gomock.Controller) (*servicemocks.MockAuthService, *servicemocks.MockAccountService, *servicemocks.MockTransferService) {
				authSvc := servicemocks.NewMockAuthService(ctrl)
				accountSvc := servicemocks.NewMockAccountService(ctrl)
				transferSvc := servicemocks.NewMockTransferService(ctrl)
				authSvc.EXPECT().VerifyToken(gomock.Any(), token).Return(user, nil)
				return authSvc, accountSvc, transferSvc
			},
			expectedStatus: http.StatusBadRequest,
			assertResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				testutil.AssertHTTPError(t, rec, http.StatusBadRequest, "invalid request body")
			},
		},
		{
			name: "validation error",
			setupAuth: func(headers map[string]string) {
				headers["Authorization"] = "Bearer " + token
			},
			requestBody: map[string]any{"to_account": "not-a-uuid", "amount": "10.00"},
			buildMocks: func(ctrl *gomock.Controller) (*servicemocks.MockAuthService, *servicemocks.MockAccountService, *servicemocks.MockTransferService) {
				authSvc := servicemocks.NewMockAuthService(ctrl)
				accountSvc := servicemocks.NewMockAccountService(ctrl)
				transferSvc := servicemocks.NewMockTransferService(ctrl)
				authSvc.EXPECT().VerifyToken(gomock.Any(), token).Return(user, nil)
				return authSvc, accountSvc, transferSvc
			},
			expectedStatus: http.StatusBadRequest,
			assertResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				testutil.AssertHTTPError(t, rec, http.StatusBadRequest, http.StatusBadRequest)
			},
		},
		{
			name: "invalid amount",
			setupAuth: func(headers map[string]string) {
				headers["Authorization"] = "Bearer " + token
			},
			requestBody: map[string]any{"to_account": toAccountID.String(), "amount": "abc"},
			buildMocks: func(ctrl *gomock.Controller) (*servicemocks.MockAuthService, *servicemocks.MockAccountService, *servicemocks.MockTransferService) {
				authSvc := servicemocks.NewMockAuthService(ctrl)
				accountSvc := servicemocks.NewMockAccountService(ctrl)
				transferSvc := servicemocks.NewMockTransferService(ctrl)
				authSvc.EXPECT().VerifyToken(gomock.Any(), token).Return(user, nil)
				return authSvc, accountSvc, transferSvc
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
				"to_account": toAccountID.String(),
				"amount":     "10.00",
			},
			buildMocks: func(ctrl *gomock.Controller) (*servicemocks.MockAuthService, *servicemocks.MockAccountService, *servicemocks.MockTransferService) {
				return servicemocks.NewMockAuthService(ctrl), servicemocks.NewMockAccountService(ctrl), servicemocks.NewMockTransferService(ctrl)
			},
			expectedStatus: http.StatusUnauthorized,
			assertResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				testutil.AssertHTTPError(t, rec, http.StatusUnauthorized, "missing authorization token")
			},
		},
		{
			name: "service bad request maps to 400",
			setupAuth: func(headers map[string]string) {
				headers["Authorization"] = "Bearer " + token
			},
			requestBody: map[string]any{"to_account": toAccountID.String(), "amount": "10.00"},
			buildMocks: func(ctrl *gomock.Controller) (*servicemocks.MockAuthService, *servicemocks.MockAccountService, *servicemocks.MockTransferService) {
				authSvc := servicemocks.NewMockAuthService(ctrl)
				accountSvc := servicemocks.NewMockAccountService(ctrl)
				transferSvc := servicemocks.NewMockTransferService(ctrl)

				amount := mustDecimal(t, "10.00")

				authSvc.EXPECT().VerifyToken(gomock.Any(), token).Return(user, nil)
				accountSvc.EXPECT().GetByUserID(gomock.Any(), userID).Return(fromAccount, nil)
				transferSvc.EXPECT().Transfer(gomock.Any(), fromAccountID, toAccountID, amount, "").Return(nil, util.NewBadRequestError("insufficient funds"))
				return authSvc, accountSvc, transferSvc
			},
			expectedStatus: http.StatusBadRequest,
			assertResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				testutil.AssertHTTPError(t, rec, http.StatusBadRequest, "insufficient funds")
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			authSvc, accountSvc, transferSvc := tc.buildMocks(ctrl)
			r := newTestRouter(t, ctrl, authSvc, accountSvc, nil, transferSvc)

			headers := map[string]string{}
			tc.setupAuth(headers)
			req := testutil.NewJSONRequest(http.MethodPost, "/api/v1/transfers", tc.requestBody, headers)
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			tc.assertResponse(t, rec)
		})
	}
}

func TestTransfers_List(t *testing.T) {
	t.Parallel()

	token := "header.payload.sig"
	userID := uuid.MustParse("00000000-0000-0000-0000-000000000050")
	accountID := uuid.MustParse("00000000-0000-0000-0000-000000000051")

	user := &model.User{ID: userID}
	account := &model.Account{ID: accountID, UserID: userID, Currency: "EUR"}

	tests := []struct {
		name           string
		path           string
		setupAuth      func(headers map[string]string)
		buildMocks     func(ctrl *gomock.Controller) (*servicemocks.MockAuthService, *servicemocks.MockAccountService, *servicemocks.MockTransferService)
		expectedStatus int
		assertResponse func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name: "binding error on non-int query param uses router ErrorHandler envelope",
			path: "/api/v1/transfers?page=abc&limit=10",
			setupAuth: func(headers map[string]string) {
				// binding fails before auth middleware runs
			},
			buildMocks: func(ctrl *gomock.Controller) (*servicemocks.MockAuthService, *servicemocks.MockAccountService, *servicemocks.MockTransferService) {
				return servicemocks.NewMockAuthService(ctrl), servicemocks.NewMockAccountService(ctrl), servicemocks.NewMockTransferService(ctrl)
			},
			expectedStatus: http.StatusBadRequest,
			assertResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				testutil.AssertHTTPError(t, rec, http.StatusBadRequest, http.StatusBadRequest)
			},
		},
		{
			name: "success with clamped pagination (page<=0 -> 1, limit>100 -> 100)",
			path: "/api/v1/transfers?page=0&limit=200",
			setupAuth: func(headers map[string]string) {
				headers["Authorization"] = "Bearer " + token
			},
			buildMocks: func(ctrl *gomock.Controller) (*servicemocks.MockAuthService, *servicemocks.MockAccountService, *servicemocks.MockTransferService) {
				authSvc := servicemocks.NewMockAuthService(ctrl)
				accountSvc := servicemocks.NewMockAccountService(ctrl)
				transferSvc := servicemocks.NewMockTransferService(ctrl)

				authSvc.EXPECT().VerifyToken(gomock.Any(), token).Return(user, nil)
				accountSvc.EXPECT().GetByUserID(gomock.Any(), userID).Return(account, nil)
				transferSvc.EXPECT().GetByAccountID(gomock.Any(), accountID, 1, 100).Return(&util.PaginatedResponse{}, nil)

				return authSvc, accountSvc, transferSvc
			},
			expectedStatus: http.StatusOK,
			assertResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				testutil.AssertHTTPStatus(t, rec, http.StatusOK)
			},
		},
		{
			name:      "missing token returns 401",
			path:      "/api/v1/transfers?page=1&limit=10",
			setupAuth: func(headers map[string]string) {},
			buildMocks: func(ctrl *gomock.Controller) (*servicemocks.MockAuthService, *servicemocks.MockAccountService, *servicemocks.MockTransferService) {
				return servicemocks.NewMockAuthService(ctrl), servicemocks.NewMockAccountService(ctrl), servicemocks.NewMockTransferService(ctrl)
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

			authSvc, accountSvc, transferSvc := tc.buildMocks(ctrl)
			r := newTestRouter(t, ctrl, authSvc, accountSvc, nil, transferSvc)

			headers := map[string]string{}
			tc.setupAuth(headers)
			req := testutil.NewJSONRequest(http.MethodGet, tc.path, nil, headers)
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			tc.assertResponse(t, rec)
		})
	}
}

