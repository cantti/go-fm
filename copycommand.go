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

func newCopyCommand(entries []dirEntry, dirView *dirView) *copyCommand {
	command := &copyCommand{
		entries: entries,
		dirView: dirView}
	return command
}

func (r *copyCommand) execute() {
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
				r.dirView.main.showExists(dst, func(a DstExistsAction) {
					r.from = i
					r.dstExistsAction = a
					r.execute()
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
