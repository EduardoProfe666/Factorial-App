package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func CreateInputContainer() *fyne.Container {
	numberInput := widget.NewLabel("Number Input")
	rangeInput := widget.NewLabel("Range Input")
	rangeInput.Hide()

	entry := NewCustomEntry()
	entry.SetPlaceHolder("Enter a number")

	lowerRangeEntry := NewCustomEntry()
	lowerRangeEntry.SetPlaceHolder("Enter lower limit")
	lowerRangeEntry.Hide()

	upperRangeEntry := NewCustomEntry()
	upperRangeEntry.SetPlaceHolder("Enter upper limit")
	upperRangeEntry.Hide()

	return container.NewVBox(
		numberInput,
		entry,
		rangeInput,
		lowerRangeEntry,
		upperRangeEntry,
	)
}

func CreateCalculateRangeCheckbox(inputContainer *fyne.Container) *widget.Check {
	numberInput := inputContainer.Objects[0].(*widget.Label)
	entry := inputContainer.Objects[1].(*CustomEntry)
	rangeInput := inputContainer.Objects[2].(*widget.Label)
	lowerRangeEntry := inputContainer.Objects[3].(*CustomEntry)
	upperRangeEntry := inputContainer.Objects[4].(*CustomEntry)

	return widget.NewCheck("Calculate Range", func(checked bool) {
		if checked {
			entry.Hide()
			lowerRangeEntry.Show()
			upperRangeEntry.Show()
			numberInput.Hide()
			rangeInput.Show()
		} else {
			entry.Show()
			lowerRangeEntry.Hide()
			upperRangeEntry.Hide()
			numberInput.Show()
			rangeInput.Hide()
		}
	})
}
