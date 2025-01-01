package gui

import (
	"database/sql"
	"fmt"
	"image/color"

	"Lembrete/models"
	repo "Lembrete/db"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func CreateDecksUI(decks []models.Deck, db *sql.DB, app fyne.App, window fyne.Window) {
	deckList := container.NewVBox()

	for _, deck := range decks {
		deck := deck
		dropdownMenu := widget.NewSelect([]string{"Show All", "Start Review"}, func(selected string) {
			if selected == "Show All" {
				if err := ShowDeckAndCards(deck, db, app, window); err != nil {
					showError(err, window)
				}
			} else if selected == "Start Review" {
				if err := StartReview(deck, db, window); err != nil {
					showError(err, window)
				}
			}
		})

		deckButton := widget.NewButton(deck.Name, func() {
			showDeckOverview(deck, db, window)
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

func showDeckOverview(deck models.Deck, db *sql.DB, window fyne.Window) {
	newCardsCount, dueCardsCount, err := repo.CountCards(db, deck.ID)
	if err != nil {
		showError(err, window)
		return
	}

	deckNameText := canvas.NewText(deck.Name, color.White)
	deckNameText.TextStyle = fyne.TextStyle{Bold: true}
	deckNameText.Alignment = fyne.TextAlignCenter
	deckNameText.TextSize = 20

	studyButton := widget.NewButton("Study", func() {
		if err := StartReview(deck, db, window); err != nil {
			showError(err, window)
		}
	})

	backButton := widget.NewButton("Back", func() {
		decks, err := repo.FetchAllDecks(db)
		if err != nil {
			showError(err, window)
			return
		}
		CreateDecksUI(decks, db, fyne.CurrentApp(), window)
	})

	addCardButton := widget.NewButton("Add Card", func() {
		showAddCardDialog(deck, db, window)
	})

	newCardsText := canvas.NewText(fmt.Sprintf("New: %d", newCardsCount), color.White)
	newCardsText.TextStyle = fyne.TextStyle{Bold: true}
	newCardsText.Alignment = fyne.TextAlignCenter

	dueCardsText := canvas.NewText(fmt.Sprintf("Due: %d", dueCardsCount), color.White)
	dueCardsText.TextStyle = fyne.TextStyle{Bold: true}
	dueCardsText.Alignment = fyne.TextAlignCenter

	statsContainer := container.New(layout.NewGridLayout(2),
		newCardsText, dueCardsText,
	)

	content := container.New(layout.NewVBoxLayout(),
		deckNameText,
		statsContainer,
		layout.NewSpacer(),
		studyButton,
		addCardButton,
		layout.NewSpacer(),
		backButton,
	)

	paddedContent := container.NewPadded(content)

	window.SetContent(paddedContent)
	window.Resize(fyne.NewSize(300, 300))
	window.CenterOnScreen()
	window.Show()
}