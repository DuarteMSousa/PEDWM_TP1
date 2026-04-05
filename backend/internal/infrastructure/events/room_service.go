package events_infrastructure

import "backend/internal/domain/room"

type RoomService interface {
	LeaveRoom(roomID, userID string) (*room.Room, error)
	DeleteRoom(roomID string) error
}
