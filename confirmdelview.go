package main

import (
	"fmt"
	"github.com/rivo/tview"
)

type ConfirmDeleteAction int

const (
	ConfirmDeleteNo ConfirmDeleteAction = iota
	ConfirmDeleteYes
)

func newConfirmDeleteView(files []dirEntry, done func(a ConfirmDeleteAction)) *tview.Modal {
	var text string
	if len(files) == 1 {
		text = fmt.Sprintf("%s \n Do you want to delete it?", files[0].path)
	} else {
		text = fmt.Sprintf("Do you want to delete %v files?", len(files))
	}
	return tview.NewModal().
		AddButtons([]string{"Yes", "No"}).
		SetText(text).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonIndex == 0 {
				done(ConfirmDeleteYes)
			} else {
				done(ConfirmDeleteNo)
			}
		})
}
