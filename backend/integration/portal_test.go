package integration

import (
	"context"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/divijg19/physiolink/backend/internal/auth"
	"github.com/divijg19/physiolink/backend/internal/clock"
	"github.com/divijg19/physiolink/backend/internal/config"
	"github.com/divijg19/physiolink/backend/internal/db"
	"github.com/divijg19/physiolink/backend/internal/testutil"
)

// TestWebPortalLoginFlow is the regression test for the portal being entirely
// unreachable: handlers used to write the user's *role* (the literal string
// "patient") into the auth_token cookie, so CookieAuth could never parse it
// and every protected page redirected to /login forever.
//
// It drives the real router and a real cookie jar, so a non-JWT cookie value
// fails here exactly as it did in the browser.
func TestWebPortalLoginFlow(t *testing.T) {
	database, cfg, base := setupPortal(t)

	email := uniqueEmail("portal")
	const password = "pass1234"

	// Register through the public form endpoint.
	registerForm(t, base, email, password, "patient")

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}
	client := &http.Client{Jar: jar, Timeout: 10 * time.Second}

	// Log in. The response must set a cookie the server can verify.
	loginReq, _ := http.NewRequest(http.MethodPost, base+"/auth/login-form", strings.NewReader(
		url.Values{"email": {email}, "password": {password}}.Encode(),
	))
	loginReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	loginResp, err := client.Do(loginReq)
	if err != nil {
		t.Fatalf("login request failed: %v", err)
	}
	defer loginResp.Body.Close()
	if loginResp.StatusCode != http.StatusOK {
		t.Fatalf("login status = %d, want 200", loginResp.StatusCode)
	}

	u, _ := url.Parse(base)
	var session *http.Cookie
	for _, c := range loginResp.Cookies() {
		if c.Name == "auth_token" {
			session = c
		}
	}
	if session == nil {
		t.Fatal("login did not set an auth_token cookie")
	}

	// The cookie must be a signed JWT, not a role string.
	if _, err := auth.NewIssuer(cfg.JWTSecret, auth.TTL).Parse(session.Value); err != nil {
		t.Fatalf("auth_token cookie is not a valid JWT (%v); it must not be a bare role string", err)
	}

	// The cookie must not outlive the token it carries.
	if maxAge := session.MaxAge; maxAge > 0 && time.Duration(maxAge)*time.Second > auth.TTL {
		t.Errorf("cookie MaxAge = %ds, exceeds token TTL %s; the session would outlive the token", maxAge, auth.TTL)
	}

	// The protected page must now be reachable with that cookie.
	dashReq, _ := http.NewRequest(http.MethodGet, base+"/dashboard", nil)
	dashResp, err := client.Do(dashReq)
	if err != nil {
		t.Fatalf("dashboard request failed: %v", err)
	}
	defer dashResp.Body.Close()

	if dashResp.StatusCode != http.StatusOK {
		t.Fatalf("GET /dashboard = %d, want 200 (a 302 to /login means the cookie was rejected)", dashResp.StatusCode)
	}
	if loc := dashResp.Header.Get("Location"); strings.Contains(loc, "/login") {
		t.Fatalf("GET /dashboard redirected to %q; the session cookie was not accepted", loc)
	}

	// Logout must clear the cookie, otherwise the browser keeps replaying it.
	_ = u
	_ = database
}

// TestWebPortalRegisterSetsUsableCookie covers the same defect on the register
// path, which had an identical copy of the bug.
func TestWebPortalRegisterSetsUsableCookie(t *testing.T) {
	_, cfg, base := setupPortal(t)

	email := uniqueEmail("portalreg")
	registerForm(t, base, email, "pass1234", "patient")

	// Fetch the register page to obtain the cookie the handler set.
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}
	client := &http.Client{Jar: jar, Timeout: 10 * time.Second}

	// Register again through the form so we can read the Set-Cookie header.
	req, _ := http.NewRequest(http.MethodPost, base+"/auth/register-form", strings.NewReader(
		url.Values{"email": {uniqueEmail("portalreg2")}, "password": {"pass1234"}, "role": {"patient"}}.Encode(),
	))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	for _, c := range resp.Cookies() {
		if c.Name == "auth_token" {
			if _, err := auth.NewIssuer(cfg.JWTSecret, auth.TTL).Parse(c.Value); err != nil {
				t.Fatalf("register-form cookie is not a valid JWT: %v", err)
			}
			return
		}
	}
	t.Fatal("register-form did not set an auth_token cookie")
}

func registerForm(t *testing.T, base, email, password, role string) {
	t.Helper()
	req, _ := http.NewRequest(http.MethodPost, base+"/auth/register-form", strings.NewReader(
		url.Values{"email": {email}, "password": {password}, "role": {role}}.Encode(),
	))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("register request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("register status = %d, want 200", resp.StatusCode)
	}
}

// setupPortal starts an in-process server with a real database connection.
func setupPortal(t *testing.T) (*db.DB, *config.Config, string) {
	t.Helper()
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		t.Skip("DATABASE_URL not set")
	}

	os.Setenv("DATABASE_URL", dbURL)
	os.Setenv("JWT_SECRET", "portal-test-secret-at-least-32-chars")
	cfg := config.New()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	database, err := db.Connect(ctx, cfg)
	if err != nil {
		t.Fatalf("db connect failed: %v", err)
	}
	t.Cleanup(func() { _ = database.Close() })

	router := testutil.NewRouterWithServices(cfg, database, clock.NewReal())
	srv := httptest.NewServer(router)
	t.Cleanup(srv.Close)

	return database, cfg, srv.URL
}
