package components

import (
	"factorial/internal/database"
	utils2 "factorial/internal/ui/utils"
	"factorial/internal/utils"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"strconv"
)

func CreateMenu(w fyne.Window, resultLabel *widget.Label, progressBarContainer *fyne.Container) *fyne.MainMenu {
	progressBar := progressBarContainer.Objects[1].(*widget.ProgressBarInfinite)

	return fyne.NewMainMenu(
		fyne.NewMenu("File",
			fyne.NewMenuItem("Export to CSV", func() {
				exportToCSV(resultLabel, progressBar)
			}),
			fyne.NewMenuItem("Quit", func() {
				w.Close()
			}),
		),
		fyne.NewMenu("Data",
			fyne.NewMenuItem("View Data", func() {
				viewData(w, resultLabel, progressBar)
			}),
			fyne.NewMenuItem("Delete Data", func() {
				deleteData(w, resultLabel, progressBar)
			}),
		),
		fyne.NewMenu("Preferences",
			fyne.NewMenuItem("Preferences", func() {
				ShowPreferencesDialog(w)
			}),
		),
	)
}

func exportToCSV(resultLabel *widget.Label, progressBar *widget.ProgressBarInfinite) {
	progressBar.Show()
	progressBar.Start()

	go func() {
		defer func() {
			progressBar.Stop()
			progressBar.Hide()
		}()

		err := database.ExportToCSV("results.csv")
		if err != nil {
			utils.LogError("Error exporting to CSV: " + err.Error())
			resultLabel.SetText("Error exporting to CSV")
		} else {
			resultLabel.SetText("Results exported to results.csv")
		}
	}()
}

func viewData(w fyne.Window, resultLabel *widget.Label, progressBar *widget.ProgressBarInfinite) {
	progressBar.Show()
	progressBar.Start()

	go func() {
		defer func() {
			progressBar.Stop()
			progressBar.Hide()
		}()

		results, err := database.GetResults()
		if err != nil {
			utils.LogError("Error retrieving results: " + err.Error())
			resultLabel.SetText("Error retrieving results")
			return
		}

		var items []fyne.CanvasObject
		for _, result := range results {
			label := widget.NewLabelWithStyle(strconv.Itoa(result.Number)+": "+utils2.TruncateString(result.Result, 20), fyne.TextAlignLeading, fyne.TextStyle{Monospace: true})
			label.Wrapping = fyne.TextWrapOff
			copyButton := widget.NewButtonWithIcon("", theme.ContentCopyIcon(), func() {
				w.Clipboard().SetContent(result.Result)
			})
			items = append(items, container.NewHBox(label, copyButton))
		}

		totalLabel := widget.NewLabel("Total results: " + strconv.Itoa(len(results)))
		scrollContainer := container.NewVScroll(container.NewVBox(append([]fyne.CanvasObject{totalLabel}, items...)...))
		scrollContainer.SetMinSize(fyne.NewSize(400, 300))

		dialog.NewCustom("Saved Results", "Close", scrollContainer, w).Show()
	}()
}

func deleteData(w fyne.Window, resultLabel *widget.Label, progressBar *widget.ProgressBarInfinite) {
	dialog.NewConfirm("Warning", "Are you sure you want to delete all data?", func(confirmed bool) {
		if confirmed {
			progressBar.Show()
			progressBar.Start()

			go func() {
				defer func() {
					progressBar.Stop()
					progressBar.Hide()
				}()

				err := database.ClearResults()
				if err != nil {
					utils.LogError("Error deleting data: " + err.Error())
					resultLabel.SetText("Error deleting data")
				} else {
					resultLabel.SetText("All data deleted")
				}
			}()
		}
	}, w).Show()
}
