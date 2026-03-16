package observer

import (
	"encoding/json"
	"strings"

	ws "backend/internal/infrastructure/transport/websocket"
	domain "backend/internal/model"
)

// Minimal bridge EventBus -> WebSocket room broadcast.
// TODO(team-eventbus): evolve payload mapping rules with domain team.
type WebSocketObserver struct {
	hub *ws.Hub
}

func NewWebSocketObserver(hub *ws.Hub) *WebSocketObserver {
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
