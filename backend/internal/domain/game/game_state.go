package game

// IGameState defines the interface for the State pattern for the game's lifecycle.
type IGameState interface {
	Enter()
	Update()
}
