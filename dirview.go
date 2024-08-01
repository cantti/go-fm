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
	list      *tview.List
	pathInput *tview.InputField

	dirPath string
	entries []dirEntry
}

func (m *mainView) drawDirs() {
	row := tview.NewFlex().SetDirection(tview.FlexRowCSS)
	for i := 0; i < 2; i++ {
		col := tview.NewFlex().SetDirection(tview.FlexColumnCSS)
		dir := m.dir0
		if i == 1 {
			dir = m.dir1
		}
		dir.pathInput = tview.NewInputField().SetText(dir.dirPath)
		dir.pathInput.SetBorder(true)
		dir.pathInput.SetDoneFunc(func(key tcell.Key) {
			m.readDir(dir, dir.pathInput.GetText())
		})
		col.AddItem(dir.pathInput, 3, 1, false)
		dir.list = tview.
			NewList().
			ShowSecondaryText(false)
		dir.list.SetBorder(true)
		dir.list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyCtrlT || event.Rune() == ' ' {
				curr := dir.list.GetCurrentItem()
				dir.entries[curr].selected = !dir.entries[curr].selected
				drawSelections(m, dir)
				dir.list.SetCurrentItem(curr + 1)
			} else if event.Rune() == 'j' {
				curr := dir.list.GetCurrentItem()
				dir.list.SetCurrentItem(curr + 1)
			} else if event.Rune() == 'k' {
				curr := dir.list.GetCurrentItem()
				dir.list.SetCurrentItem(curr - 1)
			} else if event.Key() == tcell.KeyTAB {
				if i == 0 {
					m.app.SetFocus(m.dir1.list)
				} else {
					m.app.SetFocus(m.dir0.list)
				}
				return nil
			} else if event.Key() == tcell.KeyEnter {
				m.openDirFromList(dir)
			}
			return event
		})
		col.AddItem(dir.list, 0, 1, false)
		row.AddItem(col, 0, 1, false)
		row.AddItem(tview.NewBox(), 1, 0, false)
	}
	m.flexCol.AddItem(row, 0, 1, false)
}

func (m *mainView) openDirFromList(dir *dirView) {
	mainText, _ := dir.list.GetItemText(dir.list.GetCurrentItem())
	newPath := filepath.Clean(dir.dirPath + "/" + mainText)
	stat, _ := os.Stat(newPath)
	if stat.IsDir() {
		m.readDir(dir, newPath)
		dir.pathInput.SetText(newPath)
	}
}

func (tui *mainView) readDir(dir *dirView, path string) {
	dir.list.Clear()
	dir.list.SetOffset(0, 0)
	dir.list.AddItem("..", "", 0, nil)
	files, _ := os.ReadDir(path)
	for _, e := range files {
		dir.entries = append(dir.entries, dirEntry{file: e})
		dir.list.AddItem(e.Name(), "", 0, nil)
	}
	dir.dirPath = path
}

func drawSelections(tui *mainView, panel *dirView) {
	for i := 0; i < len(panel.entries); i++ {
		if panel.entries[i].selected {
			panel.list.SetItemText(i, "[red::b]"+panel.entries[i].file.Name(), "")
		} else {
			panel.list.SetItemText(i, panel.entries[i].file.Name(), "")
		}
	}
}
