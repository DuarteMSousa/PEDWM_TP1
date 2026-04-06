// Package trick defines the Trick entity and its components.
// A trick represents a sequence of plays (up to 4) in a round,
// controlling the turn order, the leading suit, and the determination of the
// winner.
//
// It uses the Strategy pattern for rules (ITrickRuleStrategy) and scoring
// (ITrickScoringStrategy), allowing for different variants of card games.
// The current implementation supports Sueca.
package trick
