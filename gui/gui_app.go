package gui

import (
	algorithm "Lembrete/algorithm"
	repo "Lembrete/db"
	"Lembrete/models"
	"database/sql"
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type reviewUIState struct {
	db               *sql.DB
	reviewQueue      *models.ReviewQueue
	cardFrontLabel   *canvas.Text
	cardBackLabel    *widget.Label
	ratingContainer  *fyne.Container
	showAnswerButton *widget.Button
	reviewWindow     fyne.Window
	currentCard      *models.Flashcard
}

func showError(err error, window fyne.Window) {
	if err != nil {
		dialog.ShowError(err, window)
	}
}

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
		window.Content().(*fyne.Container).Objects[0], // Keep the decks column
		window.Content().(*fyne.Container).Objects[1], // Keep the cards column
		cardDetails, // Replace card details
	))
}

func ShowWorkInProgress(window fyne.Window) {
	content := container.NewVBox(
		widget.NewLabel("Work in progress"),
	)
	window.SetContent(content)
}

func StartReview(deck models.Deck, db *sql.DB, window fyne.Window) error {
	newCards, err := repo.FetchNewCards(db, deck.ID, deck.MaxNewCards)
	if err != nil {
		return fmt.Errorf("error fetching new cards: %v", err)
	}

	dueCards, err := repo.FetchDueCards(db, deck.ID, deck.MaxReviewsDaily)
	if err != nil {
		return fmt.Errorf("error fetching due cards: %v", err)
	}

	reviewQueue := models.NewReviewQueue(newCards, dueCards)

	ShowReviewWindow(deck, db, reviewQueue, window)
	return nil
}

func ShowReviewWindow(deck models.Deck, db *sql.DB, reviewQueue *models.ReviewQueue, window fyne.Window) {
	reviewWindow := fyne.CurrentApp().NewWindow("Review: " + deck.Name)
	reviewWindow.Resize(fyne.NewSize(500, 400))
	reviewWindow.SetFixedSize(true)

	cardFrontLabel, cardBackLabel, showAnswerButton, ratingContainer := setupReviewUI()

	state := &reviewUIState{
		db:               db,
		reviewQueue:      reviewQueue,
		cardFrontLabel:   cardFrontLabel,
		cardBackLabel:    cardBackLabel,
		ratingContainer:  ratingContainer,
		showAnswerButton: showAnswerButton,
		reviewWindow:     reviewWindow,
	}

	nextCard := createNextCardFunc(state)
	nextCard()

	content := container.NewVBox(
		cardFrontLabel,
		showAnswerButton,
		cardBackLabel,
		layout.NewSpacer(),
		ratingContainer,
	)

	reviewWindow.SetContent(content)
	reviewWindow.Show()
}

func setupReviewUI() (*canvas.Text, *widget.Label, *widget.Button, *fyne.Container) {
	cardFrontLabel := canvas.NewText("Front of card", nil)
	cardFrontLabel.Alignment = fyne.TextAlignCenter
	cardFrontLabel.TextStyle = fyne.TextStyle{Bold: true}

	cardBackLabel := widget.NewLabel("Back of card")
	cardBackLabel.Alignment = fyne.TextAlignCenter
	cardBackLabel.Hide()

	showAnswerButton := widget.NewButton("Show Answer", func() {
		cardBackLabel.Show()
	})

	ratingContainer := container.NewHBox()

	return cardFrontLabel, cardBackLabel, showAnswerButton, ratingContainer
}

func createNextCardFunc(state *reviewUIState) func() {
	return func() {
		currentCard := state.reviewQueue.Next()
		if currentCard == nil {
			ShowReviewCompleteMessage(state.reviewWindow)
			return
		}
		state.currentCard = currentCard // Store current card in state
		updateCardDisplay(state)
	}
}

func updateCardDisplay(state *reviewUIState) {
	state.cardFrontLabel.Text = state.currentCard.Front
	state.cardFrontLabel.Refresh()

	state.cardBackLabel.SetText(state.currentCard.Back)
	state.cardBackLabel.Hide()

	againButton := widget.NewButton("Again", func() {
		if err := UpdateCardAndGetNext(state.currentCard, 0, state); err != nil {
			showError(err, state.reviewWindow)
		}
		state.currentCard = state.reviewQueue.Next()
		if state.currentCard == nil {
			ShowReviewCompleteMessage(state.reviewWindow)
			return
		}
		updateCardDisplay(state)
	})

	hardButton := widget.NewButton("Hard", func() {
		if err := UpdateCardAndGetNext(state.currentCard, 1, state); err != nil {
			showError(err, state.reviewWindow)
		}
		state.currentCard = state.reviewQueue.Next()
		if state.currentCard == nil {
			ShowReviewCompleteMessage(state.reviewWindow)
			return
		}
		updateCardDisplay(state)
	})

	goodButton := widget.NewButton("Good", func() {
		if err := UpdateCardAndGetNext(state.currentCard, 2, state); err != nil {
			showError(err, state.reviewWindow)
		}
		state.currentCard = state.reviewQueue.Next()
		if state.currentCard == nil {
			ShowReviewCompleteMessage(state.reviewWindow)
			return
		}
		updateCardDisplay(state)
	})

	easyButton := widget.NewButton("Easy", func() {
		if err := UpdateCardAndGetNext(state.currentCard, 3, state); err != nil {
			showError(err, state.reviewWindow)
		}
		state.currentCard = state.reviewQueue.Next()
		if state.currentCard == nil {
			ShowReviewCompleteMessage(state.reviewWindow)
			return
		}
		updateCardDisplay(state)
	})

	state.ratingContainer.Objects = []fyne.CanvasObject{againButton, hardButton, goodButton, easyButton}
	state.ratingContainer.Refresh()
	state.showAnswerButton.Show()
}

func UpdateCardAndGetNext(card *models.Flashcard, rating int, state *reviewUIState) error {
	updatedCard, err := algorithm.SM2Algorithm(card, float32(rating))
	if err != nil {
		return fmt.Errorf("error applying SM2 algorithm: %v", err)
	}

	card.EaseFactor = updatedCard.EaseFactor
	card.Repetitions = updatedCard.Repetitions
	card.Interval = updatedCard.Interval
	card.NextReview = time.Now().AddDate(0, 0, int(card.Interval))

	err = repo.UpdateCardRecord(state.db, card)
	if err != nil {
		return fmt.Errorf("error updating card: %v", err)
	}
	return nil
}

func ShowReviewCompleteMessage(window fyne.Window) {
	completeWindow := fyne.CurrentApp().NewWindow("Review Complete")
	completeWindow.Resize(fyne.NewSize(300, 200))

	messageLabel := widget.NewLabel("Review complete for now!")
	messageLabel.Alignment = fyne.TextAlignCenter

	okButton := widget.NewButton("OK", func() {
		completeWindow.Close()
		window.Close()
	})
	okButton.SetText("OK")

	content := container.New(layout.NewVBoxLayout(),
		layout.NewSpacer(),
		messageLabel,
		layout.NewSpacer(),
		okButton,
		layout.NewSpacer(),
	)

	completeWindow.SetContent(content)
	completeWindow.Show()
}
