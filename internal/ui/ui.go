package ui

import (
	"factorial/internal/database"
	"factorial/internal/logic"
	"factorial/internal/ui/components"
	utils2 "factorial/internal/ui/utils"
	"factorial/internal/utils"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"strconv"
	"sync"
	"time"
)

func SetupUI(w fyne.Window) {
	database.InitDB()
	utils.LogInfo("Factorial App Started")

	w.Resize(fyne.NewSize(600, 400))

	title := widget.NewLabelWithStyle("Factorial Calculator", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	description := widget.NewLabelWithStyle("Enter a number or a range to calculate the factorial.", fyne.TextAlignCenter, fyne.TextStyle{})

	numberInput := widget.NewLabel("Number Input")
	rangeInput := widget.NewLabel("Range Input")
	rangeInput.Hide()

	entry := components.NewCustomEntry()
	entry.SetPlaceHolder("Enter a number")

	lowerRangeEntry := components.NewCustomEntry()
	lowerRangeEntry.SetPlaceHolder("Enter lower limit")
	lowerRangeEntry.Hide()

	upperRangeEntry := components.NewCustomEntry()
	upperRangeEntry.SetPlaceHolder("Enter upper limit")
	upperRangeEntry.Hide()

	calculateRangeCheckbox := widget.NewCheck("Calculate Range", func(checked bool) {
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

	calculateButton := widget.NewButton("Calculate Factorial", func() {
		resultLabel.SetText("")
		copyButton.Hide()
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
				resultLabel.SetText("Result: " + utils2.TruncateString(result, 100))
				copyButton.Show()
			} else {
				result := logic.Factorial(number)
				utils.LogResult(number, false)
				resultLabel.SetText("Result: " + utils2.TruncateString(result.String(), 100))
				copyButton.Show()

				err = database.SaveResult(number, result.String())
				if err != nil {
					resultLabel.SetText("Error saving result to database")
					utils.LogError("Error saving result to database: " + err.Error())
				}
			}
		}
	})

	inputContainer := container.NewVBox(
		numberInput,
		entry,
		rangeInput,
		lowerRangeEntry,
		upperRangeEntry,
	)

	resultContainer := container.NewHBox(
		container.NewCenter(resultLabel),
		container.NewCenter(copyButton),
	)

	content := container.NewVBox(
		title,
		description,
		inputContainer,
		calculateRangeCheckbox,
		calculateButton,
		resultContainer,
	)

	w.SetContent(content)

	w.SetOnClosed(func() {
		utils.LogInfo("Factorial App Exited")
	})
}
