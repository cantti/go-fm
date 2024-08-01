package main

import "github.com/rivo/tview"

type renameView struct {
	element *tview.Grid
	form    *tview.Form
}

func newRenameView(m *mainView) *renameView {
	renameView := &renameView{}
	renameView.form = tview.NewForm().
		AddDropDown("Title", []string{"Mr.", "Ms.", "Mrs.", "Dr.", "Prof."}, 0, nil).
		AddInputField("First name", "", 20, nil, nil).
		AddInputField("Last name", "", 20, nil, nil).
		AddTextArea("Address", "", 40, 0, 0, nil).
		AddTextView("Notes", "This is just a demo.\nYou can enter whatever you wish.", 40, 2, true, false).
		AddCheckbox("Age 18+", false, nil).
		AddPasswordField("Password", "", 10, '*', nil).
		AddButton("Save", nil)

	renameView.form.SetBorder(true).SetTitle("Enter some data").SetTitleAlign(tview.AlignCenter)

	renameView.element = tview.NewGrid().
		SetColumns(0, 50, 0).
		SetRows(0, 50, 0).
		AddItem(renameView.form, 1, 1, 1, 1, 0, 0, true)

	return renameView
}
