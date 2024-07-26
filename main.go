package main

import (
	"github.com/rivo/tview"
	"os"
	"os/user"
)

type PanelView struct {
	dirPath string
	ui_list *tview.List
}

type Tui struct {
	app            *tview.Application
	ui_mainGrid    *tview.Grid
	ui_pages       *tview.Pages
	ui_renameModal *tview.Modal
	leftPanel      *PanelView
	rightPanel     *PanelView
}

func main() {
	tui := &Tui{
		leftPanel:  &PanelView{},
		rightPanel: &PanelView{}}

	tui.app = tview.NewApplication()
	tui.app.EnableMouse(true)

	user, _ := user.Current()
	tui.leftPanel.dirPath = user.HomeDir
	tui.rightPanel.dirPath = user.HomeDir
	draw(tui)

	tui.updateDir(0)
	tui.updateDir(1)

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

func drawDir(tui *Tui, panel *PanelView, col int) {
	inputField := tview.NewInputField().SetText(panel.dirPath)
	tui.ui_mainGrid.AddItem(inputField, 0, col, 1, 1, 0, 0, false)
	panel.ui_list = tview.
		NewList().
		ShowSecondaryText(false)
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

func (tui *Tui) updateDir(dir int) {
	panel := tui.leftPanel
	if dir == 1 {
		panel = tui.rightPanel
	}
	files, _ := os.ReadDir(panel.dirPath)
	for _, e := range files {
		panel.ui_list.AddItem(e.Name(), "", 0, nil)
	}
}
