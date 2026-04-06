package websocket

import (
	"backend/internal/domain/room"
	"strings"
	"sync"
)

// Hub manages all the active RoomHubs, indexed by roomID.
// It is a thread-safe singleton.
type Hub struct {
	mu    sync.RWMutex
	rooms map[string]*RoomHub
}

var (
	hubInstance *Hub
	onceHub     sync.Once
)

// GetHubInstance returns the singleton instance of the Hub.
func GetHubInstance() *Hub {
	onceHub.Do(func() {
		hubInstance = &Hub{
			rooms: make(map[string]*RoomHub),
		}
	})
	return hubInstance
}

// CreateRoomHub creates and registers a new RoomHub for the given room.
func (h *Hub) CreateRoomHub(room *room.Room) *RoomHub {
	if h == nil || room == nil {
		return nil
	}

	roomID := strings.TrimSpace(room.ID)
	if roomID == "" {
		return nil
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	roomHub := NewRoomHub(room)
	h.rooms[roomID] = roomHub
	return roomHub
}

// GetRoomHub returns the RoomHub associated with a room.
func (h *Hub) GetRoomHub(roomID string) *RoomHub {
	roomID = strings.TrimSpace(roomID)
	if roomID == "" {
		return nil
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	if room, ok := h.rooms[roomID]; ok {
		return room
	}

	return nil
}

// GetRoom returns the Room entity of a room.
func (h *Hub) GetRoom(roomID string) *room.Room {
	roomHub := h.GetRoomHub(roomID)
	if roomHub == nil {
		return nil
	}

	return roomHub.GetRoom()

}

// AddClient registers a client in an existing RoomHub.
func (h *Hub) AddClient(roomID string, client *Client) {
	if h == nil || client == nil {
		return
	}

	roomID = strings.TrimSpace(roomID)
	if roomID == "" {
		return
	}

	roomHub := h.GetRoomHub(roomID)

	if roomHub == nil {
		return
	}

	roomHub.AddClient(client)
}

// RemoveClient removes a client from the RoomHub. If the room becomes empty, it is removed.
func (h *Hub) RemoveClient(roomID string, client *Client) {
	if h == nil || client == nil {
		return
	}

	roomID = strings.TrimSpace(roomID)
	if roomID == "" {
		return
	}

	h.mu.RLock()
	room, ok := h.rooms[roomID]
	h.mu.RUnlock()
	if !ok {
		return
	}

	isEmpty := room.RemoveClient(client)
	if !isEmpty {
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()
	if current, exists := h.rooms[roomID]; exists && current == room && current.IsEmpty() {
		delete(h.rooms, roomID)
	}

	//Aqui trata se quando o user sai e fecha se a sala e altera se oos players do game
}

// BroadcastToRoom sends a payload to all clients in a room.
func (h *Hub) BroadcastToRoom(roomID string, payload []byte) {
	if h == nil || len(payload) == 0 {
		return
	}

	roomID = strings.TrimSpace(roomID)
	if roomID == "" {
		return
	}

	h.mu.RLock()
	room, ok := h.rooms[roomID]
	h.mu.RUnlock()
	if !ok {
		return
	}

	room.Broadcast(payload)
}
