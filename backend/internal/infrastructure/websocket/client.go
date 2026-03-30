package websocket

// Client representa uma ligação WebSocket individual estabelecida entre um
// cliente remoto e o servidor. Este componente é responsável por encapsular
// a conexão ativa, gerir o envio assíncrono de mensagens e controlar o ciclo
// de vida da comunicação associada a um utilizador dentro de uma determinada sala.
//
// Cada instância de Client mantém a identificação do cliente, o identificador
// da sala a que pertence, uma referência ao Hub responsável pela gestão global
// das salas e a conexão WebSocket subjacente. Adicionalmente, dispõe de um
// canal interno de envio (send), utilizado para enfileirar mensagens a transmitir
// de forma assíncrona para o cliente.
//
// A arquitetura adotada separa explicitamente as operações de leitura e de
// escrita através dos métodos ReadPump e WritePump. O método ReadPump é
// responsável por consumir mensagens recebidas da ligação WebSocket e por
// monitorizar a continuidade da ligação através da receção de frames Pong.
// Por sua vez, o método WritePump trata do envio de mensagens pendentes e da
// transmissão periódica de frames Ping, permitindo verificar se a ligação se
// mantém ativa.
//
// Para reforçar a robustez do sistema, o componente define limites de leitura,
// tempos máximos de espera e mecanismos de heartbeat baseados em Ping/Pong.
// Esta abordagem reduz o risco de ligações bloqueadas ou inativas permanecerem
// abertas indefinidamente.
//
// O encerramento da ligação é controlado de forma segura através de sync.Once,
// garantindo que a operação de fecho é executada apenas uma única vez, mesmo
// em cenários concorrentes. Durante esse processo, o cliente é também removido
// da sala correspondente no Hub, assegurando a consistência do estado interno
// do sistema.
//
// Em conjunto, este componente constitui a unidade fundamental da comunicação
// em tempo real, servindo de abstração para cada participante ligado ao sistema
// através de WebSocket.

import (
	"sync"
	"time"

	gws "github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 8 * 1024
)

// Client wraps one websocket connection.
type Client struct {
	id         string
	roomID     string
	hub        *Hub
	conn       *gws.Conn
	send       chan []byte
	once       sync.Once
	dispatcher *CommandDispatcher
}

func NewClient(id string, roomID string, conn *gws.Conn, hub *Hub, dispatcher *CommandDispatcher) *Client {
	return &Client{
		id:         id,
		roomID:     roomID,
		hub:        hub,
		conn:       conn,
		send:       make(chan []byte, 32),
		dispatcher: dispatcher,
	}
}

func (c *Client) Enqueue(payload []byte) bool {
	if c == nil || len(payload) == 0 {
		return false
	}

	select {
	case c.send <- payload:
		return true
	default:
		return false
	}
}

func (c *Client) Close() {
	if c == nil {
		return
	}

	c.once.Do(func() {
		if c.hub != nil {
			c.hub.RemoveClient(c.roomID, c)
		}
		_ = c.conn.Close()
	})
}

func (c *Client) ReadPump() {
	if c == nil {
		return
	}
	defer c.Close()

	c.conn.SetReadLimit(maxMessageSize)
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		return c.conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			return
		}

		if c.dispatcher != nil {
			c.dispatcher.HandleMessage(c, message)
		}
	}
}

func (c *Client) WritePump() {
	if c == nil {
		return
	}
	defer c.Close()

	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for {
		select {
		case message := <-c.send:
			if err := c.write(gws.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(gws.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) write(messageType int, payload []byte) error {
	_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
	return c.conn.WriteMessage(messageType, payload)
}
