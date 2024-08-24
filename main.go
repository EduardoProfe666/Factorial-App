package main

import (
	"factorial/internal/ui"
	"fyne.io/fyne/v2/app"
)

func main() {
	a := app.New()
	w := a.NewWindow("Factorial App")

	ui.SetupUI(w)

	w.ShowAndRun()
}
