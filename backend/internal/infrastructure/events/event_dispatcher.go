package events_infrastructure

// CommandDispatcher is the central component responsible for routing
// WebSocket messages received from clients to the appropriate command handlers.
// Each message type (identified by the "type" field) is associated with a
// previously registered CommandHandler.
//
// This approach allows decoupling the WebSocket transport layer from the
// processing logic of each command, facilitating system extensibility —
// to support a new command, simply register a new handler.

import (
	"backend/internal/domain/events"
	"fmt"
	"sync"
)

// EventHandler defines the signature of a function capable of processing
// an event received via WebSocket. It receives the context of the event
// (player identification, room, and client reference) and the specific
// JSON payload of the event.
type EventHandler func(event events.Event) error

// EventDispatcher maintains a registry of handlers indexed by event type
// and routes each received event to the appropriate handler.
type EventDispatcher struct {
	mu       sync.RWMutex
	handlers map[string]EventHandler
}

var (
	eventDispatcherInstance *EventDispatcher
	onceEventDispatcher     sync.Once
)

// GetEventDispatcherInstance returns the singleton instance of EventDispatcher.
func GetEventDispatcherInstance() *EventDispatcher {
	onceEventDispatcher.Do(func() {
		eventDispatcherInstance = &EventDispatcher{
			handlers: make(map[string]EventHandler),
		}
	})
	return eventDispatcherInstance
}

// Register associates an event type with an EventHandler. If a handler
// already exists for the same type, it is replaced.
func (d *EventDispatcher) Register(messageType string, handler EventHandler) {
	if d == nil || handler == nil {
		return
	}
	d.handlers[messageType] = handler
}

// Dispatch receives an event, looks up the corresponding handler based on
// the event type, and invokes it. If no handler is found for the event type,
// an error is returned.
func (d *EventDispatcher) Dispatch(event events.Event) error {
	if d == nil {
		return fmt.Errorf("dispatcher not configured")
	}

	handler, ok := d.handlers[string(event.Type)]
	if !ok {
		return fmt.Errorf("unknown event type: %s", event.Type)
	}

	return handler(event)
}

// HandleMessage is the main entry point invoked by the client's ReadPump.
// It parses the JSON message, constructs the context, and dispatches it
// to the appropriate handler. The response (success or error) is sent back
// to the client that originated the event.
func (d *EventDispatcher) HandleMessage(event events.Event) {
	if d == nil {
		return
	}

	d.Dispatch(event)
}
