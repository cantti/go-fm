package main

import (
	"fmt"
	"sync"

	"github.com/rivo/tview"
)

type mainView struct {
	app       *tview.Application
	element   *tview.Flex
	pages     *tview.Pages
	statusBar *tview.TextView

	renameView   *renameView
	dir0         *dirView
	dir1         *dirView
	toolbar      *toolbarView
	existsView   *existsView
	existsViewWg sync.WaitGroup
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

	m.toolbar = newToolbarView(m)
	m.element.AddItem(m.toolbar.element, 1, 0, false)

	drawStatusBar(m)

	// bottom padding
	m.element.AddItem(tview.NewBox(), 1, 0, false)

	quitModal := tview.NewModal().
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
	m.pages.AddPage("modal", quitModal, true, false)

	m.renameView = newRenameView(m)
	m.pages.AddPage("rename", m.renameView.element, true, false)

	m.existsView = newExistsView(m)
	m.pages.AddPage("exists", m.existsView.element, true, false)

	m.app.SetFocus(m.dir0.list)
}

func drawStatusBar(m *mainView) {
	m.statusBar = tview.NewTextView()
	m.element.AddItem(m.statusBar, 1, 0, false)
	m.statusBar.SetDynamicColors(true)
	m.setStatus("")
}

func (tui *mainView) showRenameWin() {
	tui.pages.ShowPage("rename")
}

func (m *mainView) showExists(file string) {
	m.existsViewWg.Add(1)
	m.existsView.SetData(file)
	m.pages.ShowPage("exists")
}

func (m *mainView) hideExists() {
	m.existsViewWg.Done()
	m.pages.HidePage("exists")
}

func (m *mainView) setStatus(text string) {
	m.statusBar.Clear()
	fmt.Fprintf(m.statusBar, "[orange:]%s[-:-] %s", "Status:", text)
}
