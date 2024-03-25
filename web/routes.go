package web

import (
	"log/slog"
	"net/http"

	"github.com/mtpereira/deck/deck"
)

func addRoutes(mux *http.ServeMux, log *slog.Logger, da *deck.DeckAPI) {
	logRequests := newLoggerMiddleware(log)
	mux.Handle("POST /v1/decks", logRequests(handlePostDeck(da)))
	mux.Handle("GET /v1/decks/{deck_id}", logRequests(handleGetDeck(da)))
	mux.Handle("POST /v1/decks/{deck_id}/draw/{number}", logRequests(handlePostDeckDraw(da)))
}
