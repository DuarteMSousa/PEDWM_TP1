package websocket

import (
	"backend/internal/domain/room"
	"strings"
	"sync"
)

type Hub struct {
	mu    sync.RWMutex
	rooms map[string]*RoomHub
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

func (h *Hub) CreateRoomHub(roomID string, hostId string, hostUserName string) *RoomHub {
	roomID = strings.TrimSpace(roomID)
	if roomID == "" {
		return nil
	}

	h.mu.Lock()
	defer h.mu.Unlock()
	room, err := room.NewRoom(roomID, hostId, hostUserName)
	if err != nil {
		return nil
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

	h.mu.Lock()
	defer h.mu.Unlock()

	if room, ok := h.rooms[roomID]; ok {
		return room
	}

	return nil
}

func (h *Hub) AddClient(roomID string, client *Client) {
	if h == nil || client == nil {
		return
	}

	roomHub := h.GetRoomHub(roomID)
	if roomHub == nil {
		return
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
