package server

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"
)

//go:embed jaspr
var jasprStatic embed.FS

// spaPrefix is the path the marketing SPA is mounted at. Keeping it off the
// server root stops the SPA fallback from shadowing the server-rendered
// portal routes (/login, /register, /therapists, /dashboard, ...).
const spaPrefix = "/site"

// jasprSPAHandler serves the embedded Jaspr SPA under spaPrefix.
//
// Unknown paths fall back to index.html so client-side routes survive a hard
// refresh, mirroring how a static host would behave.
func jasprSPAHandler() http.Handler {
	sub, err := fs.Sub(jasprStatic, "jaspr")
	if err != nil {
		return http.NotFoundHandler()
	}

	if _, err := sub.Open("index.html"); err != nil {
		return http.NotFoundHandler()
	}

	fileServer := http.FileServer(http.FS(sub))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rest := strings.TrimPrefix(r.URL.Path, spaPrefix)
		if rest == "" {
			rest = "/"
		}

		name := strings.TrimPrefix(rest, "/")
		serveShell := name == "" || name == "index.html"
		if !serveShell {
			// Only real files are served directly. Directories fall through to
			// the shell so the embed cannot be walked as a listing.
			if f, err := sub.Open(name); err == nil {
				info, statErr := f.Stat()
				f.Close()
				if statErr != nil || info.IsDir() {
					serveShell = true
				}
			} else {
				serveShell = true
			}
		}
		if serveShell {
			// Serve the directory rather than "/index.html": http.FileServer
			// 301-redirects any path ending in /index.html to canonicalise it,
			// which would bounce every SPA route.
			rest = "/"
		}

		// http.FileServer routes on r.URL.Path, so give it a request whose
		// path has the mount prefix removed rather than mutating the original.
		r2 := r.Clone(r.Context())
		r2.URL.Path = rest
		fileServer.ServeHTTP(w, r2)
	})
}
