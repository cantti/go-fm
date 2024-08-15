package main

import (
	"fmt"
	"os"
)

type deleteCommand struct {
	entries []dirEntry
	dirView *dirView
}

func newDeleteCommand(dirView *dirView) command {
	command := &deleteCommand{dirView: dirView}
	return command
}

func (r *deleteCommand) execute() {
	r.entries = r.dirView.getSelected()
	r.showConfirmDelete(func(a ConfirmDeleteAction) {
		if a == ConfirmDeleteYes {
			for _, p := range r.entries {
				os.RemoveAll(p.path)
			}
			r.dirView.readDir(r.dirView.dirPath)
			r.dirView.main.setStatus(fmt.Sprintf("Delete completed, %v entries deleted", len(r.entries)))
		} else {
			r.dirView.main.setStatus("Delete canceled")
		}
	})
}

func (r *deleteCommand) showConfirmDelete(done func(a ConfirmDeleteAction)) {
	modal := newConfirmDeleteView(r.entries, func(a ConfirmDeleteAction) {
		r.dirView.main.pages.RemovePage("confirmDelete")
		r.dirView.main.app.SetFocus(r.dirView.list)
		done(a)
	})
	r.dirView.main.pages.AddPage("confirmDelete", modal, false, true)
}
