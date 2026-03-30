package websocket

import (
	"encoding/json"
	"net/http"
	"strings"

	gws "github.com/gorilla/websocket"
)

type Handler struct {
	Hub        *Hub
	Upgrader   gws.Upgrader
	Dispatcher *CommandDispatcher
}

func NewHandler(hub *Hub, dispatcher *CommandDispatcher) *Handler {
	return &Handler{
		Hub:        hub,
		Dispatcher: dispatcher,
		Upgrader: gws.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				origin := r.Header.Get("Origin")
				if origin == "" {
					return true
				}
				host := r.Host
				return strings.Contains(origin, host) || strings.Contains(origin, "localhost")
			},
		},
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h == nil || h.Hub == nil {
		http.Error(w, "websocket hub not configured", http.StatusInternalServerError)
		return
	}

	roomID := strings.TrimSpace(r.URL.Query().Get("roomId"))
	if roomID == "" {
		roomID = "lobby"
	}

	playerID := strings.TrimSpace(r.URL.Query().Get("playerId"))
	if playerID == "" {
		playerID = r.RemoteAddr
	}

	conn, err := h.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	client := NewClient(playerID, roomID, conn, h.Hub, h.Dispatcher)
	h.Hub.AddClient(roomID, client)

	go client.WritePump()
	go client.ReadPump()
}

func writeJSONError(w http.ResponseWriter, err error, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{"error": err.Error()})
}
