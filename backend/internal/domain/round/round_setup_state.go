package round

import (
	"backend/internal/domain/deck/deckFactory"
	"backend/internal/domain/events"
)

const NUMBER_OF_SUECA_CARDS_PER_PLAYER = 10

// RoundSetupState implements GameState
type RoundSetupState struct {
	round *Round
}

// NewRoundSetupState creates a new instance of RoundSetupState
func NewRoundSetupState(r *Round) *RoundSetupState {
	return &RoundSetupState{round: r}
}

// Enter initializes the round by shuffling the deck, revealing the trump suit, and dealing cards to players
func (s *RoundSetupState) Enter() {
	s.round.Deck = deckFactory.CreateSuecaDeck()
	s.round.Deck.Shuffle()
	firstCard, firstCardErr := s.round.Deck.First()

	if firstCardErr != nil {
		panic("Failed to get the first card from the deck: " + firstCardErr.Error())
	}

	s.round.TrumpSuit = firstCard.Suit
	s.round.AddEvent(events.NewTrumpRevealedEvent(s.round.gameId.String(), firstCard))

	for _, team := range s.round.Teams {
		for _, player := range team.Players {
			for i := 0; i < NUMBER_OF_SUECA_CARDS_PER_PLAYER; i++ {
				card, err := s.round.Deck.Draw()
				if err != nil {
					panic("Failed to draw a card from the deck: " + err.Error())
				}
				player.Hand.AddCard(card)
				s.round.AddEvent(events.NewCardDealtEvent(s.round.gameId.String(), player.ID, card))
			}
		}
	}

	s.round.State.Update()
}

// Update transitions the round to the playing state after setup is complete
func (s *RoundSetupState) Update() {
	s.round.AddEvent(events.NewRoundStartedEvent(s.round.gameId.String()))
	s.round.State = NewRoundPlayingState(s.round)
	s.round.State.Enter()
}
