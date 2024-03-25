package web

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"sync"
	"testing"

	"github.com/mtpereira/deck/deck"
)

func Test_handlePostDeckDraw(t *testing.T) {
	// Test it correctly draws a card.
	var deckID string
	ds := deck.NewStore(slog.New(slog.NewTextHandler(io.Discard, nil)))
	d := deck.New(false, nil)
	ds.Create(d)
	deckID = d.DeckID.String()
	cardsToDraw := 1

	req, err := http.NewRequest("POST", fmt.Sprintf("/v1/decks/%s/draw/%d", deckID, cardsToDraw), nil)
	if err != nil {
		t.Fatalf(err.Error())
	}

	rr := httptest.NewRecorder()
	handler := http.NewServeMux()
	handler.Handle("POST /v1/decks/{deck_id}/draw/{number}", handlePostDeckDraw(ds))

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %v", rr.Code)
	}

	expectedRemaining := d.Remaining - cardsToDraw
	c, err := decodeCards(rr.Body)
	if err != nil {
		t.Errorf("Expected to get Cards, got %v, and the error %v", c, err)
	}
	u, err := ds.QueryById(d.DeckID)
	if err != nil {
		t.Errorf("Deck missing from store after update: %v", err)
	}
	if u.Remaining != expectedRemaining {
		t.Errorf("Expected %d cards to remain on the deck, got %d", expectedRemaining, d.Remaining)
	}
	if !slices.Contains(c.Cards, deck.Card{Code: "2C"}) {
		t.Errorf("Expected drawn card to be 2C, got %v", c)
	}

	// Test it fails to draw more cards than the ones available.
	cardsToDraw = 52
	req, err = http.NewRequest("POST", fmt.Sprintf("/v1/decks/%s/draw/%d", deckID, cardsToDraw), nil)
	if err != nil {
		t.Fatalf(err.Error())
	}

	rr = httptest.NewRecorder()
	handler = http.NewServeMux()
	handler.Handle("POST /v1/decks/{deck_id}/draw/{number}", handlePostDeckDraw(ds))

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 Bad Request, got %v", rr.Code)
	}

	// Test it correctly handles concurrent draw requests.
	e := deck.New(false, nil)
	ds.Create(e)
	deckID = e.DeckID.String()

	cardsToDraw = 52
	numReqs := 8
	recs := []*httptest.ResponseRecorder{}
	for range numReqs {
		recs = append(recs, httptest.NewRecorder())
	}
	wg := sync.WaitGroup{}
	wg.Add(numReqs)
	for _, rr := range recs {
		handler = http.NewServeMux()
		handler.Handle("POST /v1/decks/{deck_id}/draw/{number}", handlePostDeckDraw(da))
		req, err = http.NewRequest("POST", fmt.Sprintf("/v1/decks/%s/draw/%d", deckID, cardsToDraw), nil)
		if err != nil {
			t.Fatalf(err.Error())
		}

		go func() {
			handler.ServeHTTP(rr, req)
			defer wg.Done()
		}()
	}
	wg.Wait()

	codes := []int{}
	cards := []cardsResponse{}
	for _, rr := range recs {
		c, err := decodeCards(rr.Body)
		if err != nil {
			t.Errorf("Expected to get Cards, got %v, and the error %v", c, err)
		}
		cards = append(cards, c)
		codes = append(codes, rr.Code)
	}
	badReqsCount := 0
	for _, c := range codes {
		if c == 400 {
			badReqsCount += 1
		}
	}
	if badReqsCount != numReqs-1 {
		t.Errorf("All but one requests should have 400ed, %d out of %d did", badReqsCount, numReqs-1)
	}
	reqsWithCardsCount := 0
	for _, c := range cards {
		if len(c.Cards) > 0 {
			reqsWithCardsCount += 1
		}
	}
	if reqsWithCardsCount != 1 {
		t.Errorf("All but one requests should have cards, %d did", reqsWithCardsCount)
	}

}

func Test_handlePostDeck(t *testing.T) {
	// Test it correctly handles the shuffled param when it is defined.
	shuffledParamValues := []string{
		"",
		"true",
		"false",
	}
	ds := deck.NewStore(slog.New(slog.NewTextHandler(io.Discard, nil)))
	for _, shuffledParam := range shuffledParamValues {
		req, err := http.NewRequest("POST", fmt.Sprintf("/v1/decks?shuffled=%s", shuffledParam), nil)
		if err != nil {
			t.Fatalf(err.Error())
		}

		rr := httptest.NewRecorder()
		handler := http.Handler(handlePostDeck(da))

		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected 200 OK, got %v", rr.Code)
		}

		d, err := decodeDeck(rr.Body)
		if err != nil {
			t.Errorf("Expected to get a Deck, got %v", d)
		}
	}

	// Test it correctly handles the shuffled param when it is undefined.
	req, err := http.NewRequest("POST", "/v1/decks", nil)
	if err != nil {
		t.Fatalf(err.Error())
	}

	rr := httptest.NewRecorder()
	handler := http.Handler(handlePostDeck(da))

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %v", rr.Code)
	}

	d, err := decodeDeck(rr.Body)
	if err != nil {
		t.Errorf("Expected to get a Deck, got %v", d)
	}

	// Test it returns a 400 when the shuffled param is invalid.
	req, err = http.NewRequest("POST", "/v1/decks?shuffled=asdf", nil)
	if err != nil {
		t.Fatalf(err.Error())
	}

	rr = httptest.NewRecorder()
	handler = http.Handler(handlePostDeck(da))

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 Bad Request, got %v", rr.Code)
	}

	expectedResponse := `{"code":400,"message":"Invalid shuffled parameter"}`
	if expectedResponse != strings.TrimRight(rr.Body.String(), "\n") {
		t.Errorf("Expected to get %v, got %v", expectedResponse, rr.Body.String())
	}
}

func Test_handleGetDeck(t *testing.T) {
	// Test it returns an existing deck.
	var deckID string
	da := deck.NewAPI(slog.New(slog.NewTextHandler(io.Discard, nil)))
	d := da.New(true, nil)
	deckID = d.DeckID.String()

	req, err := http.NewRequest("GET", fmt.Sprintf("/v1/decks/%s", deckID), nil)
	if err != nil {
		t.Fatalf(err.Error())
	}

	rr := httptest.NewRecorder()
	handler := http.NewServeMux()
	handler.Handle("GET /v1/decks/{deck_id}", handleGetDeck(da))

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %v", rr.Code)
	}

	dr, err := decodeDeck(rr.Body)
	if err != nil {
		t.Errorf("Expected to get a Deck, got %v", dr)
	}
	if dr.DeckID != d.DeckID {
		t.Errorf("Expected deck %v, got deck %v", d, dr)
	}

	// Test it handles correctly when the deck does not exist.
	unexistingDeckID := "14ca6cac-e933-4484-8e3f-e5acd505d11d"

	req, err = http.NewRequest("GET", fmt.Sprintf("/v1/decks/%s", unexistingDeckID), nil)
	if err != nil {
		t.Fatalf(err.Error())
	}

	rr = httptest.NewRecorder()
	handler = http.NewServeMux()
	handler.Handle("GET /v1/decks/{deck_id}", handleGetDeck(da))

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("Expected 404, got %v", rr.Code)
	}

	er, err := decodeErrorResponse(rr.Body)
	if err != nil {
		t.Errorf("Expected to get a error response, got %v", er)
	}
	if er.Code != rr.Code {
		t.Errorf("Expected %v, got %v", rr.Code, er.Code)
	}
	if strings.Contains(rr.Body.String(), "deck_id") {
		t.Errorf("Error response should not have a deck, got %v", rr.Body.String())
	}
}

func decodeDeck(b io.Reader) (deck.Deck, error) {
	var d deck.Deck
	if err := json.NewDecoder(b).Decode(&d); err != nil {
		return d, fmt.Errorf("decode json: %w", err)
	}
	return d, nil
}

func decodeErrorResponse(b io.Reader) (errorResponse, error) {
	var d errorResponse
	if err := json.NewDecoder(b).Decode(&d); err != nil {
		return d, fmt.Errorf("decode json: %w", err)
	}
	return d, nil
}

func decodeCards(b io.Reader) (cardsResponse, error) {
	var c cardsResponse
	if err := json.NewDecoder(b).Decode(&c); err != nil {
		return c, fmt.Errorf("decode json: %w", err)
	}
	return c, nil
}

type cardsResponse struct {
	Cards []deck.Card `json:"cards"`
}
