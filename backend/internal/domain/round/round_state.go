package round

// IRoundState defines the interface for the State pattern for the lifecycle of a round.
type IRoundState interface {
	Enter()
	Update()
}
