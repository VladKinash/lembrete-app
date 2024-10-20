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
	

	decks, err := repo.FetchAllDecks(db)
	if err != nil {
		log.Fatal(err)
	}
	repo.DisplayArrDecks(decks)

}
