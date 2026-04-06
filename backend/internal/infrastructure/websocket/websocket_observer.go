package websocket

import (
	"encoding/json"
	"strings"

	domain "backend/internal/domain/events"
)

// WebSocketObserver implements Observer and broadcasts domain events
// to the respective WebSocket clients in the room.
type WebSocketObserver struct {
	hub *Hub
}

// NewWebSocketObserver creates a new WebSocket observer.
func NewWebSocketObserver(hub *Hub) *WebSocketObserver {
	return &WebSocketObserver{hub: hub}
}

// Update receives a domain event and broadcasts it to the room.
func (o *WebSocketObserver) Update(event domain.Event) {
	if o == nil || o.hub == nil {
		return
	}

	roomID := strings.TrimSpace(event.RoomID)

	if roomID == "" {
		return
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return
	}

	o.hub.BroadcastToRoom(roomID, payload)
}
