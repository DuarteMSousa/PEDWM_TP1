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

const lobbyRoomID = "lobby"

var (
	hubInstance *Hub
	onceHub     sync.Once
)

func GetHubInstance() *Hub {
	onceHub.Do(func() {
		hubInstance = &Hub{
			rooms: map[string]*RoomHub{
				lobbyRoomID: NewRoomHub(nil),
			},
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

func (h *Hub) AddClient(roomID string, client *Client) {
	if h == nil || client == nil {
		return
	}

	roomID = strings.TrimSpace(roomID)
	if roomID == "" {
		return
	}

	roomHub := h.GetRoomHub(roomID)
	if roomHub == nil && roomID == lobbyRoomID {
		h.mu.Lock()
		roomHub = h.rooms[lobbyRoomID]
		if roomHub == nil {
			roomHub = NewRoomHub(nil)
			h.rooms[lobbyRoomID] = roomHub
		}
		h.mu.Unlock()
	}

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
	if roomID == lobbyRoomID {
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
