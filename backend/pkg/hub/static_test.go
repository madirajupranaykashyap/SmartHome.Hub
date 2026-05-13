package hub

import (
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
