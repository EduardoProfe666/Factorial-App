package components

import (
	utils2 "factorial/internal/ui/utils"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"time"
)

func CreateResultContainer(w fyne.Window) (*fyne.Container, func(string)) {
	resultLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	resultLabel.Wrapping = fyne.TextWrapOff

	// Variable para almacenar el resultado completo
	var fullResult string

	copyButton := widget.NewButtonWithIcon("", theme.ContentCopyIcon(), nil)
	copyButton.OnTapped = func() {
		// Usar el resultado completo almacenado
		w.Clipboard().SetContent(fullResult)
		copyButton.SetIcon(theme.ContentPasteIcon())
		time.AfterFunc(1*time.Second, func() {
			copyButton.SetIcon(theme.ContentCopyIcon())
		})
	}
	copyButton.Hide()

	// Función para actualizar el resultado completo y la etiqueta
	updateResult := func(result string) {
		fullResult = result
		resultLabel.SetText("Result: " + utils2.TruncateString(result, 100))
		copyButton.Show()
	}

	return container.NewHBox(
		container.NewCenter(resultLabel),
		container.NewCenter(copyButton),
	), updateResult
}
