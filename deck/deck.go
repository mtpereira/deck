package deck

import (
	"math/rand"

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

func New(shuffle bool) Deck {
	cards := getSortedCards()

	if shuffle {
		rand.Shuffle(len(cards), func(i, j int) {
			cards[i], cards[j] = cards[j], cards[i]
		})
	}

	return Deck{
		DeckID:    uuid.New(),
		Shuffled:  shuffle,
		Remaining: 52,
		Cards:     cards,
	}
}
