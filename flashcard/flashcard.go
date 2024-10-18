package flashcard

import (
	"database/sql"
	//"errors"
	"fmt"
	//"log"
	//"sync"

	_ "github.com/mattn/go-sqlite3"
)

type Deck struct {
	maxNewCards int32
	maxReviews  int32
}

func newDeck(maxNewCards int32, maxReviews int32) Deck {
	return Deck{
		maxNewCards: maxNewCards,
		maxReviews:  maxReviews}
}

type Flashcard struct {
	Front       string
	Back        string
	EaseFactor  float32
	Repetitions int
	Interval    float32
	NextReview  string
	DeckID      string
	CardId      string
}

func NewFlashcard(
	front string,
	back string,
	easeFactor float32,
	repetitions int,
	interval float32,
	nextReview string,
	DeckID string,
	CardId string,
) Flashcard {
	return Flashcard{
		Front:       front,
		Back:        back,
		EaseFactor:  easeFactor,
		Repetitions: repetitions,
		Interval:    interval,
		NextReview:  nextReview,
		DeckID:      DeckID,
		CardId:      CardId,
	}
}

func insertCard(card Flashcard) error {

	return nil
}

func openDB(dbName string) error {
	dbName = "./" + dbName + ".db"
	database, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return fmt.Errorf("failed to open or create database %w", err)
	}
	defer database.Close()

	if err = database.Ping(); err != nil {
		return fmt.Errorf("failed to connect to the database")
	}
	return nil
}

func mainData() {

}
