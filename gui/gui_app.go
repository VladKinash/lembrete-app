package gui

import (
	"database/sql"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"Lembrete/models"
	repo "Lembrete/db"
)

func CreateDecksUI(decks []models.Deck, db *sql.DB, app fyne.App, window fyne.Window) {
	deckList := container.NewVBox()

	for _, deck := range decks {
		deck := deck
		dropdownMenu := widget.NewSelect([]string{"Show All"}, func(selected string) {
			if selected == "Show All" {
				ShowDeckAndCards(deck, db, app, window)
			}
		})

		deckButton := widget.NewButton(deck.Name, func() {
			ShowWorkInProgress(window)
		})

		row := container.NewHBox(deckButton, dropdownMenu)
		deckList.Add(row)
	}

	centeredLayout := container.New(layout.NewVBoxLayout(),
		layout.NewSpacer(),
		container.New(layout.NewHBoxLayout(), layout.NewSpacer(), deckList, layout.NewSpacer()),
		layout.NewSpacer(),
	)
	window.SetContent(centeredLayout)
}



func ShowDeckAndCards(deck models.Deck, db *sql.DB, app fyne.App, window fyne.Window) {
	cards, err := repo.FetchAllCards(db, deck.ID)
	if err != nil {
		fmt.Println("Error fetching cards:", err)
		return
	}

	deckColumn := container.NewVBox()
	for _, d := range []models.Deck{deck} { 
		d := d
		deckButton := widget.NewButton(d.Name, func() {
			UpdateDeckAndCards(d, db, window, deckColumn)
		})
		deckColumn.Add(deckButton)
	}
	deckScroll := container.NewVScroll(deckColumn)

	cardColumn := container.NewVBox()
	for _, card := range cards {
		card := card
		cardButton := widget.NewButton(card.Front, func() {
			ShowCardDetails(card, window)
		})
		cardColumn.Add(cardButton)
	}
	cardScroll := container.NewVScroll(cardColumn)

	cardDetails := widget.NewLabel("Select a card to see details")
	detailsContainer := container.NewVBox(cardDetails)

	content := container.New(layout.NewGridLayout(3),
		deckScroll, cardScroll, detailsContainer)

	window.SetContent(content)
}



func UpdateDeckAndCards(deck models.Deck, db *sql.DB, window fyne.Window, deckColumn *fyne.Container) {
	cards, err := repo.FetchAllCards(db, deck.ID)
	if err != nil {
		fmt.Println("Error fetching cards:", err)
		return
	}

	cardColumn := container.NewVBox()
	for _, card := range cards {
		card := card
		cardButton := widget.NewButton(card.Front, func() {
			ShowCardDetails(card, window)
		})
		cardColumn.Add(cardButton)
	}

	cardScroll := container.NewVScroll(cardColumn)
	detailsContainer := widget.NewLabel("Select a card to see details")

	newContent := container.New(layout.NewGridLayout(3),
		deckColumn, cardScroll, container.NewVBox(detailsContainer))

	window.SetContent(newContent)
}

func ShowCardDetails(card models.Flashcard, window fyne.Window) {
	cardDetails := container.NewVBox(
		widget.NewLabel(fmt.Sprintf("Card ID: %d", card.ID)),
		widget.NewLabel(fmt.Sprintf("Front: %s", card.Front)),
		widget.NewLabel(fmt.Sprintf("Back: %s", card.Back)),
		widget.NewLabel(fmt.Sprintf("Ease Factor: %.2f", card.EaseFactor)),
		widget.NewLabel(fmt.Sprintf("Repetitions: %d", card.Repetitions)),
		widget.NewLabel(fmt.Sprintf("Interval: %.2f", card.Interval)),
		widget.NewLabel(fmt.Sprintf("Next Review: %s", card.NextReview.Format("2006-01-02"))),
	)

	window.SetContent(container.New(layout.NewGridLayout(3),
		window.Content().(*fyne.Container).Objects[0], // Keep the decks column
		window.Content().(*fyne.Container).Objects[1], // Keep the cards column
		cardDetails,                                   // Replace card details
	))
}


func ShowWorkInProgress(window fyne.Window) {
	content := container.NewVBox(
		widget.NewLabel("Work in progress"),
	)
	window.SetContent(content)
}
