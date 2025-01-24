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
	buttonContainer  *fyne.Container
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
	reviewWindow.Resize(fyne.NewSize(800, 600))

	cardFrontLabel, cardBackLabel, showAnswerButton, againButton, hardButton, goodButton, easyButton, ratingContainer, buttonContainer := setupReviewUI()

	state := &reviewUIState{
		db:               db,
		reviewQueue:      reviewQueue,
		cardFrontLabel:   cardFrontLabel,
		cardBackLabel:    cardBackLabel,
		ratingContainer:  ratingContainer,
		showAnswerButton: showAnswerButton,
		buttonContainer:  buttonContainer,
		reviewWindow:     reviewWindow,
	}

	rateAndProceed := func(rating int) {
		if err := UpdateCardAndGetNext(state.currentCard, rating, state); err != nil {
			showError(err, reviewWindow)
			return
		}

		state.ratingContainer.Hide()
		state.showAnswerButton.Show()
		state.cardBackLabel.Hide()
		state.buttonContainer.Refresh()

		nextCard := state.reviewQueue.Next()
		if nextCard == nil {
			ShowReviewCompleteMessage(reviewWindow)
			return
		}
		state.currentCard = nextCard
		updateCardDisplay(state)
	}

	showAnswerButton.OnTapped = func() {
		cardBackLabel.Show()
		showAnswerButton.Hide()
		ratingContainer.Show()
		buttonContainer.Refresh()
	}

	againButton.OnTapped = func() { rateAndProceed(0) }
	hardButton.OnTapped = func() { rateAndProceed(3) }
	goodButton.OnTapped = func() { rateAndProceed(4) }
	easyButton.OnTapped = func() { rateAndProceed(5) }

	state.currentCard = reviewQueue.Next()
	if state.currentCard != nil {
		updateCardDisplay(state)
	}

	content := container.NewVBox(
		layout.NewSpacer(),
		cardFrontLabel,
		layout.NewSpacer(),
		cardBackLabel,
		layout.NewSpacer(),
		buttonContainer,
	)
	reviewWindow.SetContent(content)
	reviewWindow.Show()
}

func setupReviewUI() (
	*canvas.Text,
	*widget.Label,
	*widget.Button,
	*widget.Button,
	*widget.Button,
	*widget.Button,
	*widget.Button,
	*fyne.Container,
	*fyne.Container,
) {
	cardFrontLabel := canvas.NewText("Front of card", nil)
	cardFrontLabel.Alignment = fyne.TextAlignCenter
	cardFrontLabel.TextSize = 24
	cardFrontLabel.TextStyle = fyne.TextStyle{Bold: true}

	cardBackLabel := widget.NewLabel("Back of card")
	cardBackLabel.Alignment = fyne.TextAlignCenter
	cardBackLabel.Hide()

	showAnswerButton := widget.NewButton("Show Answer", nil)

	againButton := widget.NewButton("Again", nil)
	hardButton := widget.NewButton("Hard", nil)
	goodButton := widget.NewButton("Good", nil)
	easyButton := widget.NewButton("Easy", nil)

	ratingContainer := container.NewHBox(
		layout.NewSpacer(),
		againButton,
		hardButton,
		goodButton,
		easyButton,
		layout.NewSpacer(),
	)
	ratingContainer.Hide()

	buttonContainer := container.NewVBox(
		showAnswerButton,
		ratingContainer,
	)

	return cardFrontLabel, cardBackLabel, showAnswerButton, againButton, hardButton, goodButton, easyButton, ratingContainer, buttonContainer
}
func updateCardDisplay(state *reviewUIState) {
	state.cardFrontLabel.Text = state.currentCard.Front
	state.cardFrontLabel.Refresh()

	state.cardBackLabel.SetText(state.currentCard.Back)
	state.cardBackLabel.Hide()

	state.ratingContainer.Hide()
	state.showAnswerButton.Show()
	state.buttonContainer.Refresh()
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
	completeWindow.Resize(fyne.NewSize(600, 400))

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
