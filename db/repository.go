package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func createTableDeck(db *sql.DB) error {

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

func createTableCard(db *sql.DB) error {

	createTableQuery := `CREATE TABLE IF NOT EXISTS Deck (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		MaxNewCards INTEGER,
		MaxReviewsDaily INTEGER
	);`

	/*	front string,
		back string,
		easeFactor float32,
		repetitions int,
		interval float32,
		nextReview string,
		DeckID string,
		CardId string,*/

	_, err := db.Exec(createTableQuery)

	if err != nil {
		return fmt.Errorf("failed to create table: %v", err)
	}

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
