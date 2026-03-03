package domain

import (
	"errors"
	"fmt"
)

type GameStatus string

const (
	Aguardando GameStatus = "AGUARDANDO"
	EmJogo     GameStatus = "EM_JOGO"
	FimDeJogo  GameStatus = "FIM_DE_JOGO"
)

var (
	ErrGameNotPlaying      = errors.New("game not in playing state")
	ErrNotYourTurn         = errors.New("not your turn")
	ErrPlayerNotFound      = errors.New("player not found")
	ErrTeamNotFound        = errors.New("team not found")
	ErrInvalidPlayerOrder  = errors.New("invalid player order")
	ErrStrategyNotSet      = errors.New("strategy not set")
	ErrEventBusNotSet      = errors.New("event bus not set")
	ErrRoundNotConfigured  = errors.New("round not configured")
	ErrTrickNotConfigured  = errors.New("current trick not configured")
)

// Play mantém-se como antes
type Play struct {
	PlayerID string
	Card     Card
}

type Game struct {
	ID string

	Players map[string]*Player
	Teams   map[string]*Team

	// Ordem determinística (em vez de []string solta)
	Order TurnOrder

	// Estado corrente (turno atual)
	TurnPlayer string

	// Mão (10 vazas): trunfo + vaza atual + contador
	Round *Round

	Status GameStatus

	RuleStrategy    TrickRuleStrategy
	ScoringStrategy ScoringStrategy
	BotStrategy     BotPlayStrategy

	EventBus *EventBus
}

// Start coloca o jogo em EM_JOGO (assumindo que Round já foi criado).
func (g *Game) Start() error {
	if g.Round == nil {
		return ErrRoundNotConfigured
	}
	g.Status = EmJogo
	return nil
}

// PlayCard executa uma jogada: valida estado, valida turno, remove carta da mão e adiciona à vaza.
func (g *Game) PlayCard(playerID, cardID string) error {
	if g == nil {
		return errors.New("game is nil")
	}
	if g.Status != EmJogo {
		return ErrGameNotPlaying
	}
	if g.Round == nil {
		return ErrRoundNotConfigured
	}
	if g.Round.CurrentTrick == nil {
		return ErrTrickNotConfigured
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

	// Se a tua interface TrickRuleStrategy tiver ValidatePlay, valida aqui.
	// (Se não tiveres isso implementado ainda, podes comentar este bloco.)
	if g.RuleStrategy != nil && g.Round.CurrentTrick.LeadSuit != nil {
		// leadSuit := *g.Round.CurrentTrick.LeadSuit
		// if err := g.RuleStrategy.ValidatePlay(player.Hand, leadSuit, card); err != nil {
		//     // repõe carta (opcional) ou trata erro como preferires
		//     return err
		// }
	}

	play := Play{PlayerID: playerID, Card: card}
	if err := g.Round.CurrentTrick.AddPlay(play); err != nil {
		return err
	}

	if g.EventBus != nil {
		g.EventBus.Publish(NewCardPlayedEvent(g.ID, playerID, card))
	}

	// Se 4 cartas jogadas → fechar vaza
	if g.Round.CurrentTrick.IsComplete() {
		g.endTrick()
		return nil
	}

	// avançar turno usando TurnOrder
	next, err := g.Order.Next(g.TurnPlayer)
	if err != nil {
		return err
	}
	g.TurnPlayer = next

	return nil
}

func (g *Game) endTrick() {
	if g.Round == nil || g.Round.CurrentTrick == nil {
		return
	}
	if g.RuleStrategy == nil || g.ScoringStrategy == nil {
		return
	}

	plays := g.Round.CurrentTrick.Plays
	winnerID := g.RuleStrategy.Winner(g.Round.TrumpSuit, plays)
	points := g.ScoringStrategy.TrickPoints(plays)

	winner, ok := g.Players[winnerID]
	if !ok || winner == nil {
		return
	}

	team, ok := g.Teams[winner.TeamID]
	if !ok || team == nil {
		return
	}

	_ = team.AddPoints(points)

	if g.EventBus != nil {
		g.EventBus.Publish(NewTrickEndedEvent(g.ID, winnerID, points))
	}

	// incrementa contador de vazas
	g.Round.IncrementTrick()

	// nova vaza começa com o vencedor
	g.Round.StartNewTrick(winnerID)
	g.TurnPlayer = winnerID

	// fim após 10 vazas
	if g.Round.IsFinished() {
		g.Status = FimDeJogo
		if g.EventBus != nil {
			g.EventBus.Publish(NewGameEndedEvent(g.ID))
		}
	}
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
	if !g.Order.Contains(g.TurnPlayer) {
		return ErrInvalidPlayerOrder
	}
	if g.Round == nil {
		return ErrRoundNotConfigured
	}
	if !g.Round.TrumpSuit.Valid() {
		return fmt.Errorf("%w: %q", ErrInvalidNaipe, g.Round.TrumpSuit)
	}
	if g.RuleStrategy == nil || g.ScoringStrategy == nil {
		return ErrStrategyNotSet
	}
	return nil
}