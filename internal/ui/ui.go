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
	calculateRangeCheckbox := components.CreateCalculateRangeCheckbox(inputContainer)
	progressBarContainer := components.CreateProgressBarContainer()
	resultContainer, updateResult := components.CreateResultContainer(w)

	calculateButton := components.CreateCalculateButton(w, resultContainer, progressBarContainer, inputContainer, calculateRangeCheckbox, updateResult)

	content := container.NewVBox(
		titleContainer,
		description,
		inputContainer,
		calculateRangeCheckbox,
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
