package repository

import (
	"backend/internal/domain"
	"sync"
)

type RoomRepository struct {
	mu    sync.RWMutex
	rooms map[string]*domain.Room
}

func NewRoomRepository() *RoomRepository {
	return &RoomRepository{
		rooms: make(map[string]*domain.Room),
	}
}

func (r *RoomRepository) Save(room *domain.Room) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.rooms[room.ID] = room
}

func (r *RoomRepository) FindByID(id string) (*domain.Room, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	room, ok := r.rooms[id]
	return room, ok
}

func (r *RoomRepository) FindAll() []*domain.Room {
	r.mu.RLock()
	defer r.mu.RUnlock()

	rooms := make([]*domain.Room, 0, len(r.rooms))
	for _, r := range r.rooms {
		rooms = append(rooms, r)
	}

	return rooms
}

func (r *RoomRepository) Delete(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.rooms, id)
}
