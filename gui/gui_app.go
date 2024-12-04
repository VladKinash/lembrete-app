package gui
import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"Lembrete/models"
	"fyne.io/fyne/v2/layout"
)

func CreateDecksUI(decks []models.Deck, window fyne.Window) {
	deckList := container.NewVBox()

	for _, deck := range decks {
		deck := deck // create a local copy for closure
		deckButton := widget.NewButton(deck.Name, func() {
			ShowDeckInfo(deck, window)
		})
		deckList.Add(deckButton)
	}

	centeredLayout := container.New(layout.NewVBoxLayout(),
		layout.NewSpacer(),
		container.New(layout.NewHBoxLayout(), layout.NewSpacer(), deckList, layout.NewSpacer()),
		layout.NewSpacer(),
	)
	window.SetContent(centeredLayout)
}

func ShowDeckInfo(deck models.Deck, window fyne.Window) {
	content := container.NewVBox(
		widget.NewLabel(fmt.Sprintf("Deck ID: %d", deck.ID)),
		widget.NewLabel(fmt.Sprintf("Name: %s", deck.Name)),
		widget.NewLabel(fmt.Sprintf("Max New Cards: %d", deck.MaxNewCards)),
		widget.NewLabel(fmt.Sprintf("Max Reviews Daily: %d", deck.MaxReviewsDaily)),
		widget.NewButton("Back", func() {
			// Replace with the original deck list UI
			// This will depend on how your decks are passed or re-fetched
			CreateDecksUI([]models.Deck{deck}, window)
		}),
	)

	centeredLayout := container.New(layout.NewVBoxLayout(),
		layout.NewSpacer(),
		container.New(layout.NewHBoxLayout(), layout.NewSpacer(), content, layout.NewSpacer()),
		layout.NewSpacer(),
	)

	window.SetContent(centeredLayout)
}