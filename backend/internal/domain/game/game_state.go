package game

type IGameState interface {
	Enter()
	Update()
}
