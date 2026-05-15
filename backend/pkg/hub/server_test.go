package hub

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"smarthome/hub/core/logger"
	"smarthome/hub/internal/database"
	"smarthome/hub/pkg/updater"
)

func TestUpdateCheckSkipReasonUsesLauncherCheck(t *testing.T) {
	server := &Server{
		Config: Config{
			EnableUpdateCheck:      true,
			SkipStartupUpdateCheck: true,
		},
		UpdateChecker: &updater.Checker{},
	}

	if got := server.updateCheckSkipReason(); got != "already checked by launcher" {
		t.Fatalf("expected launcher skip reason, got %q", got)
	}
}

func TestRoomsAPIRoundTrip(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret")
	logger.Init("test")

	db := openTestDB(t)
	defer db.Close()

	router := NewRouter(db, Config{
		AllowedOrigins: []string{"*"},
	}, nil)

	token := registerTestUser(t, router)

	catalog := doJSONRequest(t, router, http.MethodGet, "/api/rooms/catalog", nil, token)
	if catalog.Code != http.StatusOK {
		t.Fatalf("expected catalog status 200, got %d: %s", catalog.Code, catalog.Body.String())
	}

	var rooms []struct {
		ID      string `json:"id"`
		RoomURN string `json:"roomUrn"`
		Name    string `json:"name"`
	}
	if err := json.Unmarshal(catalog.Body.Bytes(), &rooms); err != nil {
		t.Fatal(err)
	}
	if len(rooms) != 20 {
		t.Fatalf("expected 20 seeded rooms, got %d", len(rooms))
	}

	created := doJSONRequest(
		t,
		router,
		http.MethodPost,
		"/api/rooms",
		map[string]string{"room_guid": "1f1fe586-afc6-447b-ba8b-a90bba9be12a"},
		token,
	)
	if created.Code != http.StatusCreated {
		t.Fatalf("expected create status 201, got %d: %s", created.Code, created.Body.String())
	}

	list := doJSONRequest(t, router, http.MethodGet, "/api/rooms", nil, token)
	if list.Code != http.StatusOK {
		t.Fatalf("expected list status 200, got %d: %s", list.Code, list.Body.String())
	}

	var userRooms []struct {
		ID        string `json:"id"`
		RoomGUID  string `json:"room_guid"`
		Room      string `json:"room"`
		CreatedAt string `json:"created_at"`
	}
	if err := json.Unmarshal(list.Body.Bytes(), &userRooms); err != nil {
		t.Fatal(err)
	}
	if len(userRooms) != 1 {
		t.Fatalf("expected 1 user room, got %d", len(userRooms))
	}
	if userRooms[0].Room != "Art Room" {
		t.Fatalf("expected Art Room, got %q", userRooms[0].Room)
	}

	deleted := doJSONRequest(
		t,
		router,
		http.MethodDelete,
		"/api/rooms/"+userRooms[0].ID,
		nil,
		token,
	)
	if deleted.Code != http.StatusNoContent {
		t.Fatalf("expected delete status 204, got %d: %s", deleted.Code, deleted.Body.String())
	}

	listAfterDelete := doJSONRequest(t, router, http.MethodGet, "/api/rooms", nil, token)
	if listAfterDelete.Code != http.StatusOK {
		t.Fatalf("expected list status 200 after delete, got %d: %s", listAfterDelete.Code, listAfterDelete.Body.String())
	}

	var userRoomsAfterDelete []struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(listAfterDelete.Body.Bytes(), &userRoomsAfterDelete); err != nil {
		t.Fatal(err)
	}
	if len(userRoomsAfterDelete) != 0 {
		t.Fatalf("expected 0 user rooms after delete, got %d", len(userRoomsAfterDelete))
	}
}

func openTestDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := database.Open(filepath.Join(t.TempDir(), "hub.db"))
	if err != nil {
		t.Fatal(err)
	}

	return db
}

func registerTestUser(t *testing.T, router http.Handler) string {
	t.Helper()

	res := doJSONRequest(
		t,
		router,
		http.MethodPost,
		"/api/auth/register",
		map[string]string{
			"username": "rooms-user",
			"password": "password123",
		},
		"",
	)
	if res.Code != http.StatusCreated {
		t.Fatalf("expected register status 201, got %d: %s", res.Code, res.Body.String())
	}

	var body struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(res.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Token == "" {
		t.Fatal("expected auth token")
	}

	return body.Token
}

func doJSONRequest(
	t *testing.T,
	router http.Handler,
	method string,
	path string,
	body interface{},
	token string,
) *httptest.ResponseRecorder {
	t.Helper()

	var requestBody *bytes.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}
		requestBody = bytes.NewReader(data)
	} else {
		requestBody = bytes.NewReader(nil)
	}

	req := httptest.NewRequest(method, path, requestBody)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)
	return res
}
