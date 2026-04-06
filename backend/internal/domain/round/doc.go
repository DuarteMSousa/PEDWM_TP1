// Package round defines the Round entity (a round of a Sueca game).
// Each round has a trump, a deck, tricks, and a state machine (Setup → Playing → Finished).
// The round orchestrates the distribution of cards, the playing of cards, and the calculation of scores per trick.
//
// It uses the State pattern (IRoundState) and the Strategy pattern
// (IRoundRuleStrategy) for specific rules of the Sueca game.
package round
