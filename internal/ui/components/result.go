package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"time"
)

func CreateResultContainer(w fyne.Window) *fyne.Container {
	resultLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	resultLabel.Wrapping = fyne.TextWrapOff

	copyButton := widget.NewButtonWithIcon("", theme.ContentCopyIcon(), nil)
	copyButton.OnTapped = func() {
		result := resultLabel.Text
		if len(result) > 8 {
			result = result[8:]
		}
		w.Clipboard().SetContent(result)
		copyButton.SetIcon(theme.ContentPasteIcon())
		time.AfterFunc(1*time.Second, func() {
			copyButton.SetIcon(theme.ContentCopyIcon())
		})
	}
	copyButton.Hide()

	return container.NewHBox(
		container.NewCenter(resultLabel),
		container.NewCenter(copyButton),
	)
}
