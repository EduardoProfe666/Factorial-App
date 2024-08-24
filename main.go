package main

import (
	"factorial/internal/api"
	"factorial/internal/ui"
	"fyne.io/fyne/v2/app"
)

func main() {
	go api.StartAPI()

	a := app.New()
	w := a.NewWindow("Factorial App")

	ui.SetupUI(w)

	w.ShowAndRun()
}
