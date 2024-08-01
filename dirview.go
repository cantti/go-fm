package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type dirEntry struct {
	path     string
	name     string
	selected bool
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
				log.Fatalf("key enter failed : %v", err)
			}
		} else if event.Key() == tcell.KeyF5 {
			err := d.handleCopyFileClick()
			if err != nil {
				log.Fatalf("key f5 failed : %v", err)
			}
		}
		return event
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
		return fmt.Errorf("handleOpenDirFromList failed %v", err)
	}
	if stat.IsDir() {
		d.readDir(newPath)
		d.pathInput.SetText(newPath)
	}
	return nil
}

func (d *dirView) handleCopyFileClick() error {
	var selected []dirEntry
	for _, e := range d.entries {
		if e.selected {
			selected = append(selected, e)
		}
	}
	if len(selected) == 0 {
		selected = append(selected, d.entries[d.list.GetCurrentItem()])
	}
	for _, e := range selected {
		stat, err := os.Stat(e.path)
		if err != nil {
			return fmt.Errorf("handleCopyFileClick failed : %v", err)
		}
		if !stat.IsDir() {
			src := e.path
			dest := filepath.Join(d.otherDir.dirPath, e.name)
			fsCopy(src, dest)
		}
	}
	d.readDir(d.dirPath)
	d.otherDir.readDir(d.otherDir.dirPath)
	return nil
}

func (d *dirView) readDir(path string) {
	d.list.Clear()
	d.list.SetOffset(0, 0)
	d.entries = nil
	d.entries = append(d.entries, dirEntry{name: ".."})
	files, _ := os.ReadDir(path)
	for _, e := range files {
		d.entries = append(d.entries, dirEntry{
			path: filepath.Join(path, e.Name()),
			name: e.Name()})
	}
	for _, e := range d.entries {
		d.list.AddItem(e.name, "", 0, nil)
	}
	d.dirPath = path
}

func (d *dirView) drawSelections() {
	for i := 0; i < len(d.entries); i++ {
		if d.entries[i].selected {
			d.list.SetItemText(i, "[red::b]"+d.entries[i].name, "")
		} else {
			d.list.SetItemText(i, d.entries[i].name, "")
		}
	}
}
