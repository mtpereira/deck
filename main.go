package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ardanlabs/conf/v3"
	"github.com/google/uuid"
)

func main() {
	ctx := context.Background()
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	err := run(ctx, log)
	if err != nil {
		log.Error(error.Error(err))
	}
}

func run(ctx context.Context, log *slog.Logger) error {
	cfg := struct {
		APIHost string `conf:"default:127.0.0.1:9000"`
	}{}
	prefix := "DECK"
	help, err := conf.Parse(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}
		return fmt.Errorf("error parsing config: %w", err)
	}

	out, err := conf.String(&cfg)
	if err != nil {
		return fmt.Errorf("error outputting conf: %w", err)
	}
	log.Info("startup", "configs", out)

	log.Info("startup", "status", "initalising api")
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	mux := http.NewServeMux()
	addRoutes(mux, log)

	api := http.Server{
		Handler:  mux,
		Addr:     cfg.APIHost,
		ErrorLog: slog.NewLogLogger(slog.NewTextHandler(os.Stdout, nil), slog.LevelError),
	}

	serverErrors := make(chan error, 1)
	go func() {
		log.Info("startup", "status", "api listening", "host", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case sig := <-shutdown:
		log.Info("shutdown", "status", "initiated", "signal", sig)
		log.Info("shutdown", "status", "complete")

		ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()

		err := api.Shutdown(ctx)
		if err != nil {
			api.Close()
			return fmt.Errorf("could not shutdown gracefully: %w", err)
		}
	}

	return nil
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

func addRoutes(mux *http.ServeMux, log *slog.Logger) {
	logRequests := newLoggerMiddleware(log)
	mux.Handle("POST /v1/decks/", logRequests(handlePostDeck()))
}

func handlePostDeck() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		encodeJSON(w, http.StatusOK, createDeck())
	})
}

func createDeck() deck {
	return deck{
		DeckID:    uuid.New(),
		Shuffled:  true,
		Remaining: 52,
	}
}

type deck struct {
	DeckID    uuid.UUID `json:"deck_id"`
	Shuffled  bool      `json:"shuffled"`
	Remaining int       `json:"remaining"`
	Cards     []card    `json:"cards,omitempty"`
}

type card struct {
	Code  string `json:"code"`
	Value string `json:"value,omitempty"`
	Suit  string `json:"suit,omitempty"`
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
