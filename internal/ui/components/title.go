package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	url2 "net/url"
)

func CreateTitleContainer(w fyne.Window) *fyne.Container {
	titleText := canvas.NewText("Factorial Calculator", nil)
	titleText.TextSize = 24
	titleText.Alignment = fyne.TextAlignCenter

	emojiLink := widget.NewButton("", func() {
		url, _ := url2.Parse("https://eduardoprofe666.github.io/")
		fyne.CurrentApp().OpenURL(url)
	})
	emojiLink.SetText("🎩")

	titleContainer := container.NewHBox(titleText, emojiLink)
	return container.NewCenter(titleContainer)
}
