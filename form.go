package main

import (
	"log"

	tcell "github.com/gdamore/tcell/v2"
	"github.com/google/uuid"
	"github.com/rivo/tview"
)

// helper function to make secret prompt
func lockView(app *tview.Application, verbose int) string {
	//     defer app.Stop()
	form := tview.NewForm()
	form.AddPasswordField("Password", "", 50, '*', nil)
	if err := app.SetRoot(form, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
	password := form.GetFormItemByLabel("Password").(*tview.InputField).GetText()
	if verbose > 0 {
		log.Printf("vault secret '%s'", password)
	}
	return password
}

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
func listForm(app *tview.Application, vault *Vault) {
	list := tview.NewList()
	for idx, rec := range vault.Records {
		list.AddItem(rec.ID, rec.Name, rune(idx), nil)
	}
	list.AddItem("Quit", "Press to exit", 'q', func() {
		app.Stop()
	})
	if err := app.SetRoot(list, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

// helper function to present recordForm
func recordForm(app *tview.Application, form *tview.Form, index int, vault *Vault) *tview.Form {
	rec := vault.Records[index]
	name, rurl, login, password, note := rec.Details()
	form.Clear(true) // clear the form
	form.AddInputField("Name", name, 100, nil, nil)
	form.AddInputField("Login", login, 100, nil, nil)
	form.AddPasswordField("Password", password, 100, '*', nil)
	form.AddInputField("URL", rurl, 100, nil, nil)
	form.AddInputField("Note", note, 100, nil, nil)
	form.AddButton("Save", func() {
		uid := rec.ID
		name := form.GetFormItemByLabel("Name").(*tview.InputField).GetText()
		rurl := form.GetFormItemByLabel("URL").(*tview.InputField).GetText()
		username := form.GetFormItemByLabel("Login").(*tview.InputField).GetText()
		password := form.GetFormItemByLabel("Password").(*tview.InputField).GetText()
		note := form.GetFormItemByLabel("Note").(*tview.InputField).GetText()
		recLogin := VaultItem{Name: "login", Value: username}
		recPassword := VaultItem{Name: "password", Value: password}
		recUrl := VaultItem{Name: "url", Value: rurl}
		rec := VaultRecord{ID: uid, Name: name, Items: []VaultItem{recLogin, recPassword, recUrl}, Note: note}
		vault.Update(rec)
		vault.Write()
	})
	form.SetBorder(true).SetTitle("Form").SetTitleAlign(tview.AlignLeft)
	return form
}

// helper function to build our application grid view
func gridView(app *tview.Application, vault *Vault) {
	// add search bar
	find := tview.NewForm()
	find.AddInputField("Search", "", 80, nil, nil)
	find.AddButton("Find", func() {
		app.Stop()
	})
	find.AddButton("Quit", func() {
		app.Stop()
	})
	find.SetBorder(true).SetTitle("Search").SetTitleAlign(tview.AlignLeft)

	// set current record form view
	form := tview.NewForm()
	form = recordForm(app, form, 0, vault)

	// set main record list
	main := tview.NewList()
	for _, rec := range vault.Records {
		main.AddItem(rec.Name, rec.ID, rune('-'), nil)
	}
	main.AddItem("Quit", "Press to exit", 'q', func() {
		app.Stop()
	})
	main.SetBorder(true).SetTitle("Records")
	main.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		if index < len(vault.Records) {
			form = recordForm(app, form, index, vault)
		}
	})

	// construct grid view
	grid := tview.NewGrid()
	grid.SetBorders(false)

	// Layout for screens wider than 100 cells.
	grid.AddItem(find, 0, 0, 1, 2, 20, 100, false)
	grid.AddItem(main, 1, 0, 3, 1, 0, 100, true) // default focus, index 1
	grid.AddItem(form, 1, 1, 3, 1, 0, 100, false)

	focusIndex := 1 // defaul focus index

	grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		key := event.Key()
		switch key {
		//         case tcell.KeyRune:
		//             switch event.Rune() {
		//             case 'T', 't':
		//                 app.SetFocus(tags)
		//                 return event
		//             case 'F', 'f':
		//                 app.SetFocus(form)
		//                 return event
		//             case 'M', 'm':
		//                 app.SetFocus(main)
		//                 return event
		//             }
		case tcell.KeyCtrlN:
			if focusIndex == 0 {
				app.SetFocus(main)
				focusIndex = 1
			} else if focusIndex == 1 {
				app.SetFocus(form)
				focusIndex = 2
			} else if focusIndex == 2 {
				app.SetFocus(find)
				focusIndex = 0
			}
			return event
		case tcell.KeyCtrlB:
			if focusIndex == 0 {
				app.SetFocus(form)
				focusIndex = 2
			} else if focusIndex == 1 {
				app.SetFocus(find)
				focusIndex = 0
			} else if focusIndex == 2 {
				app.SetFocus(main)
				focusIndex = 1
			}
			return event
		case tcell.KeyHome:
			app.SetFocus(main)
			return event
			//         case tcell.KeyLeft:
			//             app.SetFocus(tags)
			//             log.Println("left key")
			//             return event
			//         case tcell.KeyRight:
			//             app.SetFocus(form)
			//             log.Println("right key")
			//             return event
		}
		return event
	})
	//     if err := tview.NewApplication().SetRoot(grid, true).EnableMouse(true).Run(); err != nil {
	if err := app.SetRoot(grid, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
