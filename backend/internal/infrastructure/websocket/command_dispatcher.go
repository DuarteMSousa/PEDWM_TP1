package websocket

// CommandDispatcher is the central component responsible for routing
// WebSocket messages received from clients to the appropriate command handlers.
// Each message type (identified by the "type" field) is associated with a
// previously registered CommandHandler.
//
// This approach allows decoupling the WebSocket transport layer from the
// command processing logic, facilitating system extensibility — to support
// a new command, simply register a new handler.

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
)

// CommandHandler defines the signature of a function capable of processing
// a command received via WebSocket. It receives the command context
// (player identification, room, and reference to the client) and the specific
// JSON payload of the command.
type CommandHandler func(ctx *CommandContext, payload json.RawMessage) error

// CommandContext encapsulates the contextual information associated with a command
// received: the player identifier, the room it belongs to, and the
// reference to the WebSocket client that originated the command.
type CommandContext struct {
	PlayerID string
	RoomID   string
	Client   *Client
}

// CommandDispatcher maintains a registry of handlers indexed by message type
// and routes each received message to the appropriate handler.
type CommandDispatcher struct {
	mu       sync.RWMutex
	handlers map[string]CommandHandler
}

var (
	commandDispatcherInstance *CommandDispatcher
	onceCommandDispatcher     sync.Once
)

// GetCommandDispatcherInstance returns the singleton instance of CommandDispatcher.
func GetCommandDispatcherInstance() *CommandDispatcher {
	onceCommandDispatcher.Do(func() {
		commandDispatcherInstance = &CommandDispatcher{
			handlers: make(map[string]CommandHandler),
		}
	})
	return commandDispatcherInstance
}

// Register associates a message type with a CommandHandler. If a handler
// already exists for the same type, it is replaced.
func (d *CommandDispatcher) Register(messageType string, handler CommandHandler) {
	if d == nil || handler == nil {
		return
	}
	d.handlers[messageType] = handler
}

// Dispatch routes a ClientMessage to the registered handler for the
// respective type. If the command causes a panic (common in existing
// domain commands), it is recovered and returned as an error.
// If no handler exists for the message type, an error is returned.
func (d *CommandDispatcher) Dispatch(ctx *CommandContext, msg ClientMessage) (err error) {
	if d == nil {
		return fmt.Errorf("dispatcher not configured")
	}

	handler, ok := d.handlers[msg.Type]
	if !ok {
		return fmt.Errorf("unknown command: %s", msg.Type)
	}

	// Existing domain commands use panic to signal errors.
	// We recover here to return an error response to the client instead
	// of terminating the goroutine.
	defer func() {
		if r := recover(); r != nil {
			slog.Error("panic recovered in dispatcher", "command", msg.Type, "panic", r)
			if e, ok := r.(error); ok {
				err = e
			} else {
				err = fmt.Errorf("%v", r)
			}
		}
	}()

	return handler(ctx, msg.Payload)
}

// HandleMessage is the main entry point invoked by the client's ReadPump.
// It parses the JSON message, constructs the context, and dispatches
// it to the appropriate handler. The response (success or error) is sent back
// to the client that originated the command.
func (d *CommandDispatcher) HandleMessage(client *Client, raw []byte) {
	if d == nil || client == nil {
		return
	}

	var msg ClientMessage
	if err := json.Unmarshal(raw, &msg); err != nil {
		d.sendError(client, "parse_error", "invalid message format")
		return
	}

	if msg.Type == "" {
		d.sendError(client, "parse_error", "missing message type")
		return
	}

	ctx := &CommandContext{
		PlayerID: client.id,
		RoomID:   client.roomID,
		Client:   client,
	}

	if err := d.Dispatch(ctx, msg); err != nil {
		d.sendError(client, msg.Type, err.Error())
		return
	}

	d.sendSuccess(client, msg.Type)
}

// sendError constructs an error response message and enqueues it to the client.
func (d *CommandDispatcher) sendError(client *Client, msgType string, errMsg string) {
	resp := ServerMessage{
		Type:    msgType,
		Success: false,
		Error:   errMsg,
	}
	data, err := json.Marshal(resp)
	if err != nil {
		return
	}
	client.Enqueue(data)
}

// sendSuccess constructs a success response message and enqueues it to the client.
func (d *CommandDispatcher) sendSuccess(client *Client, msgType string) {
	resp := ServerMessage{
		Type:    msgType,
		Success: true,
	}
	data, err := json.Marshal(resp)
	if err != nil {
		return
	}
	client.Enqueue(data)
}
