package websocket

// RoomHub representa o componente responsável pela gestão das ligações
// WebSocket associadas a uma sala específica do sistema. Cada instância de
// RoomHub mantém o conjunto de clientes ligados a uma determinada sala e
// assegura a distribuição de mensagens em tempo real para todos os seus
// participantes.
//
// Este componente funciona como um mecanismo de difusão (broadcast) dentro
// do contexto de uma sala, permitindo que mensagens recebidas pelo servidor
// sejam enviadas a todos os clientes atualmente conectados a essa sala.
//
// Internamente, o RoomHub mantém um registo dos clientes ativos através de
// uma estrutura baseada num mapa, onde cada cliente WebSocket é representado
// por uma instância de Client. As operações de adição e remoção de clientes
// permitem gerir dinamicamente o conjunto de participantes à medida que
// estabelecem ou terminam a sua ligação ao sistema.
//
// Para garantir segurança em ambientes concorrentes, o acesso ao conjunto
// de clientes é protegido através de um mecanismo de sincronização
// (sync.RWMutex). Este mecanismo permite múltiplas leituras simultâneas,
// enquanto assegura exclusividade durante operações de modificação,
// como a adição ou remoção de clientes.
//
// A funcionalidade de broadcast é responsável por enviar uma determinada
// mensagem a todos os clientes registados na sala. Caso um cliente não
// consiga receber a mensagem (por exemplo, devido a lentidão ou bloqueio
// no canal de envio), este pode ser automaticamente desconectado para
// preservar o desempenho e a estabilidade da comunicação em tempo real.

import (
	"backend/internal/domain/room"
	"sync"
)

// RoomHub manages clients for a single room.
type RoomHub struct {
	room    *room.Room
	mu      sync.RWMutex
	clients map[*Client]struct{}
}

func NewRoomHub(room *room.Room) *RoomHub {
	return &RoomHub{
		room:    room,
		clients: make(map[*Client]struct{}),
	}
}

func (r *RoomHub) GetRoom() *room.Room {
	if r == nil {
		return nil
	}
	return r.room
}

func (r *RoomHub) AddClient(client *Client) {
	if r == nil || client == nil {
		return
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.clients[client] = struct{}{}
}

// RemoveClient removes a client and reports whether the room is now empty.
func (r *RoomHub) RemoveClient(client *Client) bool {
	if r == nil || client == nil {
		return true
	}

	//Por aqui a lancar evento de player left

	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.clients, client)
	return len(r.clients) == 0
}

func (r *RoomHub) IsEmpty() bool {
	if r == nil {
		return true
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.clients) == 0
}

func (r *RoomHub) Broadcast(payload []byte) {
	if r == nil || len(payload) == 0 {
		return
	}

	r.mu.RLock()
	clients := make([]*Client, 0, len(r.clients))
	for client := range r.clients {
		clients = append(clients, client)
	}
	r.mu.RUnlock()

	for _, client := range clients {
		if ok := client.Enqueue(payload); ok {
			continue
		}
		// Slow clients are disconnected to keep room broadcast healthy.
		client.Close()
	}
}
