package testutil

import (
	"net/http/httptest"
	"testing"

	"VDM2-BankBE/internal/util"
)

func AssertHTTPStatus(t *testing.T, rec *httptest.ResponseRecorder, expected int) {
	t.Helper()
	if rec.Code != expected {
		t.Fatalf("unexpected status: got=%d want=%d body=%q", rec.Code, expected, rec.Body.String())
	}
}

func AssertHTTPError(t *testing.T, rec *httptest.ResponseRecorder, expectedStatus int, expectedCodeOrMessage any) {
	t.Helper()
	AssertHTTPStatus(t, rec, expectedStatus)

	resp := DecodeJSONResponse[util.ErrorResponse](t, rec)
	if resp.Error == nil {
		t.Fatalf("expected error response body, got: %q", rec.Body.String())
	}

	switch v := expectedCodeOrMessage.(type) {
	case int:
		if resp.Error.Code != v {
			t.Fatalf("unexpected error.code: got=%d want=%d body=%q", resp.Error.Code, v, rec.Body.String())
		}
	case string:
		if resp.Error.Message != v {
			t.Fatalf("unexpected error.message: got=%q want=%q body=%q", resp.Error.Message, v, rec.Body.String())
		}
	default:
		t.Fatalf("unsupported expectedCodeOrMessage type %T", expectedCodeOrMessage)
	}
}

