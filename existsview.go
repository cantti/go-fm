package main

import (
	"fmt"
	"gofm/fsutils"

	"github.com/rivo/tview"
)

func newExistsView(file string, done func(a fsutils.DstExistsAction)) *tview.Modal {
	return tview.NewModal().
		AddButtons([]string{"Overwrite", "Skip"}).
		SetText(fmt.Sprintf("%s \n Entry exists. What to do?", file)).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonIndex == 0 {
				done(fsutils.DstExistsActionOverWrite)
			} else {
				done(fsutils.DstExistsActionSkip)
			}
		})
}
