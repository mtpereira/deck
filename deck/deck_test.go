package deck

import (
	"io"
	"log/slog"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	// Test sorted deck creation.
	da := NewAPI(slog.New(slog.NewTextHandler(io.Discard, nil)))

	sortedCards := []Card{
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

	sortedDeck := Deck{
		Shuffled:  false,
		Remaining: 52,
		Cards:     sortedCards,
	}

	d := da.New(false, nil)
	if d.Shuffled != sortedDeck.Shuffled {
		t.Fatalf("Expected deck.Shuffled to be false, it is not")
	}
	if !reflect.DeepEqual(d.Cards, sortedCards) {
		t.Fatalf("Deck cards are not sorted, got %v, want %v", d.Cards, sortedCards)
	}

	// Test shuffled deck creation.
	d = da.New(true, nil)
	if d.Shuffled == sortedDeck.Shuffled {
		t.Fatalf("Expected deck.Shuffled to be true, it is not")
	}
	if reflect.DeepEqual(d.Cards, sortedCards) {
		t.Fatalf("Deck cards are not shuffled, got %v", d.Cards)
	}

}
