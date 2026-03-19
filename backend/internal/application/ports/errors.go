package ports

import "errors"

var (
	ErrNicknameRequired = errors.New("nickname is required")
	ErrInvalidPlayerID  = errors.New("player id is required")
	ErrPlayerNotFound   = errors.New("player not found")

	ErrInvalidRoomID        = errors.New("room id is required")
	ErrRoomNameRequired     = errors.New("room name is required")
	ErrRoomNotFound         = errors.New("room not found")
	ErrRoomNotOpen          = errors.New("cannot join/leave: room is not open")
	ErrRoomFull             = errors.New("room is full")
	ErrNotRoomHost          = errors.New("only the room host can delete this room")
	ErrRoomPasswordRequired = errors.New("password is required for private room")
	ErrRoomPasswordInvalid  = errors.New("invalid room password")
)
