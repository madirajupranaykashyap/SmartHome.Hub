package rooms

import "strings"

type Service struct {
	Repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{Repo: repo}
}

func (s *Service) ListRooms() ([]Room, error) {
	return s.Repo.ListRooms()
}

func (s *Service) ListUserRooms(username string) ([]UserRoom, error) {
	return s.Repo.ListUserRooms(strings.TrimSpace(username))
}

func (s *Service) CreateUserRoom(username string, req CreateUserRoomRequest) (*UserRoom, error) {
	roomGUID := req.RoomGUID
	if roomGUID == "" {
		roomGUID = req.RoomGuid
	}
	if roomGUID == "" {
		roomGUID = req.ID
	}

	return s.Repo.CreateUserRoom(
		strings.TrimSpace(username),
		roomGUID,
	)
}

func (s *Service) DeleteUserRoom(username string, id string) error {
	return s.Repo.DeleteUserRoom(
		strings.TrimSpace(username),
		id,
	)
}
