package main

import (
	repo "Lembrete/db"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	//models "Lembrete/models"
	"fmt"
	//"log"
	//"time"
	gui "Lembrete/gui"
)

func main() {

	db, err := repo.OpenDB("your_db_name")
	if err != nil {
		fmt.Println("error fetching decks: ", err)
	}
	defer db.Close()

	decks, err := repo.FetchAllDecks(db)

	myApp := app.New()
	myWindow := myApp.NewWindow("Anki Clone")

	deckUI := gui.CreateDecksUI(decks)
	myWindow.SetContent(deckUI)

	myWindow.Resize(fyne.NewSize(400, 600))
	myWindow.ShowAndRun()

}
