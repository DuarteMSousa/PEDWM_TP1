package websocket_interfaces

import "backend/internal/domain/room"

// RoomService defines the contract for room service used by the WebSocket layer.
type RoomService interface {
	LeaveRoom(roomID, userID string) (*room.Room, error)
}
