package repository_test

import (
	repository "Lembrete/db"
	models "Lembrete/models"
	"database/sql"
	"strconv"
	"testing"

	//"fmt"
	"time"

	_ "modernc.org/sqlite"
)

func openDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}
	return db
}

func setupTestDB(t *testing.T) *sql.DB {
	db := openDB(t)

	if err := repository.CreateTableCard(db); err != nil {
		t.Fatalf("failed to create Cards table: %v", err)
	}
	if err := repository.CreateTableDeck(db); err != nil {
		t.Fatalf("failed to create Decks table: %v", err)
	}

	decks := []models.Deck{
		{ID: 1, Name: "Deck 1", MaxNewCards: 20, MaxReviewsDaily: 50},
		{ID: 2, Name: "Deck 2", MaxNewCards: 15, MaxReviewsDaily: 30},
	}

	for _, deck := range decks {
		if err := repository.InsertDeck(db, deck); err != nil {
			t.Fatalf("failed to insert deck %v: %v", deck.Name, err)
		}
	}

	for i := 1; i <= 10; i++ {
		deckID := "1"
		if i > 5 {
			deckID = "2"
		}
		card := models.Flashcard{
			Front:       "Front Text " + strconv.Itoa(i),
			Back:        "Back Text " + strconv.Itoa(i),
			EaseFactor:  2.5,
			Repetitions: i,
			Interval:    float32(i),
			NextReview:  time.Now(),
			DeckID:      deckID,
		}

		if err := repository.InsertCard(db, card); err != nil {
			t.Fatalf("failed to insert card %v: %v", card.Front, err)
		}
	}

	return db
}

func TestInsertCard(t *testing.T) {

	db := openDB(t)
	defer db.Close()

	if err := repository.CreateTableCard(db); err != nil {
		t.Fatalf("failed to create Cards table: %v", err)
	}

	card := models.Flashcard{
		Front:       "Front Text",
		Back:        "Back Text",
		EaseFactor:  2.5,
		Repetitions: 1,
		Interval:    1.0,
		NextReview:  time.Now(),
		DeckID:      "1",
	}

	err := repository.InsertCard(db, card)
	if err != nil {
		t.Errorf("InsertCard() returned an error: %v", err)
	}

	cards, err := repository.FetchAllCards(db, 1)
	if err != nil {
		t.Fatalf("failed to fetch cards: %v", err)
	}

	if len(cards) != 1 || cards[0].Front != card.Front {
		t.Errorf("expected 1 card with Front %s, got %v", card.Front, cards)
	}
}

func TestFetchCard(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	card, err := repository.FetchCard(db, 1)
	if err != nil {
		t.Fatalf("unexpected error when fetching the card: %v", err)
	}
	if card.Front == "" {
		t.Errorf("expected card to have a front text, but got empty string")
	}

	missingCard, err := repository.FetchCard(db, 999)
	if err == nil {
		t.Errorf("expected an error for a missing card, but got card with front: %v", missingCard.Front)
	} else if err.Error() != "no card was found with ID 999" {
		t.Errorf("expected 'no card found' error, but got: %v", err)
	}
}

func TestFetchDeck(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	deck, err := repository.FetchDeck(db, 1)
	if err != nil {
		t.Fatalf("unexpected error when fetching the deck: %v", err)
	}
	if deck.Name == "" {
		t.Errorf("expected deck to have a name, but got an empty string")
	}
	if deck.MaxNewCards == 0 {
		t.Errorf("expected deck to have a MaxNewCards value, but got 0")
	}
	if deck.MaxReviewsDaily == 0 {
		t.Errorf("expected deck to have a MaxReviewsDaily value, but got 0")
	}

	missingDeck, err := repository.FetchDeck(db, 999)
	if err == nil {
		t.Errorf("expected an error for a missing deck, but got deck with name: %v", missingDeck.Name)
	} else if err.Error() != "no deck found with ID: 999" {
		t.Errorf("expected 'no deck found' error, but got: %v", err)
	}
}


func TestFetchAllCards(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	cards, err := repository.FetchAllCards(db, 1)
	if err != nil {
		t.Fatalf("unexpected error when fetching all cards for deck: %v", err)
	}
	if len(cards) != 5 {
		t.Errorf("expected 5 cards, got %d", len(cards))
	}
	for _, card := range cards {
		if card.Front == "" {
			t.Errorf("expected card to have front text, but got empty string")
		}
		if card.DeckID != "1" {
			t.Errorf("expected card to have DeckID 1, but got %v", card.DeckID)
		}
	}

	missingCards, err := repository.FetchAllCards(db, 999)
	if err != nil {
		t.Fatalf("unexpected error when fetching cards for non-existent deck: %v", err)
	}
	if len(missingCards) != 0 {
		t.Errorf("expected no cards for non-existent deck, got %d", len(missingCards))
	}
}

func TestFetchAllDecks(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	decks, err := repository.FetchAllDecks(db)
	if err != nil {
		t.Fatalf("unexpected error when fetching all decks: %v", err)
	}
	if len(decks) != 2 {
		t.Errorf("expected 2 decks, got %d", len(decks))
	}
	for _, deck := range decks {
		if deck.Name == "" {
			t.Errorf("expected deck to have a name, but got empty string")
		}
		if deck.MaxNewCards == 0 {
			t.Errorf("expected deck to have MaxNewCards value, but got 0")
		}
		if deck.MaxReviewsDaily == 0 {
			t.Errorf("expected deck to have MaxReviewsDaily value, but got 0")
		}
	}
}
