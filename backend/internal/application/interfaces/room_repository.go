package interfaces

import "backend/internal/domain/room"

// RoomRepository defines the contract for room persistence.
type RoomRepository interface {
	Save(room *room.Room) error
	FindByID(id string) (*room.Room, error)
	FindAll() ([]*room.Room, error)
	Delete(id string) error
}
