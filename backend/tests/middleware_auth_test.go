package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/divijg19/physiolink/backend/internal/config"
	mware "github.com/divijg19/physiolink/backend/internal/middleware"
)

// nextHandler echoes 200 if it sees user id in context
func nextHandler(w http.ResponseWriter, r *http.Request) {
	if _, ok := r.Context().Value(mware.UserIDKey).(string); !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func makeToken(t *testing.T, secret, uid, role string) string {
	t.Helper()
	claims := jwt.MapClaims{
		"user": map[string]string{"id": uid, "role": role},
		"exp":  time.Now().Add(5 * time.Minute).Unix(),
		"iat":  time.Now().Add(-1 * time.Minute).Unix(),
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := tok.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}
	return s
}

func TestJWTAuth_AllowsAuthorizationHeader(t *testing.T) {
	cfg := config.New()
	token := makeToken(t, cfg.JWTSecret, "11111111-1111-1111-1111-111111111111", "patient")

	mw := mware.JWTAuth(cfg)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	mw(http.HandlerFunc(nextHandler)).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

func TestJWTAuth_AllowsXAuthTokenHeader(t *testing.T) {
	cfg := config.New()
	token := makeToken(t, cfg.JWTSecret, "11111111-1111-1111-1111-111111111111", "patient")

	mw := mware.JWTAuth(cfg)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("x-auth-token", token)

	mw(http.HandlerFunc(nextHandler)).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

func TestJWTAuth_RejectsMissingToken(t *testing.T) {
	cfg := config.New()
	mw := mware.JWTAuth(cfg)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)

	mw(http.HandlerFunc(nextHandler)).ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rr.Code)
	}
}
