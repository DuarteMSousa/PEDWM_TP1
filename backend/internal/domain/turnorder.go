package domain

import "errors"

var (
	ErrTurnOrderInvalidSize  = errors.New("turn order must have exactly 4 players")
	ErrTurnOrderDuplicateIDs = errors.New("turn order has duplicate player ids")
	ErrTurnOrderPlayerAbsent = errors.New("player id not present in turn order")
)

// TurnOrder garante ordem determinística de turnos (crítico em Go, não usar map iteration).
type TurnOrder struct {
	players []string // tamanho 4, ordem fixa
}

func NewTurnOrder(players []string) (TurnOrder, error) {
	if len(players) != 4 {
		return TurnOrder{}, ErrTurnOrderInvalidSize
	}
	seen := map[string]bool{}
	for _, id := range players {
		if id == "" {
			return TurnOrder{}, errors.New("turn order contains empty player id")
		}
		if seen[id] {
			return TurnOrder{}, ErrTurnOrderDuplicateIDs
		}
		seen[id] = true
	}

	cp := make([]string, 4)
	copy(cp, players)

	return TurnOrder{players: cp}, nil
}

func (o TurnOrder) Players() []string {
	cp := make([]string, len(o.players))
	copy(cp, o.players)
	return cp
}

func (o TurnOrder) Contains(playerID string) bool {
	for _, id := range o.players {
		if id == playerID {
			return true
		}
	}
	return false
}

func (o TurnOrder) Next(currentPlayerID string) (string, error) {
	for i, id := range o.players {
		if id == currentPlayerID {
			return o.players[(i+1)%4], nil
		}
	}
	return "", ErrTurnOrderPlayerAbsent
}
