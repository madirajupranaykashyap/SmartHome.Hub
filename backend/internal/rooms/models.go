package rooms

type Room struct {
	ID      string `json:"id"`
	RoomURN string `json:"roomUrn"`
	Name    string `json:"name"`
}

type UserRoom struct {
	ID         string `json:"id"`
	RoomGUID   string `json:"room_guid"`
	Room       string `json:"room"`
	CreatedAt  string `json:"created_at"`
	AssignedAt string `json:"assigned_at"`
}

type CreateUserRoomRequest struct {
	ID       string `json:"id"`
	RoomGUID string `json:"room_guid"`
	RoomGuid string `json:"roomGuid"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
