package main

import (
	"fmt"
	"gofm/fsutils"

	"github.com/rivo/tview"
)

type existsView struct {
	element *tview.Modal

	file   string
	action fsutils.DstExistsAction
}

func newExistsView(m *mainView) *existsView {
	v := &existsView{}
	v.element = tview.NewModal().
		AddButtons([]string{"Overwrite", "Skip"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonIndex == 0 {
				v.action = fsutils.DstExistsActionOverWrite
			} else {
				v.action = fsutils.DstExistsActionSkip
			}
			m.wg.Done()
		})
	return v
}

func (e *existsView) SetData(file string) {
	e.element.SetText(fmt.Sprintf("%s \n Entry exists. What to do?", file))
}
