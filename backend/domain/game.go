package domain

import (
	"errors"
	"fmt"
)

type GameStatus string

const (
    Aguardando GameStatus = "AGUARDANDO"
    EmJogo    GameStatus = "EM_JOGO"
    FimDeJogo GameStatus = "FIM_DE_JOGO"
)

var (
	ErrGameNotPlaying      = errors.New("game not in playing state")
	ErrNotYourTurn         = errors.New("not your turn")
	ErrPlayerNotFound      = errors.New("player not found")
	ErrTeamNotFound        = errors.New("team not found")
	ErrInvalidPlayerOrder  = errors.New("invalid player order")
	ErrStrategyNotSet      = errors.New("strategy not set")
	ErrEventBusNotSet      = errors.New("event bus not set")
)

type Play struct {
	PlayerID string
	Card     Card
}

type Game struct {
	ID string

	Players map[string]*Player
	Teams   map[string]*Team

	// PlayerOrder define a ordem fixa e determinística dos turnos (tamanho esperado: 4).
	// IMPORTANTÍSSIMO: não uses iteração sobre map para isto.
	PlayerOrder []string

	TurnPlayer   string
	TrumpSuit    Naipe
	CurrentTrick []Play
	TricksPlayed int
	Status       GameStatus

	RuleStrategy    TrickRuleStrategy
	ScoringStrategy ScoringStrategy
	BotStrategy     BotPlayStrategy

	EventBus *EventBus
}

// PlayCard executa uma jogada: valida estado, valida turno, remove carta da mão e adiciona à vaza.
func (g *Game) PlayCard(playerID, cardID string) error {
	if g == nil {
		return errors.New("game is nil")
	}
	if g.Status != Playing {
		return ErrGameNotPlaying
	}
	if g.TurnPlayer != playerID {
		return ErrNotYourTurn
	}

	player, ok := g.Players[playerID]
	if !ok || player == nil {
		return ErrPlayerNotFound
	}

	card, found := player.RemoveCard(cardID)
	if !found {
		return ErrCardNotFound
	}

	g.CurrentTrick = append(g.CurrentTrick, Play{
		PlayerID: playerID,
		Card:     card,
	})

	// EventBus é opcional, mas se queres obrigar, troca para: if g.EventBus == nil { return ErrEventBusNotSet }
	if g.EventBus != nil {
		g.EventBus.Publish(NewCardPlayedEvent(g.ID, playerID, card))
	}

	// Se 4 cartas jogadas → fechar vaza
	if len(g.CurrentTrick) == 4 {
		g.endTrick()
		return nil
	}

	g.advanceTurn()
	return nil
}

func (g *Game) endTrick() {
	// Estratégias são essenciais. Aqui opto por falhar “silenciosamente” com panic evitado,
	// mas podes trocar para retornar error se preferires.
	if g.RuleStrategy == nil || g.ScoringStrategy == nil {
		// Se queres comportamento estrito:
		// panic(ErrStrategyNotSet)
		return
	}

	winnerID := g.RuleStrategy.Winner(g.TrumpSuit, g.CurrentTrick)
	points := g.ScoringStrategy.TrickPoints(g.CurrentTrick)

	winner, ok := g.Players[winnerID]
	if !ok || winner == nil {
		// Estado inconsistente
		return
	}

	team, ok := g.Teams[winner.TeamID]
	if !ok || team == nil {
		// Estado inconsistente
		return
	}

	_ = team.AddPoints(points)

	if g.EventBus != nil {
		g.EventBus.Publish(NewTrickEndedEvent(g.ID, winnerID, points))
	}

	g.TricksPlayed++
	g.CurrentTrick = g.CurrentTrick[:0]
	g.TurnPlayer = winnerID

	// 10 vazas é típico para baralho de 40 cartas (4 jogadores -> 10 trick).
	if g.TricksPlayed == 10 {
		g.Status = Ended
		if g.EventBus != nil {
			g.EventBus.Publish(NewGameEndedEvent(g.ID))
		}
	}
}

func (g *Game) advanceTurn() {
	order := g.PlayerOrder
	if len(order) != 4 {
		// Se preferires: panic / error. Aqui só não avança.
		return
	}

	for i, id := range order {
		if id == g.TurnPlayer {
			g.TurnPlayer = order[(i+1)%4]
			return
		}
	}

	// Se o TurnPlayer atual não existe na ordem, tenta recuperar para o primeiro.
	g.TurnPlayer = order[0]
}

// Validate valida consistência estrutural mínima do jogo (útil ao iniciar).
func (g *Game) Validate() error {
	if g == nil {
		return errors.New("game is nil")
	}
	if len(g.Players) != 4 {
		return fmt.Errorf("expected 4 players, got %d", len(g.Players))
	}
	if len(g.Teams) == 0 {
		return errors.New("no teams configured")
	}
	if len(g.PlayerOrder) != 4 {
		return ErrInvalidPlayerOrder
	}
	if !g.TrumpSuit.Valid() {
		return fmt.Errorf("%w: %q", ErrInvalidNaipe, g.TrumpSuit)
	}
	if g.RuleStrategy == nil || g.ScoringStrategy == nil {
		return ErrStrategyNotSet
	}
	return nil
}