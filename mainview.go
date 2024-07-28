package main

import "github.com/rivo/tview"

type mainView struct {
	tvapp   *tview.Application
	tvgrid  *tview.Grid
	tvpages *tview.Pages
	tvmodal *tview.Grid

	renameView *renameView
	dir0       *dirView
	dir1       *dirView
}

func newMainView() *mainView {
	mainView := &mainView{
		dir0:       &dirView{entries: []dirEntry{}},
		dir1:       &dirView{entries: []dirEntry{}},
		renameView: &renameView{}}
	mainView.tvapp = tview.NewApplication()
	mainView.tvapp.EnableMouse(true)
	mainView.draw()
	return mainView
}

func (tui *mainView) draw() {
	tui.tvpages = tview.NewPages()
	tui.tvapp.SetRoot(tui.tvpages, true)

	tui.tvgrid = tview.NewGrid().
		SetRows(1, 0, 1).
		SetBorders(true)

	tui.drawDir(0)
	tui.drawDir(1)

	bottomToolbar := tui.drawBottomToolbar()
	tui.tvgrid.AddItem(bottomToolbar, 2, 0, 1, 2, 0, 0, false)

	modal := tview.NewModal().
		SetText("Do you want to quit the application?").
		AddButtons([]string{"Quit", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Quit" {
				tui.tvapp.Stop()
			} else {
				tui.tvpages.HidePage("modal")
			}
		})

	tui.tvpages.AddPage("main", tui.tvgrid, true, true)
	tui.tvpages.AddPage("modal", modal, true, false)

	tui.drawRenameView()

	tui.tvapp.SetFocus(tui.dir0.tvlist)
}

func (tui *mainView) drawBottomToolbar() *tview.Grid {
	buttonGrid := tview.NewGrid().SetColumns(0, 0, 0, 0).SetGap(0, 1)
	buttonGrid.AddItem(tview.NewButton("copy"), 0, 0, 1, 1, 0, 0, false)
	buttonGrid.AddItem(tview.NewButton("move"), 0, 1, 1, 1, 0, 0, false)
	buttonGrid.AddItem(tview.NewButton("rename").
		SetSelectedFunc(tui.showRenameWin), 0, 2, 1, 1, 0, 0, false)
	buttonGrid.AddItem(tview.NewButton("quit").
		SetSelectedFunc(func() { tui.tvpages.ShowPage("modal") }), 0, 3, 1, 1, 0, 0, false)
	return buttonGrid
}

func (tui *mainView) showRenameWin() {
	tui.tvpages.ShowPage("rename")
}

func modal(p tview.Primitive, width, height int) tview.Primitive {
	return tview.NewGrid().
		SetColumns(0, width, 0).
		SetRows(0, height, 0).
		AddItem(p, 1, 1, 1, 1, 0, 0, true)
}
