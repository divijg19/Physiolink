package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/divijg19/physiolink/backend/internal/config"
)

// TestSPAPrefixRouting locks the routing contract introduced when the Jaspr SPA
// moved off the server root:
//
//   - /site/* is served by the SPA
//   - /login and /register stay with the server-rendered portal
//   - unknown /api/* paths return a JSON 404 instead of the SPA shell
//
// The SPA is only embedded when built assets are present, so the SPA-specific
// assertions are skipped on a clean checkout (which ships placeholder.txt).
func TestSPAPrefixRouting(t *testing.T) {
	cfg := &config.Config{JWTSecret: "test-secret-not-used-here"}
	router := NewRouter(cfg)

	tests := []struct {
		name       string
		path       string
		wantStatus int
		wantBody   string
	}{
		{
			name:       "health is reachable at the root",
			path:       "/health",
			wantStatus: http.StatusOK,
		},
		{
			name:       "portal login page stays at the root",
			path:       "/login",
			wantStatus: http.StatusOK,
		},
		{
			name:       "portal register page stays at the root",
			path:       "/register",
			wantStatus: http.StatusOK,
		},
		{
			name:       "unknown api path is a JSON 404, not the SPA",
			path:       "/api/does-not-exist",
			wantStatus: http.StatusNotFound,
			wantBody:   `"msg"`,
		},
		{
			name:       "bare /api is also a JSON 404",
			path:       "/api",
			wantStatus: http.StatusNotFound,
			wantBody:   `"msg"`,
		},
		{
			name:       "unknown root path is a plain 404",
			path:       "/definitely-not-a-route",
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, tc.path, nil)
			router.ServeHTTP(rec, req)

			if rec.Code != tc.wantStatus {
				t.Fatalf("GET %s = %d, want %d (body: %.120q)", tc.path, rec.Code, tc.wantStatus, rec.Body.String())
			}
			if tc.wantBody != "" && !strings.Contains(rec.Body.String(), tc.wantBody) {
				t.Errorf("GET %s body = %.200q, want it to contain %q", tc.path, rec.Body.String(), tc.wantBody)
			}
		})
	}
}

// TestAPIErrorsAreJSONNotHTML guards the 404 body shape: a client that gets
// index.html (or a text/plain body) for an API path cannot parse the error.
func TestAPIErrorsAreJSONNotHTML(t *testing.T) {
	router := NewRouter(&config.Config{JWTSecret: "test-secret"})

	for _, path := range []string{"/api", "/api/", "/api/does-not-exist"} {
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))

		body := rec.Body.String()
		if strings.Contains(strings.ToLower(body), "<html") {
			t.Errorf("GET %s returned HTML to an API client: %.120q", path, body)
		}
		if ct := rec.Header().Get("Content-Type"); !strings.Contains(ct, "application/json") {
			t.Errorf("GET %s Content-Type = %q, want application/json", path, ct)
		}
	}
}

// TestSPARedirectsToSlash makes sure /site (no trailing slash) lands on the
// SPA instead of 404ing, which is what a shared marketing link would hit.
func TestSPARedirectsToSlash(t *testing.T) {
	cfg := &config.Config{JWTSecret: "test-secret-not-used-here"}
	router := NewRouter(cfg)

	for _, from := range []string{spaPrefix, "/"} {
		t.Run("from "+from, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, from, nil)
			router.ServeHTTP(rec, req)

			if rec.Code != http.StatusMovedPermanently {
				t.Fatalf("GET %s = %d, want %d", from, rec.Code, http.StatusMovedPermanently)
			}
			if got := rec.Header().Get("Location"); got != spaPrefix+"/" {
				t.Errorf("GET %s Location = %q, want %q", from, got, spaPrefix+"/")
			}
		})
	}
}

// TestRootRedirectsToSPA guards the portal's own navigation: views/layout.templ
// links to "/" twice (the wordmark and "Home"). Before /site was introduced the
// SPA answered "/" as its NotFound fallback; after the split, "/" 404ed and
// those links broke.
func TestRootRedirectsToSPA(t *testing.T) {
	router := NewRouter(&config.Config{JWTSecret: "test-secret"})

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	if rec.Code == http.StatusNotFound {
		t.Fatal(`GET / returned 404; the SPA must be reachable and layout.templ links to "/"`)
	}
	if rec.Code != http.StatusMovedPermanently {
		t.Fatalf("GET / = %d, want %d", rec.Code, http.StatusMovedPermanently)
	}
	if got := rec.Header().Get("Location"); got != spaPrefix+"/" {
		t.Errorf("GET / Location = %q, want %q", got, spaPrefix+"/")
	}
}

// TestSPAHandlerNeverLists checks the embedded FS cannot be walked even when no
// build assets are present (a clean checkout), where the handler degrades to a
// bare 404. The populated case is covered by TestSPAServesEmbeddedAssets.
func TestSPAHandlerNeverLists(t *testing.T) {
	rec := httptest.NewRecorder()
	jasprSPAHandler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, spaPrefix+"/", nil))

	if body := rec.Body.String(); strings.Contains(body, "Index of") {
		t.Errorf("SPA served a directory listing: %.200q", body)
	}
}
