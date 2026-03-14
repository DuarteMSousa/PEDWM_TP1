package domain

import (
	"errors"
	"fmt"
	"strings"
)

// Naipe representa o naipe.
// Invariante: Naipe deve pertencer ao conjunto permitido.
type Naipe string

// Rank representa o valor/figura.
// Invariante: Rank deve pertencer ao conjunto permitido.
type Rank string

const (
	Copas   Naipe = "COPAS"
	Espadas Naipe = "ESPADAS"
	Ouros   Naipe = "OUROS"
	Paus    Naipe = "PAUS"
)

const (
	A     Rank = "A"
	K     Rank = "K"
	Q     Rank = "Q"
	J     Rank = "J"
	Seven Rank = "7"
	Six   Rank = "6"
	Five  Rank = "5"
	Four  Rank = "4"
	Three Rank = "3"
	Two   Rank = "2"
)

var (
	ErrInvalidNaipe  = errors.New("naipe inválido")
	ErrInvalidRank   = errors.New("rank inválido")
	ErrInvalidCardID = errors.New("id inválido")
)

// Valid indica se o naipe pertence ao conjunto permitido.
func (n Naipe) Valid() bool {
	switch n {
	case Copas, Espadas, Ouros, Paus:
		return true
	default:
		return false
	}
}

// Valid indica se o rank pertence ao conjunto permitido.
func (r Rank) Valid() bool {
	switch r {
	case A, K, Q, J, Seven, Six, Five, Four, Three, Two:
		return true
	default:
		return false
	}
}

// Card é uma entidade do domínio que representa uma carta.
//   - ID não deve ser vazio (se for relevante no teu contexto)
//   - Naipe deve ser válido
//   - Rank deve ser válido
type Card struct {
	ID    string
	Naipe Naipe
	Rank  Rank
}

// NewCard é o construtor canónico: garante invariantes.
func NewCard(id string, naipe Naipe, rank Rank) (Card, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return Card{}, ErrInvalidCardID
	}
	if !naipe.Valid() {
		return Card{}, fmt.Errorf("%w: %q", ErrInvalidNaipe, naipe)
	}
	if !rank.Valid() {
		return Card{}, fmt.Errorf("%w: %q", ErrInvalidRank, rank)
	}
	return Card{ID: id, Naipe: naipe, Rank: rank}, nil
}

// Validate valida uma carta já existente (útil para dados vindos de fora).
func (c Card) Validate() error {
	if strings.TrimSpace(c.ID) == "" {
		return ErrInvalidCardID
	}
	if !c.Naipe.Valid() {
		return fmt.Errorf("%w: %q", ErrInvalidNaipe, c.Naipe)
	}
	if !c.Rank.Valid() {
		return fmt.Errorf("%w: %q", ErrInvalidRank, c.Rank)
	}
	return nil
}

// IsTrump indica se a carta pertence ao naipe de trunfo.
func (c Card) IsTrump(trump Naipe) bool {
	return c.Naipe == trump
}
