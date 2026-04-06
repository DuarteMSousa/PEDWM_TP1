// Package events_infrastructure implements the observers and handlers for
// events in the infrastructure layer. It bridges the domain events published
// by the EventBus and the infrastructure actions such as persistence,
// statistics update, and room cleanup.
//
// Main components:
//   - EventDispatcher: dispatches events to handlers registered by type.
//   - EventPersistanceObserver: observer that persists events and forwards
//     them to the EventDispatcher.
//   - Event handlers: PlayerLeftEventHandler, GameEndedEventHandler,
//     RoomClosedEventHandler.
package events_infrastructure
