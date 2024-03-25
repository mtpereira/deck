package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/mtpereira/deck/deck"
)

func NewMux(log *slog.Logger, da *deck.DeckAPI) *http.ServeMux {
	mux := http.NewServeMux()
	addRoutes(mux, log, da)
	return mux
}

type wrappedWriter struct {
	responseWriter http.ResponseWriter
	statusCode     int
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.responseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func newLoggerMiddleware(log *slog.Logger) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			wr := wrappedWriter{
				responseWriter: w,
				statusCode:     http.StatusOK,
			}
			h.ServeHTTP(wr.responseWriter, r)
			log.Info("api", "request", r.URL.Path, "status", wr.statusCode, "duration", time.Since(start))
		})
	}
}

func handlePostDeck(da *deck.DeckAPI) http.Handler {
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

		d := da.New(shuffled, nil)
		encodeJSON(w, http.StatusOK, deckResponse{
			DeckID:    d.DeckID,
			Shuffled:  d.Shuffled,
			Remaining: d.Remaining,
		})
	})
}

func handleGetDeck(da *deck.DeckAPI) http.Handler {
	type deckResponse struct {
		DeckID    uuid.UUID   `json:"deck_id"`
		Shuffled  bool        `json:"shuffled"`
		Remaining int         `json:"remaining"`
		Cards     []deck.Card `json:"cards"`
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		deckIDParam := r.PathValue("deck_id")
		deckID, err := uuid.Parse(deckIDParam)
		if err != nil {
			encodeJSON(w, http.StatusBadRequest, respondError(http.StatusBadRequest, "Invalid deck ID"))
			return
		}

		d, err := da.Get(deckID)
		if err != nil {
			encodeJSON(w, http.StatusNotFound, respondError(http.StatusNotFound, err.Error()))
			return
		}
		encodeJSON(w, http.StatusOK, deckResponse(*d))
	})
}

func handlePostDeckDraw(da *deck.DeckAPI) http.Handler {
	type cardsResponse struct {
		Cards []deck.Card `json:"cards"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		deckIDParam := r.PathValue("deck_id")
		deckID, err := uuid.Parse(deckIDParam)
		if err != nil {
			encodeJSON(w, http.StatusBadRequest, respondError(http.StatusBadRequest, "Invalid deck ID"))
			return
		}

		numberParam := r.PathValue("number")
		cardsToDraw, err := strconv.Atoi(numberParam)
		if err != nil {
			encodeJSON(w, http.StatusBadRequest, respondError(http.StatusBadRequest, "Invalid number of cards to draw"))
			return
		}
		if cardsToDraw < 0 || cardsToDraw > 52 {
			encodeJSON(w, http.StatusBadRequest, respondError(http.StatusBadRequest, "Invalid number of cards to draw"))
			return
		}

		cards, err := da.Draw(deckID, cardsToDraw)
		if err != nil {
			if errors.Is(err, deck.ErrUnsufficientCards) {
				encodeJSON(w, http.StatusBadRequest, respondError(http.StatusBadRequest, err.Error()))
				return
			}
			encodeJSON(w, http.StatusInternalServerError, respondError(http.StatusInternalServerError, err.Error()))
			return
		}

		encodeJSON(w, http.StatusOK, cardsResponse{Cards: cards})
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
