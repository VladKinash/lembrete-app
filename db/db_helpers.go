package repository

import (
	models "Lembrete/models"
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

func CreateTableDeck(db *sql.DB) error {
	createTableQuery := `CREATE TABLE IF NOT EXISTS Decks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		MaxNewCards INTEGER,
		MaxReviewsDaily INTEGER
	);`

	_, err := db.Exec(createTableQuery)
	if err != nil {
		return fmt.Errorf("failed to create Decks table: %v", err)
	}
	return nil
}

func CreateTableCard(db *sql.DB) error {
	createTableQuery := `CREATE TABLE IF NOT EXISTS Cards (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		Front TEXT,
		Back TEXT,
		EaseFactor REAL,
		Repetitions INTEGER,
		Interval REAL,
		NextReview TEXT,
		DeckID INTEGER,
		FOREIGN KEY (DeckID) REFERENCES Decks(id)
	);`

	_, err := db.Exec(createTableQuery)
	if err != nil {
		return fmt.Errorf("failed to create Cards table: %v", err)
	}
	return nil
}

func OpenDB(dbName string) (*sql.DB, error) {
	dbName = "./" + dbName + ".db"
	database, err := sql.Open("sqlite", dbName)
	if err != nil {
		return nil, fmt.Errorf("failed to open or create database: %w", err)
	}
	if err = database.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %v", err)
	}
	return database, nil
}

func scanDeckRow(row *sql.Rows) (models.Deck, error) {
	var deck models.Deck
	if err := row.Scan(
		&deck.ID,
		&deck.Name,
		&deck.MaxNewCards,
		&deck.MaxReviewsDaily); err != nil {
		return deck, fmt.Errorf("failed to scan deck row: %v", err)
	}
	return deck, nil
}

func scanFlashcardRow(row *sql.Rows) (models.Flashcard, error) {
	var card models.Flashcard
	var nextReviewStr string

	if err := row.Scan(
		&card.ID,
		&card.Front,
		&card.Back,
		&card.EaseFactor,
		&card.Repetitions,
		&card.Interval,
		&nextReviewStr,
		&card.DeckID); err != nil {
		return card, fmt.Errorf("failed to scan card row: %v", err)
	}
	convReview, err := time.Parse("2006-01-02", nextReviewStr)
	card.NextReview = convReview
	if err != nil {
		return card, fmt.Errorf("failed to parse NextReview date: %v", err)
	}

	return card, nil
}
