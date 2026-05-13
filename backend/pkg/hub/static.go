package hub

import (
	"bytes"
	"embed"
	"errors"
	"io/fs"
	"net/http"
	"strings"
	"time"
)

// StaticFS holds the embedded frontend assets
//
//go:embed frontend/*
var StaticFS embed.FS

func spaHandler() (http.Handler, error) {
	shellFS, err := fs.Sub(StaticFS, "frontend/shell/browser")
	if err != nil {
		return nil, err
	}

	hubappFS, err := fs.Sub(StaticFS, "frontend/hubapp/browser")
	if err != nil {
		return nil, err
	}

	shellFileServer := http.FileServer(http.FS(shellFS))
	hubappFileServer := http.StripPrefix("/hubapp/", http.FileServer(http.FS(hubappFS)))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/hubapp/") {
			hubappFileServer.ServeHTTP(w, r)
			return
		}

		path := strings.TrimPrefix(r.URL.Path, "/")
		if path == "" {
			path = "index.html"
		}

		info, err := fs.Stat(shellFS, path)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				serveShellIndex(w, r, shellFS)
			} else {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
			return
		}

		if info.IsDir() || path == "index.html" {
			serveShellIndex(w, r, shellFS)
			return
		}

		r.URL.Path = "/" + path
		shellFileServer.ServeHTTP(w, r)
	}), nil
}

func serveShellIndex(w http.ResponseWriter, r *http.Request, shellFS fs.FS) {
	data, err := fs.ReadFile(shellFS, "index.html")
	if err != nil {
		http.Error(w, "Frontend not found", http.StatusInternalServerError)
		return
	}

	http.ServeContent(w, r, "index.html", time.Time{}, bytes.NewReader(data))
}
