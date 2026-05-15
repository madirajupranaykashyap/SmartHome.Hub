package rooms

import (
	"encoding/json"
	"errors"
	"net/http"

	"smarthome/hub/internal/auth"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	Service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{Service: service}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/api/rooms", func(r chi.Router) {
		r.Use(auth.JWTMiddleware)
		r.Get("/", h.ListUserRooms)
		r.Post("/", h.CreateUserRoom)
		r.Delete("/{id}", h.DeleteUserRoom)
		r.Get("/catalog", h.ListRooms)
	})
}

// ListRooms godoc
//
// @Summary List predefined rooms
// @Description Return the predefined room catalog a user can choose from
// @Tags rooms
// @Security BearerAuth
// @Produce json
// @Success 200 {array} Room
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/rooms/catalog [get]
func (h *Handler) ListRooms(w http.ResponseWriter, r *http.Request) {
	rooms, err := h.Service.ListRooms()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load rooms")
		return
	}

	writeJSON(w, http.StatusOK, rooms)
}

// ListUserRooms godoc
//
// @Summary List user rooms
// @Description Return rooms added to the authenticated user's room list
// @Tags rooms
// @Security BearerAuth
// @Produce json
// @Success 200 {array} UserRoom
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/rooms [get]
func (h *Handler) ListUserRooms(w http.ResponseWriter, r *http.Request) {
	username, ok := usernameFromRequest(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing user")
		return
	}

	rooms, err := h.Service.ListUserRooms(username)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load user rooms")
		return
	}

	writeJSON(w, http.StatusOK, rooms)
}

// CreateUserRoom godoc
//
// @Summary Add user room
// @Description Add a predefined room to the authenticated user's room list
// @Tags rooms
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body CreateUserRoomRequest true "Room to add"
// @Success 201 {object} UserRoom
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/rooms [post]
func (h *Handler) CreateUserRoom(w http.ResponseWriter, r *http.Request) {
	username, ok := usernameFromRequest(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing user")
		return
	}

	var req CreateUserRoomRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	room, err := h.Service.CreateUserRoom(username, req)
	if err != nil {
		switch {
		case errors.Is(err, ErrRoomNotFound):
			writeError(w, http.StatusBadRequest, "room not found")
		case errors.Is(err, ErrUserRoomExists):
			writeError(w, http.StatusConflict, "room already added")
		case errors.Is(err, ErrUsernameNotFound):
			writeError(w, http.StatusUnauthorized, "user not found")
		default:
			writeError(w, http.StatusInternalServerError, "failed to add room")
		}
		return
	}

	writeJSON(w, http.StatusCreated, room)
}

// DeleteUserRoom godoc
//
// @Summary Delete user room
// @Description Remove a room from the authenticated user's room list
// @Tags rooms
// @Security BearerAuth
// @Produce json
// @Param id path string true "User room id"
// @Success 204
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/rooms/{id} [delete]
func (h *Handler) DeleteUserRoom(w http.ResponseWriter, r *http.Request) {
	username, ok := usernameFromRequest(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing user")
		return
	}

	err := h.Service.DeleteUserRoom(username, chi.URLParam(r, "id"))
	if err != nil {
		switch {
		case errors.Is(err, ErrRoomNotFound):
			writeError(w, http.StatusNotFound, "room not found")
		case errors.Is(err, ErrUsernameNotFound):
			writeError(w, http.StatusUnauthorized, "user not found")
		default:
			writeError(w, http.StatusInternalServerError, "failed to delete room")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func usernameFromRequest(r *http.Request) (string, bool) {
	username, ok := r.Context().Value(auth.UserContextKey).(string)
	return username, ok && username != ""
}

func decodeJSON(w http.ResponseWriter, r *http.Request, dst interface{}) bool {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dst); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return false
	}

	return true
}

func writeJSON(w http.ResponseWriter, status int, value interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, ErrorResponse{Error: message})
}
