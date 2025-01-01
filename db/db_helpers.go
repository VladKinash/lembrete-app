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

func OpenAndInitializeDB(dbName string) (*sql.DB, error) {

	dbName = "./" + dbName + ".db"

	database, err := sql.Open("sqlite", dbName)
	if err != nil {
		return nil, fmt.Errorf("failed to open or create database, %w", err)
	}
	if err = database.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to the database")
	}

	if err := CreateTableDeck(database); err != nil {
		return nil, fmt.Errorf("failed to create TableDeck: %v", err)
	}

	if err := CreateTableCard(database); err != nil {
		return nil, fmt.Errorf("failed to create Cards table: %v", err)
	}

	decks, err := FetchAllDecks(database)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch decks: %v", err)
	}

	if len(decks) == 0 {
		defaultDeck := models.NewDeck(5, 20, "Default", 0)
		if err := InsertDeck(database, defaultDeck); err != nil {
			return nil, fmt.Errorf("failed to create default deck: %v", err)
		}
		fmt.Println("Created default deck:", defaultDeck.Name)
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

func CountCards(db *sql.DB, deckID int32) (newCards int, dueCards int, err error) {
	err = db.QueryRow("SELECT COUNT(*) FROM Cards WHERE DeckID = ? AND Repetitions = 0", deckID).Scan(&newCards)
	if err != nil {
		return 0, 0, fmt.Errorf("error counting new cards: %v", err)
	}

	today := time.Now().Format("2006-01-02")
	err = db.QueryRow("SELECT COUNT(*) FROM Cards WHERE DeckID = ? AND NextReview <= ?", deckID, today).Scan(&dueCards)
	if err != nil {
		return 0, 0, fmt.Errorf("error counting due cards: %v", err)
	}

	return newCards, dueCards, nil
}
