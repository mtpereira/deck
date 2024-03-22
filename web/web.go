package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/mtpereira/deck/deck"
)

func NewMux(log *slog.Logger) *http.ServeMux {
	mux := http.NewServeMux()
	addRoutes(mux, log)
	return mux
}

func newLoggerMiddleware(log *slog.Logger) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Info("api", "request", "started handling", "path", r.URL.Path)
			h.ServeHTTP(w, r)
			log.Info("api", "request", "finished handling", "path", r.URL.Path)
		})
	}
}

func handlePostDeck() http.Handler {
	type deckResponse struct {
		DeckID    uuid.UUID `json:"deck_id"`
		Shuffled  bool      `json:"shuffled"`
		Remaining int       `json:"remaining"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		shuffledParam := r.URL.Query().Get("shuffled")
		shuffled := false
		err := getParam(&shuffled, shuffledParam, "true", "false")
		if err != nil {
			encodeJSON(w, http.StatusBadRequest, respondError(http.StatusBadRequest, "Invalid shuffled parameter"))
			return
		}
		d := deck.New(shuffled)
		encodeJSON(w, http.StatusOK, deckResponse{
			DeckID:    d.DeckID,
			Shuffled:  d.Shuffled,
			Remaining: d.Remaining,
		})
	})
}

func getParam(param any, paramString string, validValues ...string) error {
	if paramString == "" {
		return nil
	}

	for _, value := range validValues {
		if value == paramString {
			n, err := fmt.Sscanf(paramString, "%v", param)
			if n != 1 || err != nil {
				return err
			}
			return nil
		}
	}

	return errors.New("Couldn't parse query parameter")
}

func encodeJSON[T any](w http.ResponseWriter, status int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		return fmt.Errorf("encode json: %w", err)
	}
	return nil
}

func respondError(errorCode int, message string) errorResponse {
	return errorResponse{
		Code:    errorCode,
		Message: message,
	}
}

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
