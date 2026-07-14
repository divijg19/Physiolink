package server

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"
)

//go:embed jaspr
var jasprStatic embed.FS

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
		path := strings.TrimPrefix(r.URL.Path, "/")
		if f, err := sub.Open(path); err == nil {
			f.Close()
			fileServer.ServeHTTP(w, r)
			return
		}
		r.URL.Path = "/index.html"
		fileServer.ServeHTTP(w, r)
	})
}
