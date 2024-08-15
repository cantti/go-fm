package main

import (
	"fmt"
	"gofm/fsutils"
	"path/filepath"
)

type copyCommand struct {
	dirView         *dirView
	entries         []dirEntry
	dstExistsAction DstExistsAction
	from            int
	total           int
}

func newCopyCommand(dirView *dirView) command {
	command := &copyCommand{dirView: dirView}
	return command
}

func (r *copyCommand) execute() {
	// reset state
	r.entries = r.dirView.getSelected()
	r.from = 0
	r.total = 0
	r.copy()
}

func (r *copyCommand) copy() {
	for i := r.from; i < len(r.entries); i++ {
		// cache old action and set actual to not selected
		act := r.dstExistsAction
		r.dstExistsAction = DstExistsActionNotSelected
		e := r.entries[i]
		if e.name == ".." {
			continue
		}
		src := e.path
		dst := filepath.Join(r.dirView.otherDir.dirPath, e.name)
		if src == dst {
			r.dirView.main.setStatus("Copy failed. Source and destination are the same!")
			return
		}

		if fsutils.Exists(dst) {
			// action was not selected for current fil
			if act == DstExistsActionNotSelected {
				r.showExists(dst, func(a DstExistsAction) {
					r.from = i
					r.dstExistsAction = a
					r.copy()
				})
				return
			} else {
				// if skip continue, if not copy to override
				if act == DstExistsActionSkip {
					continue
				}
			}
		}

		r.total++

		err := fsutils.Copy(src, dst)

		if err != nil {
			r.dirView.main.setStatus(fmt.Errorf("Copy failed : %w", err).Error())
		}
	}
	r.dirView.main.setStatus(fmt.Sprintf("Copy completed, %v entries created", r.total))
	r.dirView.otherDir.readDir(r.dirView.otherDir.dirPath)
}

func (r *copyCommand) showExists(file string, done func(a DstExistsAction)) {
	modal := newExistsView(file, func(a DstExistsAction) {
		r.dirView.main.pages.RemovePage("exists")
		r.dirView.main.app.SetFocus(r.dirView.list)
		done(a)
	})
	r.dirView.main.pages.AddPage("exists", modal, false, true)
}
