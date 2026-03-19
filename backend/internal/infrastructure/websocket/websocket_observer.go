package websocket

import (
	"encoding/json"
	"strings"

	domain "backend/internal/domain/events"
)

// Minimal bridge EventBus -> WebSocket room broadcast.
type WebSocketObserver struct {
	hub *Hub
}

func NewWebSocketObserver(hub *Hub) *WebSocketObserver {
	return &WebSocketObserver{hub: hub}
}

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
