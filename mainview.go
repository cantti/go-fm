package main

import (
	// "github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type mainView struct {
	app     *tview.Application
	element *tview.Flex
	pages   *tview.Pages

	renameView *renameView
	dir0       *dirView
	dir1       *dirView
}

func newMainView() *mainView {
	mainView := &mainView{
		renameView: &renameView{}}
	mainView.app = tview.NewApplication()
	mainView.app.EnableMouse(true)
	mainView.draw()
	return mainView
}

func (m *mainView) draw() {
	m.pages = tview.NewPages()
	m.app.SetRoot(m.pages, true)

	m.element = tview.NewFlex().SetDirection(tview.FlexColumnCSS)

	m.dir0 = newDirView(m, 0)
	m.dir1 = newDirView(m, 1)
	m.dir0.otherDir = m.dir1
	m.dir1.otherDir = m.dir0

	m.element.AddItem(
		tview.NewFlex().
			SetDirection(tview.FlexRowCSS).
			AddItem(m.dir0.element, 0, 1, false).
			AddItem(m.dir1.element, 0, 1, false),
		0, 1, false)

	m.drawBottomToolbar()

	// bottom padding
	m.element.AddItem(tview.NewBox(), 1, 0, false)

	modal := tview.NewModal().
		SetText("Do you want to quit the application?").
		AddButtons([]string{"Quit", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Quit" {
				m.app.Stop()
			} else {
				m.pages.HidePage("modal")
			}
		})

	m.pages.AddPage("main", m.element, true, true)
	m.pages.AddPage("modal", modal, true, false)

	m.renameView = newRenameView(m)
	m.pages.AddPage("rename", m.renameView.element, true, false)

	m.app.SetFocus(m.dir0.list)
}

func (m *mainView) drawBottomToolbar() {
	flexCol := tview.NewFlex()

	btnCopy := tview.NewButton("Copy")

	btnMove := tview.NewButton("Move")

	btnRename := tview.NewButton("Rename")
	btnRename.SetSelectedFunc(m.showRenameWin)

	btnQuit := tview.NewButton("Quit")
	btnQuit.SetSelectedFunc(func() { m.pages.ShowPage("modal") })

	flexCol.AddItem(btnCopy, 0, 1, false)
	flexCol.AddItem(tview.NewBox(), 1, 0, false)
	flexCol.AddItem(btnMove, 0, 1, false)
	flexCol.AddItem(tview.NewBox(), 1, 0, false)
	flexCol.AddItem(btnRename, 0, 1, false)
	flexCol.AddItem(tview.NewBox(), 1, 0, false)
	flexCol.AddItem(btnQuit, 0, 1, false)

	m.element.AddItem(flexCol, 1, 0, false)
}

func (tui *mainView) showRenameWin() {
	tui.pages.ShowPage("rename")
}

func modal(p tview.Primitive, width, height int) tview.Primitive {
	return tview.NewGrid().
		SetColumns(0, width, 0).
		SetRows(0, height, 0).
		AddItem(p, 1, 1, 1, 1, 0, 0, true)
}
