package hub

import (
	"io/fs"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSPAHandlerServesShellRoot(t *testing.T) {
	handler, err := spaHandler()
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	handler.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", res.Code)
	}
	if res.Body.Len() == 0 {
		t.Fatal("expected shell index body")
	}
}

func TestSPAHandlerDoesNotRedirectIndex(t *testing.T) {
	handler, err := spaHandler()
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/index.html", nil)
	handler.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", res.Code)
	}
	if got := res.Header().Get("Cache-Control"); got != "no-store" {
		t.Fatalf("expected Cache-Control no-store, got %q", got)
	}
}

func TestSPAHandlerServesHubappRemote(t *testing.T) {
	handler, err := spaHandler()
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/hubapp/remoteEntry.json", nil)
	handler.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", res.Code)
	}
}

func TestSPAHandlerServesUnderscorePrefixedChunks(t *testing.T) {
	handler, err := spaHandler()
	if err != nil {
		t.Fatal(err)
	}

	shellChunks, err := fs.Glob(StaticFS, "frontend/shell/browser/_angular_*.js")
	if err != nil {
		t.Fatal(err)
	}
	hubappChunks, err := fs.Glob(StaticFS, "frontend/hubapp/browser/_angular_*.js")
	if err != nil {
		t.Fatal(err)
	}
	if len(shellChunks) == 0 {
		t.Fatal("expected embedded shell angular chunks")
	}
	if len(hubappChunks) == 0 {
		t.Fatal("expected embedded hubapp angular chunks")
	}

	tests := []string{
		"/" + shellChunks[0][len("frontend/shell/browser/"):],
		"/hubapp/" + hubappChunks[0][len("frontend/hubapp/browser/"):],
	}

	for _, path := range tests {
		t.Run(path, func(t *testing.T) {
			res := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, path, nil)
			handler.ServeHTTP(res, req)

			if res.Code != http.StatusOK {
				t.Fatalf("expected status 200, got %d", res.Code)
			}
		})
	}
}

func TestSPAHandlerDoesNotServeIndexForMissingAsset(t *testing.T) {
	handler, err := spaHandler()
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/missing-chunk.js", nil)
	handler.ServeHTTP(res, req)

	if res.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", res.Code)
	}
}
