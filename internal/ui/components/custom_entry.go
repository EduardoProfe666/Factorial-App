package components

import (
	"fyne.io/fyne/v2/widget"
	"unicode"
)

type CustomEntry struct {
	widget.Entry
}

func NewCustomEntry() *CustomEntry {
	entry := &CustomEntry{}
	entry.ExtendBaseWidget(entry)
	return entry
}

func (e *CustomEntry) TypedRune(r rune) {
	if unicode.IsDigit(r) {
		e.Entry.TypedRune(r)
	}
}
