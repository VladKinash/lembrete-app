package repository

import (
	models "Lembrete/models"
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

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

	for rows.Next() {
		var card models.Flashcard
		card, err := scanFlashcardRow(rows)
		if err != nil {
			return cards, fmt.Errorf("something went wrong when scanning the cards")
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
		var card models.Deck
		card, err := scanDeckRow(rows)
		if err != nil {
			return decks, fmt.Errorf("something went wrong when scanning the decks")
		}
		decks = append(decks, card)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred while iterating over rows: %v", err)
	}

	return decks, nil
}

func FetchDeck(db *sql.DB, deckID int) (models.Deck, error) {
	var deck models.Deck
	var query, err = db.Query(`SELECT ID, Name, MaxNewCards, MaxReviewsDaily FROM Decks WHERE ID = ? LIMIT 1`, deckID)
	if err != nil {
		return deck, fmt.Errorf("there was an error when preparing FetchDeck query: %v", err)
	}
	defer query.Close()

	for query.Next() {
		return scanDeckRow(query)
	}

	return models.Deck{}, fmt.Errorf("no deck found with ID: %d", deckID)
}

func FetchCard(db *sql.DB, cardID int) (models.Flashcard, error) {
	var card models.Flashcard
	query, err := db.Query(`SELECT * FROM CARDS WHERE ID = ? LIMIT 1`, cardID)

	if err != nil {
		return card, fmt.Errorf("there was an error in query of fetchard: %v", err)
	}

	defer query.Close()

	for query.Next() {
		return scanFlashcardRow(query)
	}
	return models.Flashcard{}, fmt.Errorf("no card was found with ID %d", cardID)
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

func DeleteDeck(db *sql.DB, deck models.Deck) error {

	query, err := db.Prepare(`DELETE FROM Decks WHERE ID = ?`)
	if err != nil {
		return fmt.Errorf("there was an error when preparing the delete query: %v", err)
	}
	defer query.Close()

	_, err = query.Exec(deck.ID)
	if err != nil {
		return fmt.Errorf("failed to execute the delete query: %v", err)
	}

	return nil
}

func DeleteCard(db *sql.DB, card models.Flashcard) error {

	query, err := db.Prepare(`DELETE FROM Cards WHERE ID = ?`)
	if err != nil {
		return fmt.Errorf("there was an error when preparing the delete query: %v", err)
	}
	defer query.Close()

	_, err = query.Exec(card.ID)
	if err != nil {
		return fmt.Errorf("failed to execute the delete query: %v", err)
	}

	return nil
}
func UpdateCardRecord(db *sql.DB, card *models.Flashcard) error {
    fmt.Printf("UpdateCardRecord called with card: %+v\n", *card) 

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
        fmt.Printf("Error preparing update statement: %v\n", err)
        return fmt.Errorf("failed to prepare card update statement: %v", err)
    }
    defer stmt.Close()

    res, err := stmt.Exec(card.Front, card.Back, card.EaseFactor, card.Repetitions, card.Interval, card.NextReview.Format("2006-01-02"), card.DeckID, card.ID)
    if err != nil {
        fmt.Printf("Error executing update statement: %v\n", err)
        return fmt.Errorf("failed to execute card update statement: %v", err)
    }

    rowsAffected, err := res.RowsAffected()
    if err != nil {
        fmt.Printf("Error getting affected rows: %v\n", err)
        return fmt.Errorf("failed to get affected rows: %v", err)
    }

    fmt.Printf("Rows affected by update: %d\n", rowsAffected)

    if rowsAffected == 0 {
        return fmt.Errorf("no rows updated. check if card with ID %d exists", card.ID)
    }

    return nil
}
func UpdateDeckRecord(db *sql.DB, deck models.Deck) error {
	stmt, err := db.Prepare(`
        UPDATE Decks 
        SET Name = ?, 
            MaxNewCards = ?, 
            MaxReviewsDaily = ? 
        WHERE ID = ?
    `)
	if err != nil {
		return fmt.Errorf("failed to prepare deck update statement: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(deck.Name, deck.MaxNewCards, deck.MaxReviewsDaily, deck.ID)
	if err != nil {
		return fmt.Errorf("failed to execute deck update statement: %v", err)
	}

	return nil
}

func FetchNewCards(db *sql.DB, deckID int32, limit int32) ([]*models.Flashcard, error) {
	rows, err := db.Query("SELECT * FROM Cards WHERE DeckID = ? AND Repetitions = 0 LIMIT ?", deckID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to select new cards for deck %d: %v", deckID, err)
	}
	defer rows.Close()

	var cards []*models.Flashcard // Slice of pointers
	for rows.Next() {
		card, err := scanFlashcardRow(rows)
		if err != nil {
			return nil, fmt.Errorf("error scanning card: %v", err)
		}
		cards = append(cards, &card) // Append the address of card
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %v", err)
	}

	return cards, nil
}

func FetchDueCards(db *sql.DB, deckID int32, limit int32) ([]*models.Flashcard, error) {
	today := time.Now().Format("2006-01-02")
	rows, err := db.Query("SELECT * FROM Cards WHERE DeckID = ? AND NextReview <= ? LIMIT ?", deckID, today, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to select due cards for deck %d: %v", deckID, err)
	}
	defer rows.Close()

	var cards []*models.Flashcard // Slice of pointers
	for rows.Next() {
		card, err := scanFlashcardRow(rows)
		if err != nil {
			return nil, fmt.Errorf("error scanning card: %v", err)
		}
		cards = append(cards, &card) // Append the address of card
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %v", err)
	}

	return cards, nil
}
