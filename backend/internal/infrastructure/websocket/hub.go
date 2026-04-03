package websocket

import (
	"backend/internal/domain/room"
	"strings"
	"sync"
)

type Hub struct {
	mu                   sync.RWMutex
	rooms                map[string]*RoomHub
	onClientDisconnected func(roomID string, playerID string)
}

var (
	hubInstance *Hub
	onceHub     sync.Once
)

func GetHubInstance() *Hub {
	onceHub.Do(func() {
		hubInstance = &Hub{
			rooms: make(map[string]*RoomHub),
		}
	})
	return hubInstance
}

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

	if existing, ok := h.rooms[roomID]; ok {
		existing.SetRoom(room)
		return existing
	}

	roomHub := NewRoomHub(room)
	h.rooms[roomID] = roomHub
	return roomHub
}

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

func (h *Hub) GetRoom(roomID string) *room.Room {
	roomHub := h.GetRoomHub(roomID)
	if roomHub == nil {
		return nil
	}

	return roomHub.GetRoom()

}

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
		h.mu.Lock()
		roomHub = h.rooms[roomID]
		if roomHub == nil {
			roomHub = NewRoomHub(nil)
			h.rooms[roomID] = roomHub
		}
		h.mu.Unlock()
	}

	roomHub.AddClient(client)
}

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
	disconnectHandler := h.onClientDisconnected
	h.mu.RUnlock()
	if !ok {
		return
	}

	room.RemoveClient(client)

	if disconnectHandler != nil {
		go disconnectHandler(roomID, client.id)
	}

	if room.IsEmpty() && !room.HasRoom() {
		h.mu.Lock()
		if current, exists := h.rooms[roomID]; exists && current == room && current.IsEmpty() && !current.HasRoom() {
			delete(h.rooms, roomID)
		}
		h.mu.Unlock()
	}
}

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

func (h *Hub) DeleteRoomHub(roomID string) {
	if h == nil {
		return
	}

	roomID = strings.TrimSpace(roomID)
	if roomID == "" {
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.rooms, roomID)
}

func (h *Hub) SetDisconnectHandler(handler func(roomID string, playerID string)) {
	if h == nil {
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()
	h.onClientDisconnected = handler
}
