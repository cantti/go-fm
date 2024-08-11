package main

import (
	"cmp"
	"fmt"
	"gofm/fsutils"
	"os"
	"os/user"
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
	element := tview.NewFlex().SetDirection(tview.FlexColumnCSS)
	d := &dirView{entries: []dirEntry{}, main: m, no: no}

	// add input
	d.pathInput = tview.NewInputField().SetText(d.dirPath)
	d.pathInput.SetFieldBackgroundColor(tcell.ColorBlack)
	d.pathInput.SetBorder(true)
	d.pathInput.SetDoneFunc(func(key tcell.Key) {
		d.readDir(d.pathInput.GetText())
	})
	element.AddItem(d.pathInput, 3, 1, false)

	// add toolbar
	toolbar := drawToolbar(d)
	element.AddItem(toolbar, 1, 0, false)

	// add list
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
			return nil

		} else if event.Rune() == 'j' {
			curr := d.list.GetCurrentItem()
			d.list.SetCurrentItem(curr + 1)

		} else if event.Rune() == 'k' {
			curr := d.list.GetCurrentItem()
			d.list.SetCurrentItem(curr - 1)

		} else if event.Key() == tcell.KeyTAB {
			m.app.SetFocus(d.otherDir.list)
			return nil

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
	d.list.SetSelectedFunc(func(index int, main string, second string, rune rune) {
		err := d.handleOpenDirFromList(index)
		if err != nil {
			d.main.setStatus(fmt.Sprintf("failed to call handleOpenDirFromList : %v", err))
		}
	})
	element.AddItem(d.list, 0, 1, false)

	d.element = element
	return d
}

func drawToolbar(d *dirView) *tview.Flex {
	toolbar := tview.NewFlex().SetDirection(tview.FlexRowCSS)
	toolbar.AddItem(tview.NewBox(), 1, 0, false)
	toolbar.AddItem(tview.NewButton("Home dir").SetSelectedFunc(d.openHomeDir), 0, 1, false)
	toolbar.AddItem(tview.NewBox(), 1, 0, false)
	toolbar.AddItem(tview.NewButton("Sync path").SetSelectedFunc(d.syncPath), 0, 1, false)
	toolbar.AddItem(tview.NewBox(), 1, 0, false)
	return toolbar
}

func (d *dirView) handleOpenDirFromList(index int) error {
	selected := d.entries[index]
	newPath := filepath.Clean(d.dirPath + "/" + selected.name)
	stat, err := os.Stat(newPath)
	if err != nil {
		return fmt.Errorf("failed to get file stat %v", err)
	}
	if stat.IsDir() {
		d.readDir(newPath)
	}
	return nil
}

func (d *dirView) handleCopyClick() {
	selected := d.getSelected()
	command := newCopyCommand(selected, d)
	command.execute()
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
	if !fsutils.Exists(path) {
		d.pathInput.SetText(d.dirPath)
		return
	}
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
	d.pathInput.SetText(d.dirPath)
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

func (d *dirView) syncPath() {
	d.readDir(d.otherDir.dirPath)
}

func (d *dirView) openHomeDir() {
	user, _ := user.Current()
	d.readDir(user.HomeDir)
}
