package web

import (
	"log/slog"
	"net/http"
)

func addRoutes(mux *http.ServeMux, log *slog.Logger) {
	logRequests := newLoggerMiddleware(log)
	mux.Handle("POST /v1/decks/", logRequests(handlePostDeck()))
}
