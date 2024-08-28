package components

import (
	"factorial/internal/database"
	"factorial/internal/logic"
	"factorial/internal/utils"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"strconv"
	"sync"
)

func CreateCalculateButton(w fyne.Window, resultContainer *fyne.Container, progressBarContainer *fyne.Container, inputContainer *fyne.Container, calculateRangeCheckbox *widget.Check, updateResult func(string)) *widget.Button {
	resultLabel := resultContainer.Objects[0].(*fyne.Container).Objects[0].(*widget.Label)
	copyButton := resultContainer.Objects[1].(*fyne.Container).Objects[0].(*widget.Button)
	progressBar := progressBarContainer.Objects[1].(*widget.ProgressBarInfinite)

	entry := inputContainer.Objects[1].(*CustomEntry)
	lowerRangeEntry := inputContainer.Objects[3].(*CustomEntry)
	upperRangeEntry := inputContainer.Objects[4].(*CustomEntry)

	return widget.NewButton("Calculate Factorial", func() {
		resultLabel.SetText("")
		copyButton.Hide()
		progressBar.Show()
		progressBar.Start()

		go func() {
			defer func() {
				progressBar.Stop()
				progressBar.Hide()
			}()

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
					updateResult(result)
				} else {
					result := logic.Factorial(number)
					utils.LogResult(number, false)
					updateResult(result.String())

					err = database.SaveResult(number, result.String())
					if err != nil {
						resultLabel.SetText("Error saving result to database")
						utils.LogError("Error saving result to database: " + err.Error())
					}
				}
			}
		}()
	})
}
