package websocket

// Client represents an individual WebSocket connection established between a
// remote client and the server. This component is responsible for encapsulating
// the active connection, managing asynchronous message sending, and controlling
// the lifecycle of the communication associated with a user within a specific room.
//
// Each Client instance maintains the client ID, the room ID it belongs to, a reference
// to the Hub responsible for global room management, and the underlying WebSocket
// connection. Additionally, it has an internal send channel used to queue messages
// for asynchronous transmission to the client.
//
// The adopted architecture explicitly separates read and write operations through
// the ReadPump and WritePump methods. The ReadPump method is responsible for
// consuming messages received from the WebSocket connection and monitoring the
// connection's continuity through the reception of Pong frames. Conversely, the
// WritePump method handles sending pending messages and periodically transmitting
// Ping frames, allowing verification of the connection's active status.
//
// To reinforce system robustness, the component defines read limits,
// maximum wait times, and heartbeat mechanisms based on Ping/Pong.
// This approach reduces the risk of blocked or inactive connections remaining
// open indefinitely.
//
// The closure of the connection is safely managed using sync.Once,
// ensuring that the close operation is executed only once, even
// in concurrent scenarios. During this process, the client is also removed
// from the corresponding room in the Hub, ensuring the consistency of the internal
// system state.

import (
	websocket_interfaces "backend/internal/infrastructure/websocket/interfaces"
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
	id          string
	roomID      string
	hub         *Hub
	conn        *gws.Conn
	send        chan []byte
	once        sync.Once
	dispatcher  *CommandDispatcher
	roomService websocket_interfaces.RoomService
}

// NewClient creates a new WebSocket client.
func NewClient(id string, roomID string, conn *gws.Conn, hub *Hub, dispatcher *CommandDispatcher, roomService websocket_interfaces.RoomService) *Client {
	return &Client{
		id:          id,
		roomID:      roomID,
		hub:         hub,
		conn:        conn,
		send:        make(chan []byte, 32),
		dispatcher:  dispatcher,
		roomService: roomService,
	}
}

// Enqueue adds a message to the client's send queue.
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

// Close safely closes the WebSocket connection (sync.Once).
func (c *Client) Close() {
	if c == nil {
		return
	}

	c.once.Do(func() {
		if c.roomService != nil {
			_, _ = c.roomService.LeaveRoom(c.roomID, c.id)
		}
		if c.hub != nil {
			c.hub.RemoveClient(c.roomID, c)
		}
		if c.conn != nil {
			_ = c.conn.Close()
		}
	})
}

// ReadPump reads messages from the WebSocket and forwards them to the dispatcher.
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

// WritePump sends pending messages and periodic ping/heartbeat.
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

// write sends a message to the WebSocket connection with a specified message type.
func (c *Client) write(messageType int, payload []byte) error {
	_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
	return c.conn.WriteMessage(messageType, payload)
}
