package main

import (
	"fmt"
	"gofm/fsutils"
	"sync"

	"github.com/rivo/tview"
)

type existsView struct {
	element *tview.Modal

	file   string
	action fsutils.DstExistsAction
	wg     *sync.WaitGroup
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
			v.wg.Done()
		})
	return v
}

func (e *existsView) SetData(file string) *existsView {
	e.element.SetText(fmt.Sprintf("%s \n Entry exists. What to do?", file))
	return e
}

func (e *existsView) SetWg(wg *sync.WaitGroup) *existsView {
	e.wg = wg
	return e
}
