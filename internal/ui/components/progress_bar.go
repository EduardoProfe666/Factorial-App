package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func CreateProgressBarContainer() *fyne.Container {
	progressBar := widget.NewProgressBarInfinite()
	progressBar.Stop()
	progressBar.Hide()

	return container.NewHBox(
		layout.NewSpacer(),
		progressBar,
		layout.NewSpacer(),
	)
}
