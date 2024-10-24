package repository
import (
	"database/sql"
	"fmt"

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