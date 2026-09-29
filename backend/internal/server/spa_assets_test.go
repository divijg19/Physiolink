package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/divijg19/physiolink/backend/internal/config"
)

// TestSPAServesEmbeddedAssets exercises the /site mount against real build
// output. On a clean checkout the embed directory holds only placeholder.txt
// and the handler degrades to 404, so the assertions are skipped there; run
// `make build-web` first to exercise them.
//
// The <base href="/site/"> marker is what makes jaspr_router's Link() and the
// SPA's relative asset URLs resolve under the mount prefix.
func TestSPAServesEmbeddedAssets(t *testing.T) {
	h := jasprSPAHandler()

	probe := httptest.NewRecorder()
	h.ServeHTTP(probe, httptest.NewRequest(http.MethodGet, spaPrefix+"/", nil))
	if probe.Code != http.StatusOK {
		t.Skipf("SPA build output not embedded (%d); run `make build-web`", probe.Code)
	}

	cases := []struct {
		name         string
		path         string
		wantContains string
	}{
		{"root serves the shell", spaPrefix + "/", `<base href="/site/">`},
		{"client route falls back to the shell", spaPrefix + "/about", `<base href="/site/">`},
		{"deep client route falls back to the shell", spaPrefix + "/a/b/c", `<base href="/site/">`},
		{"js bundle is served", spaPrefix + "/main.dart.js", "dartProgram"},
		{"directory does not list", spaPrefix + "/icons", `<base href="/site/">`},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, tc.path, nil))

			if rec.Code != http.StatusOK {
				t.Fatalf("GET %s = %d, want 200 (body %.120q)", tc.path, rec.Code, rec.Body.String())
			}
			// http.FileServer 301-canonicalises any /index.html path; a redirect
			// here would bounce every SPA route.
			if loc := rec.Header().Get("Location"); loc != "" {
				t.Fatalf("GET %s redirected to %q; SPA routes must not redirect", tc.path, loc)
			}
			body := rec.Body.String()
			if tc.wantContains != "" && !strings.Contains(body, tc.wantContains) {
				t.Errorf("GET %s body missing %q; got %.160q", tc.path, tc.wantContains, body)
			}
			if strings.Contains(body, "Index of") {
				t.Errorf("GET %s produced a directory listing", tc.path)
			}
		})
	}
}

// TestSPAAndPortalAreSeparate proves the two front ends no longer collide: the
// SPA answers under /site while the server-rendered portal keeps the root.
func TestSPAAndPortalAreSeparate(t *testing.T) {
	router := NewRouter(&config.Config{JWTSecret: "test-secret"})

	spaProbe := httptest.NewRecorder()
	router.ServeHTTP(spaProbe, httptest.NewRequest(http.MethodGet, spaPrefix+"/", nil))
	if spaProbe.Code != http.StatusOK {
		t.Skipf("SPA build output not embedded (%d); run `make build-web`", spaProbe.Code)
	}
	if !strings.Contains(spaProbe.Body.String(), `<base href="/site/">`) {
		t.Error("GET /site/ did not serve the SPA shell")
	}

	// The portal pages must still be reachable at the root and must NOT be the
	// SPA (they are server-rendered forms).
	for _, path := range []string{"/login", "/register"} {
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
		if rec.Code != http.StatusOK {
			t.Errorf("GET %s = %d, want 200", path, rec.Code)
			continue
		}
		body := rec.Body.String()
		if strings.Contains(body, `<base href="/site/">`) {
			t.Errorf("GET %s served the SPA; the portal must own the root", path)
		}
	}
}
