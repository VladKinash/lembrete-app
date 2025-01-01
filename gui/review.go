package gui

import (
	"database/sql"
	"fmt"
	"time"

	algorithm "Lembrete/algorithm"
	"Lembrete/models"
	repo "Lembrete/db"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
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
		state.currentCard = currentCard
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