package main

import (
	repo "Lembrete/db"
	"fmt"
	//"log"
)

func main() {

	db, err := repo.OpenDB("your_db_name")
	if err != nil {
		fmt.Println("error fetching decks: ", err)
	}
	defer db.Close()
	

	card, err := repo.FetchCard(db, 1)

	fmt.Println(card)
	fmt.Println(err)

}
