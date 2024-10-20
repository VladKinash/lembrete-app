package repository

import (
	models "Lembrete/models" // Use a descriptive alias
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

func InsertCard(db *sql.DB, card models.Flashcard) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	insertQuery := "INSERT INTO Cards (Front, Back, EaseFactor, Repetitions, Interval, NextReview, DeckID) VALUES (?, ?, ?, ?, ?, ?, ?)"
	_, err := db.Exec(insertQuery, card.Front, card.Back, card.EaseFactor, card.Repetitions, card.Interval, card.NextReview.Format("2006-01-02"), card.DeckID)
	if err != nil {
		return fmt.Errorf("failed to insert card: %v", err)
	}

	fmt.Println("Inserted card:", card.Front)
	return nil
}

func InsertDeck(db *sql.DB, deck models.Deck) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	insertQuery := "INSERT INTO Decks (MaxNewCards, MaxReviewsDaily, Name) VALUES (?, ?, ?)"
	_, err := db.Exec(insertQuery, deck.MaxNewCards, deck.MaxReviewsDaily, deck.Name)
	if err != nil {
		return fmt.Errorf("failed to insert deck: %v", err)
	}

	fmt.Println("Inserted deck:", deck.Name)
	return nil
}

func FetchAllCards(db *sql.DB, deckID int32) ([]models.Flashcard, error) {
	var cards []models.Flashcard
	rows, err := db.Query("SELECT * FROM Cards WHERE DeckID = ?", deckID)
	if err != nil {
		return nil, fmt.Errorf("failed to select cards for deck %d: %v", deckID, err)
	}
	defer rows.Close()

	var nextReviewStr string
	for rows.Next() {
		var card models.Flashcard
		if err := rows.Scan(
			&card.ID,
			&card.Front,
			&card.Back,
			&card.EaseFactor,
			&card.Repetitions,
			&card.Interval,
			&nextReviewStr,
			&card.DeckID); err != nil {
			return nil, fmt.Errorf("failed to scan card row: %v", err)
		}
		card.NextReview, err = time.Parse("2006-01-02", nextReviewStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse NextReview date: %v", err)
		}
		cards = append(cards, card)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred while iterating over rows: %v", err)
	}

	return cards, nil
}

func FetchAllDecks(db *sql.DB) ([]models.Deck, error) {
	var decks []models.Deck
	rows, err := db.Query("SELECT * FROM Decks")
	if err != nil {
		return nil, fmt.Errorf("failed to select all decks: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var deck models.Deck
		if err := rows.Scan(
			&deck.ID,
			&deck.Name,
			&deck.MaxNewCards,
			&deck.MaxReviewsDaily); err != nil {
			return nil, fmt.Errorf("failed to scan deck row: %v", err)
		}
		decks = append(decks, deck)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred while iterating over rows: %v", err)
	}

	return decks, nil
}

func DisplayArrDecks(decks []models.Deck) error {
	for _, deck := range decks {
		fmt.Printf("Deck ID: %d, Name: %s, Max New Cards: %d, Max Reviews Daily: %d\n",
			deck.ID, deck.Name, deck.MaxNewCards, deck.MaxReviewsDaily)
	}
	return nil
}

func DisplayArrCards(cards []models.Flashcard) error {
	for _, card := range cards {
		fmt.Printf("Card ID: %d, Front: %s, Back: %s, Ease Factor: %.2f, Repetitions: %d, Interval: %.2f, Next Review: %s, Deck ID: %s\n",
			card.ID, card.Front, card.Back, card.EaseFactor, card.Repetitions, card.Interval, card.NextReview.Format("2006-01-02"), card.DeckID)
	}
	return nil
}

func UpdateCardRecords(db *sql.DB, cards []models.Flashcard) error {
	stmt, err := db.Prepare(`
	UPDATE Cards 
	SET Front = ?, 
		Back = ?, 
		EaseFactor = ?, 
		Repetitions = ?, 
		Interval = ?, 
		NextReview = ?, 
		DeckID = ? 
	WHERE ID = ?
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare card update statement: %v", err)
	}
	defer stmt.Close()

	for _, card := range cards {
		_, err = stmt.Exec(card.Front, card.Back, card.EaseFactor, card.Repetitions, card.Interval, card.NextReview.Format("2006-01-02"), card.DeckID, card.ID)
		if err != nil {
			return fmt.Errorf("failed to execute card update statement: %v", err)
		}
	}

	return nil
}

func UpdateDeckRecords(db *sql.DB, decks []models.Deck) error {
	stmt, err := db.Prepare(`
        UPDATE Decks 
        SET name = ?, 
            MaxNewCards = ?, 
            MaxReviewsDaily = ? 
        WHERE id = ?
    `)
	if err != nil {
		return fmt.Errorf("failed to prepare deck update statement: %v", err)
	}
	defer stmt.Close()

	for _, deck := range decks {
		_, err = stmt.Exec(deck.Name, deck.MaxNewCards, deck.MaxReviewsDaily, deck.ID)
		if err != nil {
			return fmt.Errorf("failed to execute deck update statement: %v", err)
		}
	}

	return nil
}
