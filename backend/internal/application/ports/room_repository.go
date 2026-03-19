package ports

type RoomRepository interface {
	CreateRoom(
		name string,
		hostPlayerID string,
		maxPlayers int,
		isPrivate bool,
		password string,
	) (Room, error)
	DeleteRoom(roomID string, requesterID string) (Room, error)
	JoinRoom(roomID string, playerID string, password string) (RoomView, Room, Player, error)
	LeaveRoom(roomID string, playerID string) (RoomView, Room, error)
	GetRoom(roomID string) (Room, bool)
	ListRoomsDetailed() []Room
	ListRoomViews() []RoomView
}
