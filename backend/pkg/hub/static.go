package hub

import (
	"embed"
	"io/fs"
	"net/http"
)

// StaticFS holds the embedded frontend assets
//
//go:embed frontend
var StaticFS embed.FS

// ServeStatic serves the embedded frontend files
func serveStatic(mux *http.ServeMux) error {
	// Extract the frontend subdirectory from the embed
	frontendFS, err := fs.Sub(StaticFS, "frontend")
	if err != nil {
		return err
	}

	// Serve static files
	mux.Handle("GET /", http.FileServer(http.FS(frontendFS)))

	return nil
}
