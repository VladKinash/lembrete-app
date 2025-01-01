package gui

import (
	"database/sql"
	"fmt"
	"time"

	"Lembrete/models"
	repo "Lembrete/db"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func ShowDeckAndCards(deck models.Deck, db *sql.DB, app fyne.App, window fyne.Window) error {
	return updateDeckAndCardsContent(deck, db, window, nil)
}

func UpdateDeckAndCards(deck models.Deck, db *sql.DB, window fyne.Window, deckColumn *fyne.Container) error {
	return updateDeckAndCardsContent(deck, db, window, deckColumn)
}

func updateDeckAndCardsContent(deck models.Deck, db *sql.DB, window fyne.Window, deckColumn *fyne.Container) error {
	cards, err := repo.FetchAllCards(db, deck.ID)
	if err != nil {
		return fmt.Errorf("error fetching cards: %v", err)
	}

	if deckColumn == nil {
		deckColumn = container.NewVBox()
		for _, d := range []models.Deck{deck} {
			d := d
			deckButton := widget.NewButton(d.Name, func() {
				if err := UpdateDeckAndCards(d, db, window, deckColumn); err != nil {
					showError(err, window)
				}
			})
			deckColumn.Add(deckButton)
		}
	}
	deckScroll := container.NewVScroll(deckColumn)

	cardColumn := container.NewVBox()
	for _, card := range cards {
		cardButton := widget.NewButton(card.Front, func() {
			ShowCardDetails(card, window)
		})
		cardColumn.Add(cardButton)
	}
	cardScroll := container.NewVScroll(cardColumn)

	detailsContainer := widget.NewLabel("Select a card to see details")

	content := container.New(layout.NewGridLayout(3),
		deckScroll, cardScroll, container.NewVBox(detailsContainer))

	window.SetContent(content)
	return nil
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
		window.Content().(*fyne.Container).Objects[0],
		window.Content().(*fyne.Container).Objects[1],
		cardDetails,
	))
}

func showAddCardDialog(deck models.Deck, db *sql.DB, window fyne.Window) {
	frontEntry := widget.NewEntry()
	backEntry := widget.NewEntry()

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Front", Widget: frontEntry},
			{Text: "Back", Widget: backEntry},
		},
		OnSubmit: func() {
			card := models.NewFlashcard(
				frontEntry.Text,
				backEntry.Text,
				2.5,
				0,
				0,
				time.Now(),
				fmt.Sprint(deck.ID),
				0,
			)

			err := repo.InsertCard(db, card)
			if err != nil {
				dialog.ShowError(err, window)
				return
			}

			dialog.ShowInformation("Success", "Card added successfully!", window)
			frontEntry.SetText("")
			backEntry.SetText("")
		},
		OnCancel: func() {
		},
		SubmitText: "Add Card",
	}

	dialog.ShowCustomConfirm("Add New Card", "Add", "Cancel", form, func(bool) {}, window)
}