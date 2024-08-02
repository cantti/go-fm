package main

import (
	"fmt"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type mainView struct {
	app       *tview.Application
	element   *tview.Flex
	pages     *tview.Pages
	statusBar *tview.TextView
	quitView  *tview.Modal

	modalWg sync.WaitGroup

	renameView        *renameView
	dir0              *dirView
	dir1              *dirView
	lastFocusedDir    *dirView
	toolbar           *toolbarView
	existsView        *existsView
	confirmDeleteView *confirmDeleteView
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

	m.element.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyF10 {
			m.showQuit()
		}
		return event
	})

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
	m.pages.AddPage("main", m.element, true, true)

	m.drawStatusBar()

	// bottom padding
	m.element.AddItem(tview.NewBox(), 1, 0, false)

	m.drawQuit()

	m.renameView = newRenameView(m)
	m.pages.AddPage("rename", m.renameView.element, true, false)

	m.existsView = newExistsView(m)
	m.pages.AddPage("exists", m.existsView.element, true, false)

	m.confirmDeleteView = newConfirmDeleteView(m)
	m.pages.AddPage("confirmDelete", m.confirmDeleteView.element, true, false)

	m.app.SetFocus(m.dir0.list)
}

func (m *mainView) drawStatusBar() {
	m.statusBar = tview.NewTextView()
	m.element.AddItem(m.statusBar, 1, 0, false)
	m.statusBar.SetDynamicColors(true)
	m.setStatus("")
}

func (m *mainView) drawQuit() {
	m.quitView = tview.NewModal().
		SetText("Do you want to quit the application?").
		AddButtons([]string{"Quit", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Quit" {
				m.app.Stop()
			} else {
				m.pages.HidePage("quit")
			}
		})
	m.pages.AddPage("quit", m.quitView, true, false)
}

func (tui *mainView) showRenameWin() {
	tui.pages.ShowPage("rename")
}

func (m *mainView) showExists(file string) {
	m.modalWg.Add(1)
	m.existsView.SetData(file)
	m.pages.ShowPage("exists")
}

func (m *mainView) hideExists() {
	m.modalWg.Done()
	m.pages.HidePage("exists")
	m.app.SetFocus(m.lastFocusedDir.list)
}

func (m *mainView) showQuit() {
	m.pages.ShowPage("quit")
}

func (m *mainView) hideQuit() {
	m.pages.HidePage("quit")
	m.app.SetFocus(m.lastFocusedDir.list)
}

func (m *mainView) showConfirmDelete(files []string) {
	m.modalWg.Add(1)
	m.confirmDeleteView.SetData(files)
	m.pages.ShowPage("confirmDelete")
}

func (m *mainView) hideConfirmDelete() {
	m.modalWg.Done()
	m.pages.HidePage("confirmDelete")
	m.app.SetFocus(m.lastFocusedDir.list)
}

func (m *mainView) setStatus(text string) {
	m.statusBar.Clear()
	fmt.Fprintf(m.statusBar, "[orange:]%s[-:-] %s", "Status:", text)
}
