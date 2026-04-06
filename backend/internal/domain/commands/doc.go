// Package command implements the Command pattern to encapsulate player actions
// as executable objects. Each command receives the necessary context during construction
// and executes on the appropriate aggregate.
//
// Available commands:
//   - PlayCardCommand: plays a card in an active game.
//   - ChangeBotStrategyCommand: changes the bot strategy in a room.
package command
