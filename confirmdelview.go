package main

import (
	"fmt"

	"github.com/rivo/tview"
)

type ConfirmDeleteAction int

const (
	ConfirmDeleteYes ConfirmDeleteAction = iota
	ConfirmDeleteNo
)

type confirmDeleteView struct {
	main    *mainView
	element *tview.Modal

	action ConfirmDeleteAction
}

func newConfirmDeleteView(m *mainView) *confirmDeleteView {
	v := &confirmDeleteView{}
	v.element = tview.NewModal().
		AddButtons([]string{"Yes", "No"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonIndex == 0 {
				v.action = ConfirmDeleteYes
			} else {
				v.action = ConfirmDeleteNo
			}
			m.hideConfirmDelete()
		})
	v.SetData(nil)
	return v
}

func (e *confirmDeleteView) SetData(files []string) {
	e.action = ConfirmDeleteYes
	if len(files) == 1 {
		e.element.SetText(fmt.Sprintf("%s \n Do you want to delete it?", files[0]))
	} else {
		e.element.SetText(fmt.Sprintf("Do you want to delete %v files?", len(files)))
	}
}
