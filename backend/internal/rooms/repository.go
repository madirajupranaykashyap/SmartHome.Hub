package rooms

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/google/uuid"
)

var (
	ErrRoomNotFound     = errors.New("room not found")
	ErrUserRoomExists   = errors.New("room already added")
	ErrUsernameNotFound = errors.New("user not found")
)

type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) ListRooms() ([]Room, error) {
	rows, err := r.DB.Query(`
		SELECT id, room_urn, name
		FROM rooms
		ORDER BY name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rooms := []Room{}
	for rows.Next() {
		var room Room
		if err := rows.Scan(&room.ID, &room.RoomURN, &room.Name); err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}

	return rooms, rows.Err()
}

func (r *Repository) ListUserRooms(username string) ([]UserRoom, error) {
	rows, err := r.DB.Query(`
		SELECT ur.id, ur.room_guid, rooms.name, ur.created_at
		FROM user_rooms ur
		JOIN users ON users.id = ur.user_id
		JOIN rooms ON rooms.id = ur.room_guid
		WHERE users.username = ?
		ORDER BY ur.created_at, rooms.name
	`, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	userRooms := []UserRoom{}
	for rows.Next() {
		var room UserRoom
		if err := rows.Scan(&room.ID, &room.RoomGUID, &room.Room, &room.CreatedAt); err != nil {
			return nil, err
		}
		room.AssignedAt = room.CreatedAt
		userRooms = append(userRooms, room)
	}

	return userRooms, rows.Err()
}

func (r *Repository) CreateUserRoom(username string, roomGUID string) (*UserRoom, error) {
	roomGUID = strings.TrimSpace(roomGUID)
	if roomGUID == "" {
		return nil, ErrRoomNotFound
	}

	userID, err := r.userID(username)
	if err != nil {
		return nil, err
	}

	if ok, err := r.roomExists(roomGUID); err != nil {
		return nil, err
	} else if !ok {
		return nil, ErrRoomNotFound
	}

	id := uuid.NewString()
	_, err = r.DB.Exec(`
		INSERT INTO user_rooms (id, user_id, room_guid)
		VALUES (?, ?, ?)
	`, id, userID, roomGUID)
	if err != nil {
		if strings.Contains(err.Error(), "constraint failed") || strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return nil, ErrUserRoomExists
		}
		return nil, err
	}

	return r.GetUserRoom(username, id)
}

func (r *Repository) DeleteUserRoom(username string, id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return ErrRoomNotFound
	}

	userID, err := r.userID(username)
	if err != nil {
		return err
	}

	result, err := r.DB.Exec(`
		DELETE FROM user_rooms
		WHERE id = ? AND user_id = ?
	`, id, userID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrRoomNotFound
	}

	return nil
}

func (r *Repository) GetUserRoom(username string, id string) (*UserRoom, error) {
	var room UserRoom
	err := r.DB.QueryRow(`
		SELECT ur.id, ur.room_guid, rooms.name, ur.created_at
		FROM user_rooms ur
		JOIN users ON users.id = ur.user_id
		JOIN rooms ON rooms.id = ur.room_guid
		WHERE users.username = ? AND ur.id = ?
	`, username, id).Scan(&room.ID, &room.RoomGUID, &room.Room, &room.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRoomNotFound
		}
		return nil, err
	}

	room.AssignedAt = room.CreatedAt
	return &room, nil
}

func (r *Repository) userID(username string) (int64, error) {
	var id int64
	err := r.DB.QueryRow(`
		SELECT id
		FROM users
		WHERE username = ?
	`, username).Scan(&id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrUsernameNotFound
		}
		return 0, err
	}

	return id, nil
}

func (r *Repository) roomExists(roomGUID string) (bool, error) {
	var exists bool
	err := r.DB.QueryRow(`
		SELECT EXISTS(
			SELECT 1
			FROM rooms
			WHERE id = ?
		)
	`, roomGUID).Scan(&exists)

	return exists, err
}
