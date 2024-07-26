package main

import (
	"io/fs"
	"os"
	"os/user"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type DirEntry struct {
	file     fs.DirEntry
	selected bool
}

type Panel struct {
	dirPath string
	ui_list *tview.List
	entries []DirEntry
}

type Tui struct {
	app            *tview.Application
	ui_mainGrid    *tview.Grid
	ui_pages       *tview.Pages
	ui_renameModal *tview.Modal
	leftPanel      *Panel
	rightPanel     *Panel
}

var sel int

func main() {
	tui := &Tui{leftPanel: &Panel{entries: []DirEntry{}}, rightPanel: &Panel{entries: []DirEntry{}}}

	tui.app = tview.NewApplication()
	tui.app.EnableMouse(true)

	user, _ := user.Current()
	tui.leftPanel.dirPath = user.HomeDir
	tui.rightPanel.dirPath = user.HomeDir
	draw(tui)

	tui.readDir(0)
	tui.readDir(1)

	tui.app.Run()
}

func draw(tui *Tui) {
	tui.ui_pages = tview.NewPages()
	tui.app.SetRoot(tui.ui_pages, true)

	tui.ui_mainGrid = tview.NewGrid().
		SetRows(1, 0, 1).
		SetBorders(true)

	drawDir(tui, tui.leftPanel, 0)
	drawDir(tui, tui.rightPanel, 1)

	bottomToolbar := tui.drawBottomToolbar()
	tui.ui_mainGrid.AddItem(bottomToolbar, 2, 0, 1, 2, 0, 0, false)

	modal := tview.NewModal().
		SetText("Do you want to quit the application?").
		AddButtons([]string{"Quit", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Quit" {
				tui.app.Stop()
			} else {
				tui.ui_pages.HidePage("modal")
			}
		})

	tui.ui_renameModal = tview.NewModal().
		SetText("hello")

	tui.ui_pages.AddPage("main", tui.ui_mainGrid, true, true)
	tui.ui_pages.AddPage("modal", modal, true, false)
	tui.ui_pages.AddPage("rename", tui.ui_renameModal, true, false)

	tui.app.SetFocus(tui.leftPanel.ui_list)
}

func drawDir(tui *Tui, panel *Panel, col int) {
	inputField := tview.NewInputField().SetText(panel.dirPath)
	tui.ui_mainGrid.AddItem(inputField, 0, col, 1, 1, 0, 0, false)
	panel.ui_list = tview.
		NewList().
		ShowSecondaryText(false)
	panel.ui_list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlT || event.Rune() == ' ' {
			curr := tui.leftPanel.ui_list.GetCurrentItem()
			tui.leftPanel.entries[curr].selected = !tui.leftPanel.entries[curr].selected
			drawSelections(tui, tui.leftPanel)
			tui.leftPanel.ui_list.SetCurrentItem(curr + 1)
		} else if event.Rune() == 'j' {
			curr := tui.leftPanel.ui_list.GetCurrentItem()
			tui.leftPanel.ui_list.SetCurrentItem(curr + 1)
		} else if event.Rune() == 'k' {
			curr := tui.leftPanel.ui_list.GetCurrentItem()
			tui.leftPanel.ui_list.SetCurrentItem(curr - 1)
		}
		return event
	})
	tui.ui_mainGrid.AddItem(panel.ui_list, 1, col, 1, 1, 0, 0, false)
}

func (tui *Tui) drawBottomToolbar() *tview.Grid {
	buttonGrid := tview.NewGrid().SetColumns(0, 0, 0, 0).SetGap(0, 1)
	buttonGrid.AddItem(tview.NewButton("copy"), 0, 0, 1, 1, 0, 0, false)
	buttonGrid.AddItem(tview.NewButton("move"), 0, 1, 1, 1, 0, 0, false)
	buttonGrid.AddItem(tview.NewButton("rename").
		SetSelectedFunc(tui.showRenameWin), 0, 2, 1, 1, 0, 0, false)
	buttonGrid.AddItem(tview.NewButton("quit").
		SetSelectedFunc(func() { tui.ui_pages.ShowPage("modal") }), 0, 3, 1, 1, 0, 0, false)
	return buttonGrid
}

func (tui *Tui) showRenameWin() {
	i := tui.leftPanel.ui_list.GetCurrentItem()
	text, _ := tui.leftPanel.ui_list.GetItemText(i)
	tui.ui_renameModal.SetText(text)
	tui.ui_pages.ShowPage("rename")
}

func drawSelections(tui *Tui, panel *Panel) {
	for i := 0; i < len(panel.entries); i++ {
		if panel.entries[i].selected {
			panel.ui_list.SetItemText(i, "[red::b]"+panel.entries[i].file.Name(), "")
		} else {
			panel.ui_list.SetItemText(i, panel.entries[i].file.Name(), "")
		}
	}
}

func (tui *Tui) readDir(dir int) {
	panel := tui.leftPanel
	if dir == 1 {
		panel = tui.rightPanel
	}

	panel.ui_list.Clear()

	files, _ := os.ReadDir(panel.dirPath)
	for _, e := range files {
		panel.entries = append(panel.entries, DirEntry{file: e})
		panel.ui_list.AddItem(e.Name(), "", 0, nil)
	}
}
