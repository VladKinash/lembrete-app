package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func showError(err error, window fyne.Window) {
	if err != nil {
		dialog.ShowError(err, window)
	}
}

func ShowWorkInProgress(window fyne.Window) {
	content := container.NewVBox(
		widget.NewLabel("Work in progress"),
	)
	window.SetContent(content)
}


