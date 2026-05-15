package hub

import (
	"bytes"
	"embed"
	"io/fs"
	"net/http"
	"strings"
	"time"
)

//go:embed all:frontend
var StaticFS embed.FS

func spaHandler() (http.Handler, error) {

	// Shell frontend
	shellFS, err := fs.Sub(StaticFS, "frontend/shell/browser")
	if err != nil {
		return nil, err
	}

	// Hubapp frontend
	hubappFS, err := fs.Sub(StaticFS, "frontend/hubapp/browser")
	if err != nil {
		return nil, err
	}

	shellFileServer := http.FileServer(http.FS(shellFS))
	hubappFileServer := http.FileServer(http.FS(hubappFS))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// =========================================================
		// HUBAPP REMOTE ROUTES
		// =========================================================

		if strings.HasPrefix(r.URL.Path, "/hubapp/") {

			req := r.Clone(r.Context())

			req.URL.Path = strings.TrimPrefix(r.URL.Path, "/hubapp")

			if req.URL.Path == "" {
				req.URL.Path = "/"
			}

			hubappFileServer.ServeHTTP(w, req)
			return
		}

		// =========================================================
		// NORMAL SHELL ROUTES
		// =========================================================

		path := strings.TrimPrefix(r.URL.Path, "/")

		// Root SPA entry
		if path == "" || path == "index.html" {
			serveShellIndex(w, r, shellFS)
			return
		}

		// =========================================================
		// SHELL STATIC ASSETS
		// =========================================================

		if _, err := fs.Stat(shellFS, path); err == nil {

			req := r.Clone(r.Context())
			req.URL.Path = "/" + path

			shellFileServer.ServeHTTP(w, req)
			return
		}

		// =========================================================
		// HUBAPP STATIC ASSETS
		// =========================================================

		if _, err := fs.Stat(hubappFS, path); err == nil {

			req := r.Clone(r.Context())
			req.URL.Path = "/" + path

			hubappFileServer.ServeHTTP(w, req)
			return
		}

		// =========================================================
		// MISSING STATIC FILE
		// =========================================================

		if isAssetRequest(path) {
			http.NotFound(w, r)
			return
		}

		// =========================================================
		// ANGULAR SPA FALLBACK
		// =========================================================

		serveShellIndex(w, r, shellFS)

	}), nil
}

func isAssetRequest(path string) bool {

	parts := strings.Split(path, "/")
	name := parts[len(parts)-1]

	return strings.Contains(name, ".")
}

func serveShellIndex(
	w http.ResponseWriter,
	r *http.Request,
	shellFS fs.FS,
) {

	data, err := fs.ReadFile(shellFS, "index.html")
	if err != nil {
		http.Error(w, "Frontend not found", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")

	http.ServeContent(
		w,
		r,
		"index.html",
		time.Time{},
		bytes.NewReader(data),
	)
}
