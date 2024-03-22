package deck

import (
	"errors"
	"log/slog"

	"github.com/google/uuid"
)

func NewStore(log *slog.Logger) *DeckStore {
	ds := DeckStore{
		log:   log,
		store: make(store, 1),
	}
	return &ds
}

type store map[uuid.UUID]*Deck

type DeckStore struct {
	log   *slog.Logger
	store store
}

func (ds *DeckStore) Create(d Deck) {
	ds.store[d.DeckID] = &d
}

func (ds *DeckStore) QueryById(u uuid.UUID) (Deck, error) {
	d := ds.store[u]
	if d == nil {
		return Deck{}, errors.New("Deck not found")
	}
	return *d, nil
}
