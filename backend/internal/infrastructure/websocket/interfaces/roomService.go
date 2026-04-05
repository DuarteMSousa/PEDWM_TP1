package websocket_interfaces

import "backend/internal/domain/room"

type RoomService interface {
	LeaveRoom(roomID, userID string) (*room.Room, error)
}
