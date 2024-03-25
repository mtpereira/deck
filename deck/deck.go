package deck

import (
	"errors"
	"log/slog"
	"math/rand"
	"slices"
	"sync"

	"github.com/google/uuid"
)

type Deck struct {
	DeckID    uuid.UUID `json:"deck_id"`
	Shuffled  bool      `json:"shuffled"`
	Remaining int       `json:"remaining"`
	Cards     []Card    `json:"cards,omitempty"`
}

type Card struct {
	Code  string `json:"code"`
	Value string `json:"value,omitempty"`
	Suit  string `json:"suit,omitempty"`
}

type DeckAPI struct {
	store *DeckStore
	mu    sync.Mutex
	log   *slog.Logger
}

var ErrUnsufficientCards error = errors.New("Deck doesn't have that many cards to draw")
var ErrDeckNotFound error = errors.New("Deck not found")

func getSortedCards() []Card {
	return []Card{
		// Clubs
		{
			Code: "2C",
		},
		{
			Code: "3C",
		},
		{
			Code: "4C",
		},
		{
			Code: "5C",
		},
		{
			Code: "6C",
		},
		{
			Code: "7C",
		},
		{
			Code: "8C",
		},
		{
			Code: "9C",
		},
		{
			Code: "10C",
		},
		{
			Code: "JC",
		},
		{
			Code: "QC",
		},
		{
			Code: "KC",
		},
		{
			Code: "AC",
		},

		// Diamonds
		{
			Code: "2D",
		},
		{
			Code: "3D",
		},
		{
			Code: "4D",
		},
		{
			Code: "5D",
		},
		{
			Code: "6D",
		},
		{
			Code: "7D",
		},
		{
			Code: "8D",
		},
		{
			Code: "9D",
		},
		{
			Code: "10D",
		},
		{
			Code: "JD",
		},
		{
			Code: "QD",
		},
		{
			Code: "KD",
		},
		{
			Code: "AD",
		},

		// Hearts
		{
			Code: "2H",
		},
		{
			Code: "3H",
		},
		{
			Code: "4H",
		},
		{
			Code: "5H",
		},
		{
			Code: "6H",
		},
		{
			Code: "7H",
		},
		{
			Code: "8H",
		},
		{
			Code: "9H",
		},
		{
			Code: "10H",
		},
		{
			Code: "JH",
		},
		{
			Code: "QH",
		},
		{
			Code: "KH",
		},
		{
			Code: "AH",
		},

		// Spades
		{
			Code: "2H",
		},
		{
			Code: "3H",
		},
		{
			Code: "4H",
		},
		{
			Code: "5H",
		},
		{
			Code: "6H",
		},
		{
			Code: "7H",
		},
		{
			Code: "8H",
		},
		{
			Code: "9H",
		},
		{
			Code: "10H",
		},
		{
			Code: "JH",
		},
		{
			Code: "QH",
		},
		{
			Code: "KH",
		},
		{
			Code: "AH",
		},
	}
}

func NewAPI(log *slog.Logger) *DeckAPI {
	return &DeckAPI{
		log:   log,
		store: NewStore(log),
	}
}

func (da *DeckAPI) New(shuffle bool, cards []Card) *Deck {
	if cards == nil {
		cards = getSortedCards()
	}

	if shuffle {
		rand.Shuffle(len(cards), func(i, j int) {
			cards[i], cards[j] = cards[j], cards[i]
		})
	}

	d := &Deck{
		DeckID:    uuid.New(),
		Shuffled:  shuffle,
		Remaining: len(cards),
		Cards:     cards,
	}

	da.store.Create(*d)

	return d
}

func (da *DeckAPI) Get(u uuid.UUID) (*Deck, error) {
	d, err := da.store.QueryById(u)
	if err != nil {
		return nil, ErrDeckNotFound
	}
	return &d, nil
}

func (da *DeckAPI) Draw(u uuid.UUID, n int) ([]Card, error) {
	da.mu.Lock()
	defer da.mu.Unlock()

	d, err := da.store.QueryById(u)
	if err != nil {
		return nil, ErrDeckNotFound
	}

	if n > d.Remaining {
		return nil, ErrUnsufficientCards
	}

	var drawn []Card
	drawn = append(drawn, d.Cards[0:n]...)
	d.Cards = slices.Delete(d.Cards, 0, n)
	d.Remaining = len(d.Cards)
	updatedDeck := Deck{DeckID: d.DeckID, Cards: d.Cards, Remaining: len(d.Cards)}
	da.store.Update(u, updatedDeck)
	return drawn, nil
}
