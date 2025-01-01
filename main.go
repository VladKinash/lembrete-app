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
	db, err := repo.OpenAndInitializeDB("main")
	if err != nil {
		fmt.Println("Error initializing database:", err)
		return
	}
	defer db.Close()

	decks, err := repo.FetchAllDecks(db)
	if err != nil {
		fmt.Println("Error fetching decks:", err)
		return
	}

	myApp := app.New()
	myWindow := myApp.NewWindow("Lembrete")

	gui.CreateDecksUI(decks, db, myApp, myWindow)

	myWindow.Resize(fyne.NewSize(600, 400))
	myWindow.ShowAndRun()
}
