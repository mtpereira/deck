package web

import (
	"log/slog"
	"net/http"

	"github.com/mtpereira/deck/deck"
)

func addRoutes(mux *http.ServeMux, log *slog.Logger, ds *deck.DeckStore) {
	logRequests := newLoggerMiddleware(log)
	mux.Handle("POST /v1/decks", logRequests(handlePostDeck(ds)))
	mux.Handle("GET /v1/decks/{deck_id}", logRequests(handleGetDeck(ds)))
	mux.Handle("POST /v1/decks/{deck_id}/draw/{number}", logRequests(handlePostDeckDraw(ds)))
}
