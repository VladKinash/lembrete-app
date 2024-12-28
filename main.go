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
		return
	}
	defer db.Close()

	decks, err := repo.FetchAllDecks(db)
	if err != nil {
		fmt.Println("error fetching decks: ", err)
		return
	}

	myApp := app.New()
	myWindow := myApp.NewWindow("Lembrete")

	gui.CreateDecksUI(decks, db, myApp, myWindow)

	myWindow.Resize(fyne.NewSize(600, 400))
	myWindow.ShowAndRun()
}
