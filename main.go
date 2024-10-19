package main

import (
	repo "Lembrete/db"
	models "Lembrete/models"
	"fmt"
	"log"
	"time"
)

func main() {

	// Create sample decks
	deck1 := models.NewDeck(10, 20, "Spanish Vocabulary", 1)

	// Create sample flashcards
	flashcard1 := models.NewFlashcard(
		"Hola",                      // Front
		"Hello",                     // Back
		2.5,                         // EaseFactor
		5,                           // Repetitions
		1.0,                         // Interval
		time.Now().AddDate(0, 0, 2), // NextReview (2 days from now)
		"1",                         // DeckID (as string)
		1,                           // Id
	)

	flashcard2 := models.NewFlashcard(
		"Gracias",                   // Front
		"Thank you",                 // Back
		3.0,                         // EaseFactor
		3,                           // Repetitions
		1.5,                         // Interval
		time.Now().AddDate(0, 0, 1), // NextReview (1 day from now)
		"1",                         // DeckID (as string)
		2,                           // Id
	)

	flashcard3 := models.NewFlashcard(
		"¿Cómo estás?",              // Front
		"How are you?",              // Back
		2.0,                         // EaseFactor
		4,                           // Repetitions
		2.0,                         // Interval
		time.Now().AddDate(0, 0, 3), // NextReview (3 days from now)
		"1",                         // DeckID (as string)
		3,                           // Id
	)

	// Print the deck and flashcards
	fmt.Printf("Deck: %+v\n", deck1)
	fmt.Printf("Flashcard 1: %+v\n", flashcard1)
	fmt.Printf("Flashcard 2: %+v\n", flashcard2)
	fmt.Printf("Flashcard 3: %+v\n", flashcard3)
	fmt.Println("Hello world!")

	db, err := repo.OpenDB("your_db_name") // Open the database
	if err != nil {
		log.Fatal(err) // Handle any error
	}
	defer db.Close()

	

}
