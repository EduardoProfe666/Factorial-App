package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	url2 "net/url"
	"strconv"
	"time"
)

func ShowPreferencesDialog(w fyne.Window) {
	darkModeCheckbox := widget.NewCheck("Dark Mode", nil)
	darkModeCheckbox.OnChanged = func(checked bool) {
		if checked {
			fyne.CurrentApp().Settings().SetTheme(theme.DarkTheme())
			darkModeCheckbox.Text = "Dark Mode"
		} else {
			fyne.CurrentApp().Settings().SetTheme(theme.LightTheme())
			darkModeCheckbox.Text = "Light Mode"
		}
		darkModeCheckbox.Refresh()
	}
	fyne.CurrentApp().Settings().SetTheme(theme.DarkTheme())
	darkModeCheckbox.Checked = true

	preferencesContent := container.NewVBox(
		widget.NewLabel("Main Preferences:"),
		darkModeCheckbox,
	)

	footer := container.NewBorder(nil, nil, nil, nil,
		container.NewVBox(
			widget.NewSeparator(),
			widget.NewLabelWithStyle("©️ "+strconv.Itoa(time.Now().Year())+" EduardoProfe666 🎩", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		),
	)

	emojiLink := widget.NewButton("", func() {
		url, _ := url2.Parse("https://eduardoprofe666.github.io/")
		fyne.CurrentApp().OpenURL(url)
	})
	emojiLink.SetText("🎩")

	preferencesDialog := dialog.NewCustom("Preferences", "Close", container.NewVBox(preferencesContent, footer), w)
	preferencesDialog.Show()
}
