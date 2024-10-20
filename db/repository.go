package repository

import (
	m "Lembrete/models"
	"database/sql"
	"fmt"
	"strconv"

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
		return fmt.Errorf("failed to create table: %v", err)
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
		return fmt.Errorf("failed to create table: %v", err)
	}

	return nil
}

func OpenDB(dbName string) (*sql.DB, error) {
	dbName = "./" + dbName + ".db"
	database, err := sql.Open("sqlite", dbName)
	if err != nil {
		return nil, fmt.Errorf("failed to open or create database %w", err)
	}
	if err = database.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %v", err)
	}
	return database, nil
}

func InsertCard(db *sql.DB, card m.Flashcard) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	insertQuery := "INSERT INTO Cards (Front, Back, EaseFactor, Repetitions, Interval, NextReview, DeckID) VALUES (?, ?, ?, ?, ?, ?, ?)"

	_, err := db.Exec(insertQuery, card.Front, card.Back, card.EaseFactor, card.Repetitions, card.Interval, card.NextReview.Format("2006-01-02"), card.DeckID)
	if err != nil {
		return err
	}
	fmt.Println("Inserted card:", card.Front)
	return nil
}

func InsertDeck(db *sql.DB, deck m.Deck) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	insertQuery := "INSERT INTO decks (MaxNewCards, MaxReviewsDaily, Name) VALUES (?, ?, ?)"

	_, err := db.Exec(insertQuery, deck.MaxNewCards, deck.MaxReviewsDaily, deck.Name)
	if err != nil {
		return err
	}

	fmt.Println("Inserted deck", deck.Name)
	return nil
}

func FetchAllCards(db *sql.DB, deckId int32) ([]m.Flashcard, error) {
	var cards []m.Flashcard
	rows, err := db.Query("SELECT * FROM Cards WHERE DeckID= ?", strconv.Itoa(int(deckId)))
	if err != nil {
		return nil, fmt.Errorf("the query that selects all cards has failed")
	}
	defer rows.Close()

	for rows.Next() {
		var card m.Flashcard
		if err := rows.Scan(
			&card.Id,
			&card.Front,
			&card.Back,
			&card.EaseFactor,
			&card.Repetitions,
			&card.Interval,
			&card.NextReview,
			&card.DeckID); err != nil {
			return nil, fmt.Errorf("failed to scan rows of the card")
		}
		cards = append(cards, card)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("there was an error when iterating over the rows: %v", err)
	}
	if len(cards) == 0 {
		return nil, fmt.Errorf("the list of decks is empty")
	}
	return cards, nil
}

func FetchAllDecks(db *sql.DB) ([]m.Deck, error) {

	var Decks []m.Deck
	rows, err := db.Query("SELECT * FROM Decks")
	if err != nil {
		return nil, fmt.Errorf("the query that was to select all decks has failed: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var deck m.Deck
		if err := rows.Scan(
			&deck.Id,
			&deck.Name,
			&deck.MaxNewCards,
			&deck.MaxReviewsDaily,
		); err != nil {
			return nil, fmt.Errorf("there was an error when scanning the decks: %v", err)
		}
		Decks = append(Decks, deck)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("there was an error when iterating over rows")
	}
	if len(Decks) == 0 {
		return nil, fmt.Errorf("the list of decks is empty")
	}
	return Decks, nil
}

func DisplayArrDecks(decks []m.Deck) error {
	for _, deck := range decks {
		fmt.Printf("Deck ID: %d, Name: %s, Max New Cards: %d, Max Reviews Daily: %d\n",
			deck.Id, deck.Name, deck.MaxNewCards, deck.MaxReviewsDaily)
	}
	return nil
}

func DisplayArrCards(cards []m.Flashcard) error {
	for _, card := range cards {
		fmt.Printf("Card ID: %d, Front: %s, Back: %s, Ease Factor: %.2f, Repetitions: %d, Interval: %.2f, Next Review: %s, Deck ID: %s\n",
			card.Id, card.Front, card.Back, card.EaseFactor, card.Repetitions, card.Interval, card.NextReview.Format("2006-01-02"), card.DeckID)
	}
	return nil
}

func UpdateCardRecords(db *sql.DB, cards []m.Flashcard) error {

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
		return fmt.Errorf("there ws an error when preparing the statement to update card records: %v", err)
	}
	defer stmt.Close()

	for _, card := range cards {
		_, err = stmt.Exec(card.Id, card.Front, card.Back, card.EaseFactor, card.Repetitions, card.Interval, card.NextReview.Format("2006-01-02"), card.DeckID)
		if err != nil {
			return fmt.Errorf("there was an error when executing the statement: %v", err)
		}
	}

	return nil
}

func UpdateDeckRecords(db *sql.DB, decks []m.Deck) error {
	stmt, err := db.Prepare(`
        UPDATE Decks 
        SET name = ?, 
            MaxNewCards = ?, 
            MaxReviewsDaily = ? 
        WHERE id = ?
    `)
	if err != nil {
		return fmt.Errorf("there was an error when preparing the statement to update deck records: %v", err)
	}
	defer stmt.Close()

	for _, deck := range decks {
		_, err = stmt.Exec(deck.Name, deck.MaxNewCards, deck.MaxReviewsDaily, deck.Id)
		if err != nil {
			return fmt.Errorf("there was an error when executing the statement: %v", err)
		}
	}

	return nil
}
