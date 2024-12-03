package gui
import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"Lembrete/models" // Update with the actual path to your models package
)

func CreateDecksUI(decks []models.Deck) fyne.CanvasObject {
	deckList := container.NewVBox()

	for _, deck := range decks {
		deckInfo := widget.NewLabel(
			fmt.Sprintf("ID: %d, Name: %s, Max New: %d, Max Reviews: %d",
				deck.ID, deck.Name, deck.MaxNewCards, deck.MaxReviewsDaily),
		)
		deckList.Add(deckInfo)
	}

	return deckList
}