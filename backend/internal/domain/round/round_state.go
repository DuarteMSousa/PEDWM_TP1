package round

type IRoundState interface {
	Enter()
	Update()
}
