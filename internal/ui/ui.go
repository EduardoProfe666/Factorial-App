package ui

import (
	"factorial/internal/database"
	"factorial/internal/logic"
	"factorial/internal/utils"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"strconv"
	"sync"
)

func SetupUI(w fyne.Window) {
	database.InitDB()
	utils.LogInfo("Factorial App Started")

	// Establecer el tamaño inicial y mínimo de la ventana
	w.Resize(fyne.NewSize(600, 400))
	w.SetFixedSize(true)

	entry := widget.NewEntry()
	entry.SetPlaceHolder("Enter a number")

	lowerRangeEntry := widget.NewEntry()
	lowerRangeEntry.SetPlaceHolder("Enter lower limit")
	lowerRangeEntry.Hide() // Inicialmente oculto

	upperRangeEntry := widget.NewEntry()
	upperRangeEntry.SetPlaceHolder("Enter upper limit")
	upperRangeEntry.Hide() // Inicialmente oculto

	calculateRangeCheckbox := widget.NewCheck("Calculate Range", func(checked bool) {
		if checked {
			entry.Hide()
			lowerRangeEntry.Show()
			upperRangeEntry.Show()
		} else {
			entry.Show()
			lowerRangeEntry.Hide()
			upperRangeEntry.Hide()
		}
	})

	resultLabel := widget.NewLabel("")

	calculateButton := widget.NewButton("Calculate Factorial", func() {
		if calculateRangeCheckbox.Checked {
			lowerLimit, err1 := strconv.Atoi(lowerRangeEntry.Text)
			upperLimit, err2 := strconv.Atoi(upperRangeEntry.Text)
			if err1 != nil || err2 != nil || lowerLimit > upperLimit {
				resultLabel.SetText("Invalid range values")
				utils.LogError("Invalid range values")
				return
			}

			var wg sync.WaitGroup
			for i := lowerLimit; i <= upperLimit; i++ {
				wg.Add(1)
				go func(num int) {
					defer wg.Done()
					_, err := database.GetFactorial(num)
					if err != nil {
						result := logic.Factorial(num)
						err = database.SaveResult(num, result.String())
						if err != nil {
							utils.LogError("Error saving result to database: " + err.Error())
						}
					}
				}(i)
			}
			wg.Wait()
			resultLabel.SetText("Range factorial calculation completed.")
		} else {
			number, err := strconv.Atoi(entry.Text)
			if err != nil {
				resultLabel.SetText("Invalid input")
				utils.LogError("Invalid input")
				return
			}

			result, err := database.GetFactorial(number)
			if err == nil {
				utils.LogResult(number, true)
				resultLabel.SetText("Result: " + result)
			} else {
				result := logic.Factorial(number)
				utils.LogResult(number, false)
				resultLabel.SetText("Result: " + result.String())

				err = database.SaveResult(number, result.String())
				if err != nil {
					resultLabel.SetText("Error saving result to database")
					utils.LogError("Error saving result to database: " + err.Error())
				}
			}
		}
	})

	inputContainer := container.NewVBox(entry, lowerRangeEntry, upperRangeEntry)

	calculateRangeCheckbox.OnChanged = func(checked bool) {
		if checked {
			entry.Hide()
			lowerRangeEntry.Show()
			upperRangeEntry.Show()
		} else {
			entry.Show()
			lowerRangeEntry.Hide()
			upperRangeEntry.Hide()
		}
		inputContainer.Refresh()
	}

	content := container.NewVBox(
		inputContainer,
		calculateRangeCheckbox,
		calculateButton,
		resultLabel,
	)

	w.SetContent(content)

	w.SetOnClosed(func() {
		utils.LogInfo("Factorial App Exited")
	})
}
