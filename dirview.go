package main

import (
	"io/fs"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type dirEntry struct {
	file     fs.DirEntry
	selected bool
}

type dirView struct {
	tvlist *tview.List

	dirPath string
	entries []dirEntry
}

func (tui *mainView) drawDir(no int) {
	dir := tui.dir0
	if no != 0 {
		dir = tui.dir1
	}
	inputField := tview.NewInputField().SetText(dir.dirPath)
	tui.tvgrid.AddItem(inputField, 0, no, 1, 1, 0, 0, false)
	dir.tvlist = tview.
		NewList().
		ShowSecondaryText(false)
	dir.tvlist.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlT || event.Rune() == ' ' {
			curr := tui.dir0.tvlist.GetCurrentItem()
			tui.dir0.entries[curr].selected = !tui.dir0.entries[curr].selected
			drawSelections(tui, tui.dir0)
			tui.dir0.tvlist.SetCurrentItem(curr + 1)
		} else if event.Rune() == 'j' {
			curr := tui.dir0.tvlist.GetCurrentItem()
			tui.dir0.tvlist.SetCurrentItem(curr + 1)
		} else if event.Rune() == 'k' {
			curr := tui.dir0.tvlist.GetCurrentItem()
			tui.dir0.tvlist.SetCurrentItem(curr - 1)
		}
		return event
	})
	tui.tvgrid.AddItem(dir.tvlist, 1, no, 1, 1, 0, 0, false)
}

func (tui *mainView) readDir(dir *dirView) {
	dir.tvlist.Clear()
	files, _ := os.ReadDir(dir.dirPath)
	for _, e := range files {
		dir.entries = append(dir.entries, dirEntry{file: e})
		dir.tvlist.AddItem(e.Name(), "", 0, nil)
	}
}

func drawSelections(tui *mainView, panel *dirView) {
	for i := 0; i < len(panel.entries); i++ {
		if panel.entries[i].selected {
			panel.tvlist.SetItemText(i, "[red::b]"+panel.entries[i].file.Name(), "")
		} else {
			panel.tvlist.SetItemText(i, panel.entries[i].file.Name(), "")
		}
	}
}
