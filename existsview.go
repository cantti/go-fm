package main

import (
	"fmt"

	"github.com/rivo/tview"
)

type DstExistsAction int

const (
	DstExistsActionNotSelected DstExistsAction = iota
	DstExistsActionOverWrite
	DstExistsActionSkip
)

func newExistsView(file string, done func(a DstExistsAction)) *tview.Modal {
	return tview.NewModal().
		AddButtons([]string{"Overwrite", "Skip"}).
		SetText(fmt.Sprintf("%s \n Entry exists. What to do?", file)).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonIndex == 0 {
				done(DstExistsActionOverWrite)
			} else {
				done(DstExistsActionSkip)
			}
		})
}
