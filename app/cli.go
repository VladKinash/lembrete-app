package app

import (
	repo "Lembrete/db"
	models "Lembrete/models"
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"

	_ "modernc.org/sqlite"
)

func mainMenu(db *sql.DB) error{

	reader := bufio.NewReader(os.Stdin)

	decks, err := repo.FetchAllDecks(db)
	if err != nil{
		return fmt.Errorf("failed to fetch all decks")
	}

	for {
		fmt.Println("Available Decks:")
		for i, deck := range decks {
			fmt.Printf("%d. %s\n", i+1, deck.Name)
		}
		fmt.Println("0. Exit")

		fmt.Print("Select a deck by number: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		input = strings.TrimSpace(input)
		choice, err := strconv.Atoi(input)
		if err != nil || choice < 0 || choice > len(decks) {
			fmt.Println("Invalid choice. Please enter a valid number.")
			continue
		}

		if choice == 0{
			fmt.Println("Bye!")
			return nil
		}
		

	}
}

func deckMenu(db *sql.DB, deck models.Deck) error{


	return nil
}