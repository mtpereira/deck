package web

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mtpereira/deck/deck"
)

func Test_handlePostDeck(t *testing.T) {
	// Test it correctly handles the shuffled param when it is defined.
	shuffledParamValues := []string{
		"",
		"true",
		"false",
	}
	for _, shuffledParam := range shuffledParamValues {
		req, err := http.NewRequest("POST", fmt.Sprintf("/v1/decks?shuffled=%s", shuffledParam), nil)
		if err != nil {
			t.Fatalf(err.Error())
		}

		rr := httptest.NewRecorder()
		handler := http.Handler(handlePostDeck())

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
	handler := http.Handler(handlePostDeck())

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
	handler = http.Handler(handlePostDeck())

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 Bad Request, got %v", rr.Code)
	}

	expectedResponse := `{"code":400,"message":"Invalid shuffled parameter"}`
	if expectedResponse != strings.TrimRight(rr.Body.String(), "\n") {
		t.Errorf("Expected to get %v, got %v", expectedResponse, rr.Body.String())
	}

}

func decodeDeck(b io.Reader) (deck.Deck, error) {
	var d deck.Deck
	if err := json.NewDecoder(b).Decode(&d); err != nil {
		return d, fmt.Errorf("decode json: %w", err)
	}
	return d, nil
}
