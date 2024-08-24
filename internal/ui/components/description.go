package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func CreateDescription() *widget.Label {
	return widget.NewLabelWithStyle("Enter a number or a range to calculate the factorial.", fyne.TextAlignCenter, fyne.TextStyle{})
}
