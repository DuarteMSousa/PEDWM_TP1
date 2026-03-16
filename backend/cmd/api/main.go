package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"

	wstransport "backend/internal/infrastructure/transport/websocket"
)

type testBroadcastRequest struct {
	RoomID string         `json:"roomId"`
	Event  map[string]any `json:"event"`
}

type roomMembershipRequest struct {
	PlayerID string `json:"playerId"`
}

type roomView struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	PlayersCount int    `json:"playersCount"`
	MaxPlayers   int    `json:"maxPlayers"`
	IsPrivate    bool   `json:"isPrivate"`
}

type roomRecord struct {
	ID         string
	Name       string
	MaxPlayers int
	IsPrivate  bool
	Players    map[string]struct{}
}

type lobbyStore struct {
	mu    sync.RWMutex
	rooms map[string]*roomRecord
}

var (
	errRoomNotFound = errors.New("room not found")
	errRoomFull     = errors.New("room is full")
)

func main() {
	hub := wstransport.NewHub()
	store := newLobbyStore()

	mux := http.NewServeMux()
	mux.Handle("/ws", wstransport.NewHandler(hub))
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/rooms", roomsHandler(store))
	mux.HandleFunc("/rooms/", roomActionHandler(store))
	mux.HandleFunc("/ws/broadcast", testBroadcastHandler(hub))

	addr := os.Getenv("API_ADDR")
	if strings.TrimSpace(addr) == "" {
		addr = ":4000"
	}

	log.Printf("api listening on %s", addr)
	if err := http.ListenAndServe(addr, withCORS(mux)); err != nil {
		log.Fatal(err)
	}
}

func newLobbyStore() *lobbyStore {
	return &lobbyStore{
		rooms: map[string]*roomRecord{
			"room_1": {
				ID:         "room_1",
				Name:       "Mesa 1",
				MaxPlayers: 4,
				IsPrivate:  false,
				Players: map[string]struct{}{
					"seed_a": {},
					"seed_b": {},
				},
			},
			"room_2": {
				ID:         "room_2",
				Name:       "Mesa 2",
				MaxPlayers: 4,
				IsPrivate:  false,
				Players: map[string]struct{}{
					"seed_1": {},
					"seed_2": {},
					"seed_3": {},
					"seed_4": {},
				},
			},
			"room_3": {
				ID:         "room_3",
				Name:       "Treino",
				MaxPlayers: 4,
				IsPrivate:  true,
				Players: map[string]struct{}{
					"seed_only": {},
				},
			},
		},
	}
}

func (s *lobbyStore) listRooms() []roomView {
	s.mu.RLock()
	defer s.mu.RUnlock()

	rooms := make([]roomView, 0, len(s.rooms))
	for _, room := range s.rooms {
		rooms = append(rooms, roomView{
			ID:           room.ID,
			Name:         room.Name,
			PlayersCount: len(room.Players),
			MaxPlayers:   room.MaxPlayers,
			IsPrivate:    room.IsPrivate,
		})
	}

	sort.Slice(rooms, func(i, j int) bool {
		return rooms[i].ID < rooms[j].ID
	})

	return rooms
}

func (s *lobbyStore) joinRoom(roomID string, playerID string) (roomView, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	room, ok := s.rooms[roomID]
	if !ok {
		return roomView{}, errRoomNotFound
	}
	if _, exists := room.Players[playerID]; exists {
		return roomView{
			ID:           room.ID,
			Name:         room.Name,
			PlayersCount: len(room.Players),
			MaxPlayers:   room.MaxPlayers,
			IsPrivate:    room.IsPrivate,
		}, nil
	}
	if len(room.Players) >= room.MaxPlayers {
		return roomView{}, errRoomFull
	}

	room.Players[playerID] = struct{}{}
	return roomView{
		ID:           room.ID,
		Name:         room.Name,
		PlayersCount: len(room.Players),
		MaxPlayers:   room.MaxPlayers,
		IsPrivate:    room.IsPrivate,
	}, nil
}

func (s *lobbyStore) leaveRoom(roomID string, playerID string) (roomView, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	room, ok := s.rooms[roomID]
	if !ok {
		return roomView{}, errRoomNotFound
	}
	delete(room.Players, playerID)

	return roomView{
		ID:           room.ID,
		Name:         room.Name,
		PlayersCount: len(room.Players),
		MaxPlayers:   room.MaxPlayers,
		IsPrivate:    room.IsPrivate,
	}, nil
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, errors.New("method not allowed"), http.StatusMethodNotAllowed)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"status": "ok"})
}

func roomsHandler(store *lobbyStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(w, errors.New("method not allowed"), http.StatusMethodNotAllowed)
			return
		}

		writeJSON(w, http.StatusOK, store.listRooms())
	}
}

func roomActionHandler(store *lobbyStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeError(w, errors.New("method not allowed"), http.StatusMethodNotAllowed)
			return
		}

		roomID, action := parseRoomActionPath(r.URL.Path)
		if roomID == "" || action == "" {
			writeError(w, errors.New("invalid room action path"), http.StatusBadRequest)
			return
		}

		var req roomMembershipRequest
		if err := decodeJSON(r, &req); err != nil {
			writeError(w, err, http.StatusBadRequest)
			return
		}

		req.PlayerID = strings.TrimSpace(req.PlayerID)
		if req.PlayerID == "" {
			writeError(w, errors.New("playerId is required"), http.StatusBadRequest)
			return
		}

		switch action {
		case "join":
			room, err := store.joinRoom(roomID, req.PlayerID)
			if err != nil {
				handleRoomStoreError(w, err)
				return
			}
			writeJSON(w, http.StatusOK, room)
		case "leave":
			room, err := store.leaveRoom(roomID, req.PlayerID)
			if err != nil {
				handleRoomStoreError(w, err)
				return
			}
			writeJSON(w, http.StatusOK, room)
		default:
			writeError(w, errors.New("unknown room action"), http.StatusNotFound)
		}
	}
}

func parseRoomActionPath(path string) (string, string) {
	trimmed := strings.Trim(strings.TrimPrefix(path, "/rooms/"), "/")
	if trimmed == "" {
		return "", ""
	}
	parts := strings.Split(trimmed, "/")
	if len(parts) != 2 {
		return "", ""
	}
	return parts[0], parts[1]
}

func handleRoomStoreError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, errRoomNotFound):
		writeError(w, err, http.StatusNotFound)
	case errors.Is(err, errRoomFull):
		writeError(w, err, http.StatusConflict)
	default:
		writeError(w, err, http.StatusInternalServerError)
	}
}

func testBroadcastHandler(hub *wstransport.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeError(w, errors.New("method not allowed"), http.StatusMethodNotAllowed)
			return
		}

		var req testBroadcastRequest
		if err := decodeJSON(r, &req); err != nil {
			writeError(w, err, http.StatusBadRequest)
			return
		}

		roomID := strings.TrimSpace(req.RoomID)
		if roomID == "" {
			writeError(w, errors.New("roomId is required"), http.StatusBadRequest)
			return
		}

		payload, err := json.Marshal(req.Event)
		if err != nil {
			writeError(w, errors.New("invalid event payload"), http.StatusBadRequest)
			return
		}

		hub.BroadcastToRoom(roomID, payload)
		writeJSON(w, http.StatusAccepted, map[string]any{"status": "broadcasted"})
	}
}

func decodeJSON(r *http.Request, v any) error {
	defer r.Body.Close()

	decoder := json.NewDecoder(io.LimitReader(r.Body, 1<<20))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(v); err != nil {
		return err
	}

	if decoder.More() {
		return errors.New("request body must contain a single JSON object")
	}

	return nil
}

func writeError(w http.ResponseWriter, err error, status int) {
	writeJSON(w, status, map[string]any{"error": err.Error()})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
