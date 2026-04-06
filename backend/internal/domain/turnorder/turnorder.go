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

// TurnOrder implements a deterministic turn order.

// TurnOrder implements a circular queue for turn order.
type TurnOrder struct {
	players        []*player.Player // queue of players
	playersHistory []*player.Player // history of players for re-entry
}

// NewTurnOrder creates a turn order starting from the leader and following the original order.
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
	// Reorder players starting from leader
	cp := make([]*player.Player, len(orderedPlayers))
	for i := 0; i < len(orderedPlayers); i++ {
		cp[i] = orderedPlayers[(leaderIdx+i)%len(orderedPlayers)]
	}
	return TurnOrder{players: cp}, nil
}

// Contains checks if the player is in the queue.
func (o TurnOrder) Contains(playerID string) bool {
	for _, p := range o.players {
		if p.ID == playerID {
			return true
		}
	}
	return false
}

// Next returns the next player
func (o TurnOrder) Next() (string, error) {
	if len(o.players) == 0 {
		return "", errors.New("turn order is empty")
	}
	return o.players[0].ID, nil
}

// Enqueue adds a player to the end of the queue.
func (o *TurnOrder) Enqueue(player *player.Player) {
	o.players = append(o.players, player)
}

// Dequeue removes and returns the first player in the queue.
func (o *TurnOrder) Dequeue() (string, error) {
	if len(o.players) == 0 {
		return "", errors.New("turn order is empty")
	}
	first := o.players[0]
	o.players = o.players[1:]
	o.playersHistory = append(o.playersHistory, first)
	return first.ID, nil
}

// Remove removes a player from the queue by ID.
func (o *TurnOrder) Remove(playerID string) error {
	for i, p := range o.players {
		if p.ID == playerID {
			o.players = append(o.players[:i], o.players[i+1:]...)
			return nil
		}
	}
	return ErrTurnOrderPlayerAbsent
}

// AddPlayer adds a player back to the queue.
// If the player has already played, they are added to the end of the queue.
// Otherwise, they are added in their original position.
func (o *TurnOrder) AddPlayer(player *player.Player) {
	alreadyPlayed := false
	for _, p := range o.playersHistory {
		if p.Sequence == player.Sequence {
			alreadyPlayed = true
			break
		}
	}

	if alreadyPlayed {
		o.playersHistory = append(o.playersHistory, player)
	} else {
		o.players = append(o.players, player)
	}

	sort.Slice(o.players, func(i, j int) bool {
		if o.players[i].Sequence == o.players[j].Sequence {
			return o.players[i].ID < o.players[j].ID
		}
		return o.players[i].Sequence < o.players[j].Sequence
	})

}
