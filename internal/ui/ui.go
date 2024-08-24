package ui

import (
	"factorial/internal/ui/components"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func SetupUI(w fyne.Window) {
	initDatabase()
	initLogging()

	w.Resize(fyne.NewSize(500, 400))

	titleContainer := components.CreateTitleContainer(w)
	description := components.CreateDescription()
	inputContainer := components.CreateInputContainer()
	resultContainer := components.CreateResultContainer(w)
	progressBarContainer := components.CreateProgressBarContainer()

	calculateButton := components.CreateCalculateButton(w, resultContainer, progressBarContainer, inputContainer)

	content := container.NewVBox(
		titleContainer,
		description,
		inputContainer,
		components.CreateCalculateRangeCheckbox(inputContainer),
		calculateButton,
		resultContainer,
		container.NewCenter(progressBarContainer),
	)

	resultLabel := resultContainer.Objects[0].(*fyne.Container).Objects[0].(*widget.Label)

	w.SetContent(content)
	w.SetMainMenu(components.CreateMenu(w, resultLabel, progressBarContainer))
	w.SetOnClosed(func() {
		closeApp()
	})
}
