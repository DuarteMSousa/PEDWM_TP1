package turnorder

import (
	"backend/internal/domain/player"
	"errors"
	"sort"
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
	orderedPlayers := make([]*player.Player, len(players))
	copy(orderedPlayers, players)
	sort.Slice(orderedPlayers, func(i, j int) bool {
		if orderedPlayers[i].Sequence == orderedPlayers[j].Sequence {
			return orderedPlayers[i].ID < orderedPlayers[j].ID
		}
		return orderedPlayers[i].Sequence < orderedPlayers[j].Sequence
	})

	seen := map[string]bool{}
	leaderIdx := -1
	for i, p := range orderedPlayers {
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
	cp := make([]*player.Player, len(orderedPlayers))
	for i := 0; i < len(orderedPlayers); i++ {
		cp[i] = orderedPlayers[(leaderIdx+i)%len(orderedPlayers)]
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

func (o *TurnOrder) Advance() (string, error) {
	if o == nil || len(o.players) == 0 {
		return "", errors.New("turn order is empty")
	}

	first := o.players[0]
	o.players = append(o.players[1:], first)
	return first.ID, nil
}
