// Package game defines the Game aggregate and its state machine
// (Starting → Playing → Finished). It uses the State pattern to manage
// transitions and the Strategy pattern for score calculation.
//
// Main components:
//   - Game: root aggregate with players, teams, score, and rounds.
//   - IGameState: interface for the State pattern (GameStartingState,
//     GamePlayingState, GameFinishedState).
//   - IGameScoringStrategy: interface for the Strategy pattern for game-level
//     scoring (e.g., SuecaGameScoringStrategy).
package game
