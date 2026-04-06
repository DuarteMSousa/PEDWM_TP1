package events_infrastructure

import "backend/internal/domain/room"

// RoomService defines the interface for interacting with the room service.
type RoomService interface {
	LeaveRoom(roomID, userID string) (*room.Room, error)
	DeleteRoom(roomID string) error
}
