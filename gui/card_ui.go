package gui

import (
	"database/sql"
	"fmt"
	"time"
	"strconv"

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

	detailsContainer := container.NewVBox(widget.NewLabel("Select a card to see details"))

	cardButtonsColumn := container.NewVBox()
	for _, c := range cards {
		card := c
		cardButton := widget.NewButton(card.Front, func() {
			detailsContainer.Objects = []fyne.CanvasObject{
				widget.NewLabel(fmt.Sprintf("Front: %s", card.Front)),
				widget.NewLabel(fmt.Sprintf("Back: %s", card.Back)),
				widget.NewLabel(fmt.Sprintf("Ease Factor: %.2f", card.EaseFactor)),
				widget.NewLabel(fmt.Sprintf("Repetitions: %d", card.Repetitions)),
				widget.NewLabel(fmt.Sprintf("Interval: %.2f", card.Interval)),
				widget.NewLabel(fmt.Sprintf("Next Review: %s", card.NextReview.Format("2006-01-02"))),
				widget.NewButton("Edit", func() {
					showEditCardDialog(card, db, window)
				}),
			}
			detailsContainer.Refresh()
		})
		cardButtonsColumn.Add(cardButton)
	}
	cardScroll := container.NewVScroll(cardButtonsColumn)

	content := container.New(layout.NewGridLayout(3),
		deckScroll,
		cardScroll,
		detailsContainer,
	)

	window.SetContent(content)
	return nil
}



func ShowCardDetails(deck models.Deck, card models.Flashcard, db *sql.DB, window fyne.Window) {
    backButton := widget.NewButton("Back", func() {
        if err := updateDeckAndCardsContent(deck, db, window, nil); err != nil {
            showError(err, window)
        }
    })

    deleteButton := widget.NewButton("Delete Card", func() {
        dialog.ShowConfirm("Delete Card",
            "Are you sure you want to delete this card?",
            func(confirmed bool) {
                if confirmed {
                    if err := repo.DeleteCard(db, card); err != nil {
                        showError(err, window)
                        return
                    }
                    if err := updateDeckAndCardsContent(deck, db, window, nil); err != nil {
                        showError(err, window)
                    }
                }
            }, 
            window)
    })

    topButtons := container.NewHBox(backButton, deleteButton)

    cardDetails := container.NewVBox(
        topButtons, 
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
	frontEntry.SetPlaceHolder("Enter front text")
	backEntry := widget.NewEntry()
	backEntry.SetPlaceHolder("Enter back text")

	formContainer := container.NewVBox(
		widget.NewLabel("Front"),
		frontEntry,
		widget.NewLabel("Back"),
		backEntry,
	)

	var addCardDialog dialog.Dialog

	addButton := widget.NewButton("Add Card", func() {
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
	})

	cancelButton := widget.NewButton("Cancel", func() {
		addCardDialog.Hide()
	})

	buttonContainer := container.NewHBox(addButton, cancelButton)

	content := container.NewVBox(formContainer, buttonContainer)

	addCardDialog = dialog.NewCustom("Add New Card", "Close", content, window)
	addCardDialog.Show()
}




func showEditCardDialog(card models.Flashcard, db *sql.DB, parentWindow fyne.Window) {
	editWindow := fyne.CurrentApp().NewWindow("Edit Card")
	editWindow.Resize(fyne.NewSize(600, 800))

	frontEntry := widget.NewEntry()
	frontEntry.SetText(card.Front)

	backEntry := widget.NewEntry()
	backEntry.SetText(card.Back)

	easeFactorEntry := widget.NewEntry()
	easeFactorEntry.SetText(fmt.Sprintf("%.2f", card.EaseFactor))

	repetitionsEntry := widget.NewEntry()
	repetitionsEntry.SetText(fmt.Sprintf("%d", card.Repetitions))

	intervalEntry := widget.NewEntry()
	intervalEntry.SetText(fmt.Sprintf("%.2f", card.Interval))

	nextReviewEntry := widget.NewEntry()
	nextReviewEntry.SetText(card.NextReview.Format("2006-01-02"))

	deckIDEntry := widget.NewEntry()
	deckIDEntry.SetText(card.DeckID)

	idLabel := widget.NewLabel(fmt.Sprintf("ID: %d", card.ID))

	saveButton := widget.NewButton("Save", func() {
		newFront := frontEntry.Text
		newBack := backEntry.Text

		newEaseFactor, err := strconv.ParseFloat(easeFactorEntry.Text, 32)
		if err != nil {
			showError(fmt.Errorf("invalid Ease Factor"), editWindow)
			return
		}

		newRepetitions, err := strconv.Atoi(repetitionsEntry.Text)
		if err != nil {
			showError(fmt.Errorf("invalid Repetitions"), editWindow)
			return
		}

		newInterval, err := strconv.ParseFloat(intervalEntry.Text, 32)
		if err != nil {
			showError(fmt.Errorf("invalid Interval"), editWindow)
			return
		}

		newNextReview, err := time.Parse("2006-01-02", nextReviewEntry.Text)
		if err != nil {
			showError(fmt.Errorf("invalid Next Review date"), editWindow)
			return
		}

		newDeckID, err := strconv.Atoi(deckIDEntry.Text)
		if err != nil {
			showError(fmt.Errorf("invalid Deck ID"), editWindow)
			return
		}

		updatedCard := models.Flashcard{
			Front:       newFront,
			Back:        newBack,
			EaseFactor:  float32(newEaseFactor),
			Repetitions: newRepetitions,
			Interval:    float32(newInterval),
			NextReview:  newNextReview,
			DeckID:      strconv.Itoa(newDeckID),
			ID:          card.ID,
		}

		err = repo.UpdateCardRecord(db, &updatedCard)
		if err != nil {
			showError(err, editWindow)
			return
		}

		editWindow.Close()

		deckIDInt, err := strconv.Atoi(updatedCard.DeckID)
		if err != nil {
			showError(fmt.Errorf("invalid Deck ID after update"), parentWindow)
			return
		}

		err = updateDeckAndCardsContent(models.Deck{ID: int32(deckIDInt)}, db, parentWindow, nil)
		if err != nil {
			showError(err, parentWindow)
		}
	})

	cancelButton := widget.NewButton("Cancel", func() {
		editWindow.Close()
	})

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Front", Widget: frontEntry},
			{Text: "Back", Widget: backEntry},
			{Text: "Ease Factor", Widget: easeFactorEntry},
			{Text: "Repetitions", Widget: repetitionsEntry},
			{Text: "Interval", Widget: intervalEntry},
			{Text: "Next Review (YYYY-MM-DD)", Widget: nextReviewEntry},
			{Text: "Deck ID", Widget: deckIDEntry},
			{Text: "ID", Widget: idLabel},
		},
	}

	buttons := container.NewHBox(saveButton, cancelButton)
	content := container.NewVBox(form, buttons)

	editWindow.SetContent(content)
	editWindow.Show()
}
