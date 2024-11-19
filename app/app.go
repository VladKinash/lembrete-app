package app

import(
	models "Lembrete/models"
	repo "Lembrete/db"
	"database/sql"
	_ "modernc.org/sqlite"
	"fmt"
)

/*db, err := repo.OpenDB("default")
if err != nil{
	return fmt.Errorf("there was an error when opening the database: %v", err)
} */

func displayDecks(db *sql.DB) error{

	decks, err := repo.FetchAllDecks(db)
	if err != nil{
		return fmt.Errorf("failed to fetch all decks")
	}

	repo.DisplayArrDecks(decks)

	return nil
}

func cardReview(db *sql.DB, card models.Flashcard) error{



	return nil
}