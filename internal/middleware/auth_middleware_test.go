package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"VDM2-BankBE/internal/generated"
	"VDM2-BankBE/internal/middleware"
	"VDM2-BankBE/internal/model"
	servicemocks "VDM2-BankBE/internal/service/mocks"
	"VDM2-BankBE/internal/testutil"
	"VDM2-BankBE/internal/util"
)

func TestAuthMiddleware_AuthenticateIfRequiredFunc(t *testing.T) {
	t.Parallel()

	const (
		jwtToken    = "header.payload.sig"
		pasetoToken = "v2.local.test"
	)

	userID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	user := &model.User{ID: userID}

	type response struct {
		UserID string `json:"user_id"`
	}

	tests := []struct {
		name               string
		setSecurityMarkers bool
		setupAuth          func(h http.Header)
		buildMocks         func(ctrl *gomock.Controller) *servicemocks.MockAuthService
		expectedStatus     int
		assertResponse     func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name:               "public operation skips auth",
			setSecurityMarkers: false,
			setupAuth:          func(h http.Header) {},
			buildMocks: func(ctrl *gomock.Controller) *servicemocks.MockAuthService {
				return servicemocks.NewMockAuthService(ctrl)
			},
			expectedStatus: http.StatusOK,
			assertResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				testutil.AssertHTTPStatus(t, rec, http.StatusOK)
			},
		},
		{
			name:               "missing Authorization header returns 401",
			setSecurityMarkers: true,
			setupAuth:          func(h http.Header) {},
			buildMocks: func(ctrl *gomock.Controller) *servicemocks.MockAuthService {
				return servicemocks.NewMockAuthService(ctrl)
			},
			expectedStatus: http.StatusUnauthorized,
			assertResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				testutil.AssertHTTPError(t, rec, http.StatusUnauthorized, "missing authorization token")
			},
		},
		{
			name:               "wrong scheme returns 401",
			setSecurityMarkers: true,
			setupAuth: func(h http.Header) {
				h.Set("Authorization", "Basic abc")
			},
			buildMocks: func(ctrl *gomock.Controller) *servicemocks.MockAuthService {
				return servicemocks.NewMockAuthService(ctrl)
			},
			expectedStatus: http.StatusUnauthorized,
			assertResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				testutil.AssertHTTPError(t, rec, http.StatusUnauthorized, "missing authorization token")
			},
		},
		{
			name:               "invalid token returns 401",
			setSecurityMarkers: true,
			setupAuth: func(h http.Header) {
				h.Set("Authorization", "Bearer "+jwtToken)
			},
			buildMocks: func(ctrl *gomock.Controller) *servicemocks.MockAuthService {
				m := servicemocks.NewMockAuthService(ctrl)
				m.EXPECT().
					VerifyToken(gomock.AssignableToTypeOf(&gin.Context{}), jwtToken).
					Return(nil, util.NewUnauthorizedError("invalid token"))
				return m
			},
			expectedStatus: http.StatusUnauthorized,
			assertResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				testutil.AssertHTTPError(t, rec, http.StatusUnauthorized, "invalid or expired token")
			},
		},
		{
			name:               "expired token returns 401",
			setSecurityMarkers: true,
			setupAuth: func(h http.Header) {
				h.Set("Authorization", "Bearer "+jwtToken)
			},
			buildMocks: func(ctrl *gomock.Controller) *servicemocks.MockAuthService {
				m := servicemocks.NewMockAuthService(ctrl)
				m.EXPECT().
					VerifyToken(gomock.AssignableToTypeOf(&gin.Context{}), jwtToken).
					Return(nil, util.NewUnauthorizedError("token expired"))
				return m
			},
			expectedStatus: http.StatusUnauthorized,
			assertResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				testutil.AssertHTTPError(t, rec, http.StatusUnauthorized, "invalid or expired token")
			},
		},
		{
			name:               "valid JWT returns 200 and sets user",
			setSecurityMarkers: true,
			setupAuth: func(h http.Header) {
				h.Set("Authorization", "Bearer "+jwtToken)
			},
			buildMocks: func(ctrl *gomock.Controller) *servicemocks.MockAuthService {
				m := servicemocks.NewMockAuthService(ctrl)
				m.EXPECT().
					VerifyToken(gomock.AssignableToTypeOf(&gin.Context{}), jwtToken).
					Return(user, nil)
				return m
			},
			expectedStatus: http.StatusOK,
			assertResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				testutil.AssertHTTPStatus(t, rec, http.StatusOK)
				got := testutil.DecodeJSONResponse[response](t, rec)
				if got.UserID != userID.String() {
					t.Fatalf("unexpected user_id: got=%q want=%q", got.UserID, userID.String())
				}
			},
		},
		{
			name:               "valid PASETO returns 200 and sets user",
			setSecurityMarkers: true,
			setupAuth: func(h http.Header) {
				h.Set("Authorization", "Bearer "+pasetoToken)
			},
			buildMocks: func(ctrl *gomock.Controller) *servicemocks.MockAuthService {
				m := servicemocks.NewMockAuthService(ctrl)
				m.EXPECT().
					VerifyToken(gomock.AssignableToTypeOf(&gin.Context{}), pasetoToken).
					Return(user, nil)
				return m
			},
			expectedStatus: http.StatusOK,
			assertResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				testutil.AssertHTTPStatus(t, rec, http.StatusOK)
				got := testutil.DecodeJSONResponse[response](t, rec)
				if got.UserID != userID.String() {
					t.Fatalf("unexpected user_id: got=%q want=%q", got.UserID, userID.String())
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

			authSvc := tc.buildMocks(ctrl)
			mw := middleware.NewAuthMiddleware(authSvc, zap.NewNop())
			authFn := mw.AuthenticateIfRequiredFunc()

			r := gin.New()
			r.GET("/test", func(c *gin.Context) {
				if tc.setSecurityMarkers {
					c.Set(generated.BearerJWTScopes, []string{})
					c.Set(generated.BearerPASETOScopes, []string{})
				}

				authFn(c)
				if c.IsAborted() {
					return
				}

				u, _ := c.Get("user")
				userModel, _ := u.(*model.User)
				if userModel == nil {
					c.JSON(http.StatusOK, gin.H{})
					return
				}
				c.JSON(http.StatusOK, response{UserID: userModel.ID.String()})
			})

			rec := httptest.NewRecorder()
			req := testutil.NewJSONRequest(http.MethodGet, "/test", nil, nil)
			tc.setupAuth(req.Header)
			r.ServeHTTP(rec, req)

			if tc.assertResponse != nil {
				tc.assertResponse(t, rec)
			} else {
				testutil.AssertHTTPStatus(t, rec, tc.expectedStatus)
			}
		})
	}
}

