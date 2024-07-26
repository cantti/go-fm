package main

import (
	"github.com/rivo/tview"
	"os"
	"os/user"
)

type PanelView struct {
	dirPath string
	dir     *tview.List
}

type Tui struct {
	app         *tview.Application
	pages       *tview.Pages
	leftPanel   *PanelView
	rightPanel  *PanelView
	renameModal *tview.Modal
}

func main() {
	test := 123
	tui := Tui{
		leftPanel:  &PanelView{},
		rightPanel: &PanelView{}}

	tui.app = tview.NewApplication()
	tui.app.EnableMouse(true)

	user, _ := user.Current()
	tui.leftPanel.dirPath = user.HomeDir
	tui.rightPanel.dirPath = user.HomeDir
	test = test + 1
	tui.render()

	tui.updateDir(0)
	tui.updateDir(1)

	tui.app.Run()
}

func (tui *Tui) render() {
	tui.pages = tview.NewPages()
	tui.app.SetRoot(tui.pages, true)

	mainGrid := tview.NewGrid().
		SetRows(0, 1).
		SetBorders(true)

	tui.leftPanel.dir = tui.renderDir()
	tui.rightPanel.dir = tui.renderDir()

	mainGrid.AddItem(tui.leftPanel.dir, 0, 0, 1, 1, 0, 0, false)
	mainGrid.AddItem(tui.rightPanel.dir, 0, 1, 1, 1, 0, 0, false)

	bottomToolbar := tui.renderBottomToolbar()
	mainGrid.AddItem(bottomToolbar, 1, 0, 1, 2, 0, 0, false)

	modal := tview.NewModal().
		SetText("Do you want to quit the application?").
		AddButtons([]string{"Quit", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Quit" {
				tui.app.Stop()
			} else {
				tui.pages.HidePage("modal")
			}
		})

	tui.renameModal = tview.NewModal().
		SetText("hello")

	tui.pages.AddPage("main", mainGrid, true, true)
	tui.pages.AddPage("modal", modal, true, false)
	tui.pages.AddPage("rename", tui.renameModal, true, false)

	tui.app.SetFocus(tui.leftPanel.dir)
}

func (tui *Tui) renderDir() *tview.List {
	list := tview.NewList()
	list.ShowSecondaryText(false)
	return list
}

func (tui *Tui) renderBottomToolbar() *tview.Grid {
	buttonGrid := tview.NewGrid().SetColumns(0, 0, 0, 0).SetGap(0, 1)
	buttonGrid.AddItem(tview.NewButton("copy"), 0, 0, 1, 1, 0, 0, false)
	buttonGrid.AddItem(tview.NewButton("move"), 0, 1, 1, 1, 0, 0, false)
	buttonGrid.AddItem(tview.NewButton("rename").
		SetSelectedFunc(tui.showRenameWin), 0, 2, 1, 1, 0, 0, false)
	buttonGrid.AddItem(tview.NewButton("quit").
		SetSelectedFunc(func() { tui.pages.ShowPage("modal") }), 0, 3, 1, 1, 0, 0, false)
	return buttonGrid
}

func (tui *Tui) showRenameWin() {
	i := tui.leftPanel.dir.GetCurrentItem()
	text, _ := tui.leftPanel.dir.GetItemText(i)
	tui.renameModal.SetText(text)
	tui.pages.ShowPage("rename")
}

func (tui *Tui) updateDir(dir int) {
	var list *tview.List
	var path string
	if dir == 0 {
		list = tui.leftPanel.dir
		path = tui.leftPanel.dirPath
	} else {
		list = tui.rightPanel.dir
		path = tui.rightPanel.dirPath
	}
	files, _ := os.ReadDir(path)
	for _, e := range files {
		list.AddItem(e.Name(), "", 0, nil)
	}
}
