package websocket

// Hub representa o componente responsável pela gestão da comunicação
// WebSocket associada às diferentes salas do sistema. Este elemento atua
// como ponto central de distribuição de mensagens em tempo real, permitindo
// organizar a comunicação por contexto de sala (room-scoped fan-out).
//
// A principal responsabilidade deste componente consiste em manter um registo
// das salas ativas de comunicação, assegurando que cada identificador de sala
// está associado a uma instância de RoomHub. Cada RoomHub é responsável por
// gerir os clientes WebSocket ligados a essa sala específica e por efetuar a
// difusão de mensagens para todos os clientes nela registados.
//
// O Hub disponibiliza operações para obter ou criar dinamicamente uma sala de
// comunicação, adicionar clientes, remover clientes e difundir mensagens para
// todos os participantes de uma sala. Esta abordagem permite encapsular a
// gestão da comunicação em tempo real e separar essa responsabilidade da lógica
// de negócio da aplicação.
//
// De forma a garantir segurança em ambientes concorrentes, o acesso ao mapa
// interno de salas é protegido através de um mecanismo de sincronização
// baseado em sync.RWMutex. Esta solução permite múltiplas leituras simultâneas
// e assegura exclusividade durante operações de modificação, como a criação
// ou remoção de salas.
//
// Adicionalmente, quando o último cliente abandona uma sala, a respetiva
// instância de RoomHub pode ser removida do registo interno, contribuindo para
// uma gestão mais eficiente dos recursos do sistema.

import (
	"strings"
	"sync"
)

// Hub is the entrypoint for room-scoped websocket fan-out.
type Hub struct {
	mu    sync.RWMutex
	rooms map[string]*RoomHub
}

func NewHub() *Hub {
	return &Hub{rooms: make(map[string]*RoomHub)}
}

func (h *Hub) GetOrCreateRoom(roomID string) *RoomHub {
	roomID = strings.TrimSpace(roomID)
	if roomID == "" {
		return nil
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	room, ok := h.rooms[roomID]
	if ok {
		return room
	}

	room = NewRoomHub(roomID)
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
