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
	ds.log.Info("store", "create", "started", "deckID", d.DeckID)
	ds.store[d.DeckID] = &d
	ds.log.Info("store", "create", "finished", "deckID", d.DeckID)
}

func (ds *DeckStore) QueryById(u uuid.UUID) (Deck, error) {
	ds.log.Info("store", "query", "started", "deckID", u)
	d := ds.store[u]
	if d == nil {
		ds.log.Info("store", "query", "not found", "deckID", u)
		return Deck{}, errors.New("Deck not found")
	}
	ds.log.Info("store", "query", "finished", "deckID", u)
	return *d, nil
}
