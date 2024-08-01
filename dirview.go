package main

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type dirEntry struct {
	file     fs.DirEntry
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

func newDir(m *mainView, no int) *dirView {
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
			d.handleOpenDirFromList()
		} else if event.Key() == tcell.KeyF5 {
			d.handleCopyFileClick()
		}
		return event
	})
	col.AddItem(d.list, 0, 1, false)
	d.element = col
	return d
}

func (d *dirView) handleOpenDirFromList() {
	mainText, _ := d.list.GetItemText(d.list.GetCurrentItem())
	newPath := filepath.Clean(d.dirPath + "/" + mainText)
	stat, _ := os.Stat(newPath)
	if stat.IsDir() {
		d.readDir(newPath)
		d.pathInput.SetText(newPath)
	}
}

func (d *dirView) handleCopyFileClick() {
	mainText, _ := d.list.GetItemText(d.list.GetCurrentItem())
	src := filepath.Join(d.dirPath, mainText)
	dest := filepath.Join(d.otherDir.dirPath, mainText)
	fsCopy(src, dest)
	d.readDir(d.dirPath)
	d.otherDir.readDir(d.otherDir.dirPath)
}

func (d *dirView) readDir(path string) {
	d.list.Clear()
	d.list.SetOffset(0, 0)
	d.list.AddItem("..", "", 0, nil)
	files, _ := os.ReadDir(path)
	for _, e := range files {
		d.entries = append(d.entries, dirEntry{file: e})
		d.list.AddItem(e.Name(), "", 0, nil)
	}
	d.dirPath = path
}

func (d *dirView) drawSelections() {
	for i := 0; i < len(d.entries); i++ {
		if d.entries[i].selected {
			d.list.SetItemText(i, "[red::b]"+d.entries[i].file.Name(), "")
		} else {
			d.list.SetItemText(i, d.entries[i].file.Name(), "")
		}
	}
}
