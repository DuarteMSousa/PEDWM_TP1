// Package bot_strategy defines the IBotStrategy interface and concrete implementations
// of strategies for bot players. It uses the Strategy pattern to allow changing
// the bot's behavior at runtime.
//
// Available strategies:
//   - EasyBotStrategy: plays the first card of the leading suit, or the first card in hand.
//   - HardBotStrategy: plays the strongest card of the leading suit, or the strongest card in hand.
package bot_strategy
