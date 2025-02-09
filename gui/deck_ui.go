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
	"fyne.io/fyne/v2/dialog"
	
)

func CreateDecksUI(decks []models.Deck, db *sql.DB, app fyne.App, window fyne.Window) {
	deckList := container.NewVBox()

	for _, deck := range decks {
		deckList.Add(createDeckRow(deck, db, app, window))
	}

	createDeckButton := widget.NewButton("Create Deck", func() {
		showCreateDeckDialog(db, app, window)
	})
	deckList.Add(container.NewHBox(createDeckButton, layout.NewSpacer()))

	window.SetContent(container.NewCenter(deckList))
}

func createDeckRow(deck models.Deck, db *sql.DB, app fyne.App, window fyne.Window) *fyne.Container {
	deckButton := widget.NewButton(deck.Name, func() {
		showDeckOverview(deck, db, window)
	})

	actions := []string{"Show All", "Start Review", "Delete Deck", "Edit Deck"}
	actionsDropdown := widget.NewSelect(actions, func(selected string) {
		handleDeckAction(selected, deck, db, app, window)
	})

	return container.NewHBox(deckButton, actionsDropdown)
}

func handleDeckAction(action string, deck models.Deck, db *sql.DB, app fyne.App, window fyne.Window) {
	switch action {
	case "Edit Deck":
        showEditDeckDialog(deck, db, app, window)
	case "Show All":
		showAllCards(deck, db, app)
	case "Start Review":
		startDeckReview(deck, db, window)
	case "Delete Deck":
		confirmAndDeleteDeck(deck, db, app, window)
	}
}

func showAllCards(deck models.Deck, db *sql.DB, app fyne.App) {
	newWindow := app.NewWindow("All Cards: " + deck.Name)
	newWindow.Resize(fyne.NewSize(800, 600))
	if err := ShowDeckAndCards(deck, db, app, newWindow); err != nil {
		showError(err, newWindow)
	}
	newWindow.Show()
}

func startDeckReview(deck models.Deck, db *sql.DB, window fyne.Window) {
	if err := StartReview(deck, db, window); err != nil {
		showError(err, window)
	}
}

func confirmAndDeleteDeck(deck models.Deck, db *sql.DB, app fyne.App, window fyne.Window) {
	dialog.ShowConfirm("Delete Deck", "Are you sure you want to delete this deck?", func(confirmed bool) {
		if confirmed {
			deleteDeckAndRefreshUI(deck, db, app, window)
		}
	}, window)
}

func deleteDeckAndRefreshUI(deck models.Deck, db *sql.DB, app fyne.App, window fyne.Window) {
	if err := repo.DeleteDeck(db, deck); err != nil {
		showError(err, window)
		return
	}
	decks, err := repo.FetchAllDecks(db)
	if err != nil {
		showError(err, window)
		return
	}
	CreateDecksUI(decks, db, app, window)
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

    // Stats Container
    newCardsText := canvas.NewText(fmt.Sprintf("New: %d", newCardsCount), color.White)
    newCardsText.TextStyle = fyne.TextStyle{Bold: true}
    newCardsText.Alignment = fyne.TextAlignCenter

    dueCardsText := canvas.NewText(fmt.Sprintf("Due: %d", dueCardsCount), color.White)
    dueCardsText.TextStyle = fyne.TextStyle{Bold: true}
    dueCardsText.Alignment = fyne.TextAlignCenter

    statsContainer := container.New(layout.NewGridLayout(2),
        newCardsText, dueCardsText,
    )

    // Create buttons with fixed size
    studyButton := widget.NewButton("Study", func() {
        if err := StartReview(deck, db, window); err != nil {
            showError(err, window)
        }
    })

    addCardButton := widget.NewButton("Add Card", func() {
        showAddCardDialog(deck, db, window)
    })

    editDeckButton := widget.NewButton("Edit", func() {
        showEditDeckDialog(deck, db, fyne.CurrentApp(), window)
    })

    backButton := widget.NewButton("Back", func() {
        decks, err := repo.FetchAllDecks(db)
        if err != nil {
            showError(err, window)
            return
        }
        CreateDecksUI(decks, db, fyne.CurrentApp(), window)
    })

    buttonContainer := container.NewHBox(
        layout.NewSpacer(),
        container.NewHBox(
            studyButton,
            addCardButton,
            editDeckButton,
            backButton,
        ),
        layout.NewSpacer(),
    )

    content := container.NewBorder(
        nil, //top
        buttonContainer, //bottom
        nil, //left
        nil, // right
        container.NewVBox(
            deckNameText,
            statsContainer,
        ),
    )

    window.SetContent(content)
    window.Resize(fyne.NewSize(800, 600))
    window.CenterOnScreen()
    window.Show()
}


func showCreateDeckDialog(db *sql.DB, app fyne.App, window fyne.Window) {
    deckNameEntry := widget.NewEntry()
    deckNameEntry.SetPlaceHolder("Enter deck name")
    maxNewCardsEntry := widget.NewEntry()
    maxNewCardsEntry.SetPlaceHolder("Optional")
    maxReviewsEntry := widget.NewEntry()
    maxReviewsEntry.SetPlaceHolder("Optional")

    formContainer := container.NewVBox(
        widget.NewLabel("Deck Name"),
        deckNameEntry,
        widget.NewLabel("Max New Cards (Optional)"),
        maxNewCardsEntry,
        widget.NewLabel("Max Reviews (Optional)"),
        maxReviewsEntry,
    )

    var createDeckDialog dialog.Dialog

    onSubmit := func() {
        deckName := deckNameEntry.Text
        if deckName == "" {
            dialog.ShowError(fmt.Errorf("deck name cannot be empty"), window)
            return
        }

        var maxNewCards int32 = 0
        if maxNewCardsEntry.Text != "" {
            var newCards int
            _, err := fmt.Sscan(maxNewCardsEntry.Text, &newCards)
            if err != nil {
                dialog.ShowError(fmt.Errorf("invalid Max New Cards value"), window)
                return
            }
            maxNewCards = int32(newCards)
        }

        var maxReviews int32 = 0
        if maxReviewsEntry.Text != "" {
            var reviews int
            _, err := fmt.Sscan(maxReviewsEntry.Text, &reviews)
            if err != nil {
                dialog.ShowError(fmt.Errorf("invalid Max Reviews value"), window)
                return
            }
            maxReviews = int32(reviews)
        }

        newDeck := models.Deck{
            Name:           deckName,
            MaxNewCards:    maxNewCards,
            MaxReviewsDaily: maxReviews,
        }

        err := repo.InsertDeck(db, newDeck)
        if err != nil {
            dialog.ShowError(err, window)
            return
        }

        dialog.ShowInformation("Success", "Deck created successfully!", window)
        decks, err := repo.FetchAllDecks(db)
        if err != nil {
            dialog.ShowError(fmt.Errorf("error updating deck list: %v", err), window)
            return
        }

        CreateDecksUI(decks, db, app, window)
        createDeckDialog.Hide()
    }

    createButton := widget.NewButton("Create", onSubmit)
    buttonContainer := container.NewHBox(layout.NewSpacer(), createButton)

    content := container.NewVBox(
        formContainer,
        buttonContainer,
    )

    createDeckDialog = dialog.NewCustom("Create New Deck", "Cancel", content, window)
    createDeckDialog.Show()

    deckNameEntry.OnSubmitted = func(_ string) { onSubmit() }
    maxNewCardsEntry.OnSubmitted = func(_ string) { onSubmit() }
    maxReviewsEntry.OnSubmitted = func(_ string) { onSubmit() }
}


func showEditDeckDialog(deck models.Deck, db *sql.DB, app fyne.App, window fyne.Window) {
    deckNameEntry := widget.NewEntry()
    deckNameEntry.SetText(deck.Name)
    maxNewCardsEntry := widget.NewEntry()
    maxNewCardsEntry.SetText(fmt.Sprintf("%d", deck.MaxNewCards))
    maxReviewsEntry := widget.NewEntry()
    maxReviewsEntry.SetText(fmt.Sprintf("%d", deck.MaxReviewsDaily))

    formContainer := container.NewVBox(
        widget.NewLabel("Deck Name"),
        deckNameEntry,
        widget.NewLabel("Max New Cards (Optional)"),
        maxNewCardsEntry,
        widget.NewLabel("Max Reviews (Optional)"),
        maxReviewsEntry,
    )

    var editDeckDialog dialog.Dialog

    onSubmit := func() {
        deckName := deckNameEntry.Text
        if deckName == "" {
            dialog.ShowError(fmt.Errorf("deck name cannot be empty"), window)
            return
        }

        var maxNewCards int32 = 0
        if maxNewCardsEntry.Text != "" {
            var newCards int
            _, err := fmt.Sscan(maxNewCardsEntry.Text, &newCards)
            if err != nil {
                dialog.ShowError(fmt.Errorf("invalid Max New Cards value"), window)
                return
            }
            maxNewCards = int32(newCards)
        }

        var maxReviews int32 = 0
        if maxReviewsEntry.Text != "" {
            var reviews int
            _, err := fmt.Sscan(maxReviewsEntry.Text, &reviews)
            if err != nil {
                dialog.ShowError(fmt.Errorf("invalid Max Reviews value"), window)
                return
            }
            maxReviews = int32(reviews)
        }

        updatedDeck := models.Deck{
            ID:              deck.ID,
            Name:            deckName,
            MaxNewCards:     maxNewCards,
            MaxReviewsDaily: maxReviews,
        }

        err := repo.UpdateDeckRecord(db, updatedDeck)
        if err != nil {
            dialog.ShowError(err, window)
            return
        }

        dialog.ShowInformation("Success", "Deck updated successfully!", window)
        decks, err := repo.FetchAllDecks(db)
        if err != nil {
            dialog.ShowError(fmt.Errorf("error updating deck list: %v", err), window)
            return
        }

        CreateDecksUI(decks, db, app, window)
        editDeckDialog.Hide()
    }

    updateButton := widget.NewButton("Update", onSubmit)
    buttonContainer := container.NewHBox(layout.NewSpacer(), updateButton)

    content := container.NewVBox(
        formContainer,
        buttonContainer,
    )

    editDeckDialog = dialog.NewCustom("Edit Deck", "Cancel", content, window)
    editDeckDialog.Show()

    deckNameEntry.OnSubmitted = func(_ string) { onSubmit() }
    maxNewCardsEntry.OnSubmitted = func(_ string) { onSubmit() }
    maxReviewsEntry.OnSubmitted = func(_ string) { onSubmit() }
}