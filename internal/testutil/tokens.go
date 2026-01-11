package testutil

import (
	"crypto/sha256"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/o1egl/paseto"
)

const (
	DefaultValidExpTime   = "2030-01-01T00:00:00Z"
	DefaultExpiredExpTime = "2000-01-01T00:00:00Z"
)

func ValidJWT(t *testing.T, secret string, userID string) string {
	t.Helper()
	return jwtWithExpRFC3339(t, secret, userID, DefaultValidExpTime)
}

func ExpiredJWT(t *testing.T, secret string, userID string) string {
	t.Helper()
	return jwtWithExpRFC3339(t, secret, userID, DefaultExpiredExpTime)
}

func jwtWithExpRFC3339(t *testing.T, secret string, userID string, expRFC3339 string) string {
	t.Helper()
	exp, err := time.Parse(time.RFC3339, expRFC3339)
	if err != nil {
		t.Fatalf("bad exp time: %v", err)
	}

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     exp.Unix(),
		"iat":     time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
		"iss":     "VDM2-Bank",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("failed to sign JWT: %v", err)
	}
	return s
}

func ValidPASETO(t *testing.T, secret string, userID string) string {
	t.Helper()
	return pasetoV2LocalWithExpRFC3339(t, secret, userID, DefaultValidExpTime)
}

func ExpiredPASETO(t *testing.T, secret string, userID string) string {
	t.Helper()
	return pasetoV2LocalWithExpRFC3339(t, secret, userID, DefaultExpiredExpTime)
}

func pasetoV2LocalWithExpRFC3339(t *testing.T, secret string, userID string, expRFC3339 string) string {
	t.Helper()
	key := sha256.Sum256([]byte(secret))
	claims := map[string]any{
		"user_id": userID,
		"exp":     expRFC3339,
	}
	v2 := paseto.NewV2()
	token, err := v2.Encrypt(key[:], claims, nil)
	if err != nil {
		t.Fatalf("failed to encrypt PASETO: %v", err)
	}
	return token
}

