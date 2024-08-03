package main

import (
	"github.com/rivo/tview"
)

func newHelpView(done func()) *tview.Modal {
	text := "Ctrl-T - select file"
	return tview.NewModal().
		AddButtons([]string{"Close"}).
		SetText(text).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			done()
		})
}
