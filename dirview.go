package main

import (
	"cmp"
	"fmt"
	"gofm/fsutils"
	"log"
	"os"
	"path/filepath"
	"slices"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type dirEntry struct {
	path        string
	name        string
	displayName string
	isDir       bool
	selected    bool
}

type dirView struct {
	main     *mainView
	otherDir *dirView

	element   *tview.Flex
	list      *tview.List
	pathInput *tview.InputField

	no      int
	dirPath string
	entries []dirEntry
}

func newDirView(m *mainView, no int) *dirView {
	col := tview.NewFlex().SetDirection(tview.FlexColumnCSS)
	d := &dirView{entries: []dirEntry{}, main: m, no: no}
	d.pathInput = tview.NewInputField().SetText(d.dirPath)
	d.pathInput.SetBorder(true)
	d.pathInput.SetDoneFunc(func(key tcell.Key) {
		d.readDir(d.pathInput.GetText())
	})
	col.AddItem(d.pathInput, 3, 1, false)
	d.list = tview.
		NewList().
		ShowSecondaryText(false)
	d.list.SetBorder(true)
	d.list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlT || event.Rune() == ' ' {
			curr := d.list.GetCurrentItem()
			d.entries[curr].selected = !d.entries[curr].selected
			d.drawSelections()
			d.list.SetCurrentItem(curr + 1)
		} else if event.Rune() == 'j' {
			curr := d.list.GetCurrentItem()
			d.list.SetCurrentItem(curr + 1)
		} else if event.Rune() == 'k' {
			curr := d.list.GetCurrentItem()
			d.list.SetCurrentItem(curr - 1)
		} else if event.Key() == tcell.KeyTAB {
			m.app.SetFocus(d.otherDir.list)
			return nil
		} else if event.Key() == tcell.KeyEnter {
			err := d.handleOpenDirFromList()
			if err != nil {
				log.Fatalf("failed to call handleOpenDirFromList : %v", err)
			}
		} else if event.Key() == tcell.KeyF5 {
			d.handleCopyClick()
		} else if event.Key() == tcell.KeyF8 {
			d.handleDeleteClick()
		}
		return event
	})
	d.list.SetFocusFunc(func() {
		d.main.lastFocusedDir = d
	})
	col.AddItem(d.list, 0, 1, false)
	d.element = col
	return d
}

func (d *dirView) handleOpenDirFromList() error {
	selected := d.entries[d.list.GetCurrentItem()]
	newPath := filepath.Clean(d.dirPath + "/" + selected.name)
	stat, err := os.Stat(newPath)
	if err != nil {
		return fmt.Errorf("failed to get file stat %v", err)
	}
	if stat.IsDir() {
		d.readDir(newPath)
		d.pathInput.SetText(newPath)
	}
	return nil
}

func (d *dirView) handleCopyClick() {
	selected := d.getSelected()
	d.copyFrom(selected, 0, nil, 0)
}

func (d *dirView) copyFrom(selected []dirEntry, start int, act *fsutils.DstExistsAction, total int) {
	for i := start; i < len(selected); i++ {
		e := selected[i]
		if e.name == ".." {
			continue
		}
		src := e.path
		dst := filepath.Join(d.otherDir.dirPath, e.name)
		if src == dst {
			d.main.setStatus("Copy failed. Source and destination are the same!")
			return
		}

		if fsutils.Exists(dst) {
			// action was not selected for current fil
			if act == nil {
				d.main.showExists(dst, func(a fsutils.DstExistsAction) {
					d.copyFrom(selected, i, &a, total)
				})
				return
			} else {
				// if skip continue, if not copy to override
				if *act == fsutils.DstExistsActionSkip {
					act = nil
					continue
				}
				act = nil
			}
		}

		total++

		err := fsutils.Copy(src, dst)

		if err != nil {
			d.main.setStatus(fmt.Errorf("Copy failed : %w", err).Error())
		}
	}
	d.main.setStatus(fmt.Sprintf("Copy completed, %v entries created", total))
	d.otherDir.readDir(d.otherDir.dirPath)
}

func (d *dirView) getSelected() []dirEntry {
	var selected []dirEntry
	for _, e := range d.entries {
		if e.selected {
			selected = append(selected, e)
		}
	}
	if len(selected) == 0 {
		selected = append(selected, d.entries[d.list.GetCurrentItem()])
	}
	return selected
}

func (d *dirView) handleDeleteClick() {
	selected := d.getSelected()
	d.main.showConfirmDelete(selected, func(a ConfirmDeleteAction) {
		if a == ConfirmDeleteYes {
			for _, p := range selected {
				os.RemoveAll(p.path)
			}
			d.readDir(d.dirPath)
			d.main.setStatus(fmt.Sprintf("Delete completed, %v entries deleted", len(selected)))
		} else {
			d.main.setStatus("Delete canceled")
		}
	})
}

func (d *dirView) readDir(path string) {
	d.list.Clear()
	d.list.SetOffset(0, 0)
	d.entries = nil
	d.entries = append(d.entries, dirEntry{name: "..", displayName: "/.."})
	files, _ := os.ReadDir(path)
	for _, e := range files {
		displayName := e.Name()
		if e.IsDir() {
			displayName = "📁" + displayName
		} else {
			displayName = " " + displayName
		}
		d.entries = append(d.entries, dirEntry{
			path:        filepath.Join(path, e.Name()),
			name:        e.Name(),
			displayName: displayName,
			isDir:       e.IsDir()})
	}
	slices.SortFunc(d.entries, func(a, b dirEntry) int {
		return cmp.Or(
			func() int {
				if a.name == ".." && b.name != ".." {
					return -1
				} else if a.name != ".." && b.name == ".." {
					return 1
				} else {
					return 0
				}
			}(),
			func() int {
				if a.isDir && !b.isDir {
					return -1
				} else if !a.isDir && b.isDir {
					return 1
				} else {
					return 0
				}
			}(),
			cmp.Compare(a.name, b.name),
		)
	})
	for _, e := range d.entries {
		d.list.AddItem(e.displayName, "", 0, nil)
	}
	d.dirPath = path
}

func (d *dirView) drawSelections() {
	for i := 0; i < len(d.entries); i++ {
		if d.entries[i].selected {
			d.list.SetItemText(i, "[red::b]"+d.entries[i].displayName, "")
		} else {
			d.list.SetItemText(i, d.entries[i].displayName, "")
		}
	}
}
