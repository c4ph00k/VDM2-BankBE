package util_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"VDM2-BankBE/internal/util"
)

func TestHandleError(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		err            error
		expectedStatus int
		expectedCode   int
		expectedMsg    string
	}{
		{
			name:           "APIError preserved",
			err:            util.NewNotFoundError("account not found"),
			expectedStatus: http.StatusNotFound,
			expectedCode:   http.StatusNotFound,
			expectedMsg:    "account not found",
		},
		{
			name:           "generic error becomes 500 internal server error",
			err:            errors.New("boom"),
			expectedStatus: http.StatusInternalServerError,
			expectedCode:   http.StatusInternalServerError,
			expectedMsg:    "internal server error",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

			util.HandleError(c, tc.err)

			if rec.Code != tc.expectedStatus {
				t.Fatalf("unexpected status: got=%d want=%d body=%q", rec.Code, tc.expectedStatus, rec.Body.String())
			}

			var got util.ErrorResponse
			if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
				t.Fatalf("failed to unmarshal response: %v body=%q", err, rec.Body.String())
			}
			if got.Error == nil {
				t.Fatalf("expected error object, got: %q", rec.Body.String())
			}
			if got.Error.Code != tc.expectedCode {
				t.Fatalf("unexpected code: got=%d want=%d", got.Error.Code, tc.expectedCode)
			}
			if got.Error.Message != tc.expectedMsg {
				t.Fatalf("unexpected message: got=%q want=%q", got.Error.Message, tc.expectedMsg)
			}
		})
	}
}

