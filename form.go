package main

import (
	"github.com/google/uuid"
	"github.com/rivo/tview"
)

// helper function to provide input form, returns vault record
func inputForm(app *tview.Application) VaultRecord {
	var vrec VaultRecord
	form := tview.NewForm()
	form.AddInputField("Name", "", 20, nil, nil)
	form.AddInputField("Login", "", 20, nil, nil)
	form.AddPasswordField("Password", "", 10, '*', nil)
	form.AddInputField("URL", "", 100, nil, nil)
	form.AddInputField("Note", "", 20, nil, nil)
	form.AddButton("Save", func() {
		vrec = saveForm(form)
		app.Stop()
	})
	form.AddButton("Quit", func() {
		app.Stop()
	})
	form.SetBorder(true).SetTitle("Form").SetTitleAlign(tview.AlignLeft)
	if err := app.SetRoot(form, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
	return vrec
}

// helper function to save input form
func saveForm(form *tview.Form) VaultRecord {
	name := form.GetFormItemByLabel("Name").(*tview.InputField).GetText()
	rurl := form.GetFormItemByLabel("URL").(*tview.InputField).GetText()
	username := form.GetFormItemByLabel("Login").(*tview.InputField).GetText()
	password := form.GetFormItemByLabel("Password").(*tview.InputField).GetText()
	note := form.GetFormItemByLabel("Note").(*tview.InputField).GetText()
	recLogin := VaultItem{Name: "login", Value: username}
	recPassword := VaultItem{Name: "password", Value: password}
	recUrl := VaultItem{Name: "url", Value: rurl}
	uid := uuid.NewString()
	rec := VaultRecord{ID: uid, Name: name, Items: []VaultItem{recLogin, recPassword, recUrl}, Note: note}
	return rec
}

// helper function to list vault records
func listForm(app *tview.Application, records []VaultRecord) {
	list := tview.NewList()
	for idx, rec := range records {
		list.AddItem(rec.ID, rec.Name, rune(idx), nil)
	}
	list.AddItem("Quit", "Press to exit", 'q', func() {
		app.Stop()
	})
	if err := app.SetRoot(list, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
