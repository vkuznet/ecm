package main

import (
	"log"

	tcell "github.com/gdamore/tcell/v2"
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

func gridView(app *tview.Application, records []VaultRecord) {
	// set tag list form
	tagList := []string{"Tag1", "tag2", "bla"}
	tags := tview.NewList()
	for _, tag := range tagList {
		tags.AddItem(tag, "", rune('-'), nil)
	}
	tags.AddItem("Quit", "press to exit", rune('q'), func() {
		app.Stop()
	})
	tags.SetBorder(true).SetTitle("Tags")

	// set main record list
	main := tview.NewList()
	for _, rec := range records {
		main.AddItem(rec.Name, "", rune('-'), nil)
	}
	main.AddItem("Quit", "Press to exit", 'q', func() {
		app.Stop()
	})
	main.SetBorder(true).SetTitle("Records")

	// set current record form view
	rec := records[0]
	name, rurl, login, password, note := rec.Details()
	form := tview.NewForm()
	form.AddInputField("Name", name, 20, nil, nil)
	form.AddInputField("Login", login, 20, nil, nil)
	form.AddPasswordField("Password", password, 10, '*', nil)
	form.AddInputField("URL", rurl, 100, nil, nil)
	form.AddInputField("Note", note, 20, nil, nil)
	form.AddButton("Save", func() {
		vrec := saveForm(form)
		log.Println("saved record", vrec)
		// TODO: update our records
	})
	form.AddButton("Quit", func() {
		app.Stop()
	})
	form.SetBorder(true).SetTitle("Form").SetTitleAlign(tview.AlignLeft)

	// construct grid view
	grid := tview.NewGrid()
	grid.SetRows(1)
	grid.SetColumns(10, 0, 50)
	grid.SetBorders(false)

	// Layout for screens wider than 100 cells.
	grid.AddItem(tags, 1, 0, 1, 1, 0, 100, false)
	grid.AddItem(main, 1, 1, 1, 1, 0, 100, true) // default focus, index 1
	grid.AddItem(form, 1, 2, 1, 1, 0, 100, false)

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
				app.SetFocus(tags)
				focusIndex = 0
			}
			return event
		case tcell.KeyCtrlB:
			if focusIndex == 0 {
				app.SetFocus(form)
				focusIndex = 2
			} else if focusIndex == 1 {
				app.SetFocus(tags)
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
