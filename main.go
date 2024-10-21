package main

import (
	repo "Lembrete/db"
	//	model "Lembrete/models"
	"fmt"
	"log"
	//	"time"
)

func main() {

	db, err := repo.OpenDB("your_db_name")
	if err != nil {
		fmt.Println("error fetching decks: ", err)
	}
	defer db.Close()
	

	cards, err := repo.FetchAllCards(db, 1)
	if err != nil {
		log.Fatal(err)
	}
	repo.DisplayArrCards(cards)

}
