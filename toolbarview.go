package main

import (
	"fmt"

	"github.com/rivo/tview"
)

type toolbarView struct {
	element *tview.Flex
	main    *mainView
}

func newToolbarView(m *mainView) *toolbarView {
	t := &toolbarView{main: m}
	flex := tview.NewFlex().SetDirection(tview.FlexRowCSS)

	btnHelp := tview.NewButton(t.fmtBtn("F1", "Help")).
		SetSelectedFunc(m.showHelp)
	flex.AddItem(btnHelp, 0, 1, false).
		AddItem(tview.NewBox(), 1, 0, false)

	btnCopy := tview.NewButton(t.fmtBtn("F5", "Copy")).
		SetSelectedFunc(func() { m.lastFocusedDir.handleCopyClick() })
	flex.AddItem(btnCopy, 0, 1, false).
		AddItem(tview.NewBox(), 1, 0, false)

	// btnMove := tview.NewButton(t.fmtBtn("F6", "Move (not implemented)"))
	// flex.AddItem(btnMove, 0, 1, false)
	// flex.AddItem(tview.NewBox(), 1, 0, false)

	// btnRename := tview.NewButton(t.fmtBtn("Shift-F6", "Rename (not implemented)"))
	// btnRename.SetSelectedFunc(m.showRenameWin)
	// flex.AddItem(btnRename, 0, 1, false)
	// flex.AddItem(tview.NewBox(), 1, 0, false)

	btnDelete := tview.NewButton(t.fmtBtn("F8", "Delete")).
		SetSelectedFunc(func() { m.lastFocusedDir.handleDeleteClick() })
	flex.AddItem(btnDelete, 0, 1, false).
		AddItem(tview.NewBox(), 1, 0, false)

	btnQuit := tview.NewButton(t.fmtBtn("F10", "Quit")).
		SetSelectedFunc(m.showQuit)
	flex.AddItem(btnQuit, 0, 1, false).
		AddItem(tview.NewBox(), 1, 0, false)

	t.element = flex

	return t
}

func (t *toolbarView) fmtBtn(key string, text string) string {
	return fmt.Sprintf("[black:orange]%s[-:-] %s", key, text)
}
