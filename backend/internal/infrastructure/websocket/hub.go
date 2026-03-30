package websocket

import (
	"strings"
	"sync"
)

type Hub struct {
	mu    sync.RWMutex
	rooms map[string]*RoomHub
}

var (
	instance *Hub
	once     sync.Once
)

func GetInstance() *Hub {
	once.Do(func() {
		instance = &Hub{
			rooms: make(map[string]*RoomHub),
		}
	})
	return instance
}

func (h *Hub) GetOrCreateRoom(roomID string) *RoomHub {
	roomID = strings.TrimSpace(roomID)
	if roomID == "" {
		return nil
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	if room, ok := h.rooms[roomID]; ok {
		return room
	}

	room := NewRoomHub(roomID)
	h.rooms[roomID] = room
	return room
}

func (h *Hub) AddClient(roomID string, client *Client) {
	if h == nil || client == nil {
		return
	}

	room := h.GetOrCreateRoom(roomID)
	if room == nil {
		return
	}

	room.AddClient(client)
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
