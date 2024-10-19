package repository

import (
	m "Lembrete/models"
	"database/sql"
	"fmt"

	//"log"
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
