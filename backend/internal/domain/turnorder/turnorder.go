package turnorder

import (
	"backend/internal/domain/player"
	"errors"
)

var (
	ErrTurnOrderInvalidSize  = errors.New("turn order must have exactly 4 players")
	ErrTurnOrderDuplicateIDs = errors.New("turn order has duplicate player ids")
	ErrTurnOrderPlayerAbsent = errors.New("player id not present in turn order")
)

// TurnOrder garante ordem determinística de turnos (crítico em Go, não usar map iteration).

// TurnOrder implementa uma fila circular (queue) para ordem de turnos.
type TurnOrder struct {
	players []*player.Player // fila de jogadores
}

// NewTurnOrder cria uma fila de turnos começando pelo líder e seguindo a ordem original, circular.
func NewTurnOrder(leaderId string, players []*player.Player) (TurnOrder, error) {
	if len(players) != 4 {
		return TurnOrder{}, ErrTurnOrderInvalidSize
	}
	seen := map[string]bool{}
	leaderIdx := -1
	for i, p := range players {
		if p.ID == "" {
			return TurnOrder{}, errors.New("turn order contains empty player id")
		}
		if seen[p.ID] {
			return TurnOrder{}, ErrTurnOrderDuplicateIDs
		}
		seen[p.ID] = true
		if p.ID == leaderId {
			leaderIdx = i
		}
	}
	if leaderIdx == -1 {
		return TurnOrder{}, errors.New("leader id not found in players")
	}
	// Rearranjar para começar pelo líder
	cp := make([]*player.Player, len(players))
	for i := 0; i < len(players); i++ {
		cp[i] = players[(leaderIdx+i)%len(players)]
	}
	return TurnOrder{players: cp}, nil
}

// Contains verifica se o jogador está na fila.
func (o TurnOrder) Contains(playerID string) bool {
	for _, p := range o.players {
		if p.ID == playerID {
			return true
		}
	}
	return false
}

// Next retorna o próximo jogador
func (o TurnOrder) Next() (string, error) {
	if len(o.players) == 0 {
		return "", errors.New("turn order is empty")
	}
	return o.players[0].ID, nil
}

func (o *TurnOrder) Enqueue(player *player.Player) {
	o.players = append(o.players, player)
}

func (o *TurnOrder) Dequeue() (string, error) {
	if len(o.players) == 0 {
		return "", errors.New("turn order is empty")
	}
	first := o.players[0]
	o.players = o.players[1:]
	return first.ID, nil
}
