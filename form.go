package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	tcell "github.com/gdamore/tcell/v2"
	uuid "github.com/google/uuid"
	"github.com/rivo/tview"
)

// initGrid controls when we read our grid view
var initGrid bool

// helper function to start our UI app
func gpgApp(vault *Vault, interval int) {

	// create vault app and run it
	app := tview.NewApplication()

	pages := tview.NewPages()
	input, auth := authView(app, pages, vault, interval)
	pages.AddPage("auth", auth, true, true)
	go lockGPM(app, pages, input, vault, interval)

	// Start the application.
	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

// helper function to lock gpm
func lockGPM(app *tview.Application, pages *tview.Pages, input *tview.InputField, vault *Vault, interval int) {
	for {
		if time.Since(vault.Start).Seconds() > float64(interval) {
			log.Println("time to lock the screen")
			pages.HidePage("grid")
			pages.ShowPage("auth")
			pages.SwitchToPage("auth")
			vault.Start = time.Now()
			// TODO: I need to add action to prese the key in order for screen to lock
			input.SetText("")
			app.ForceDraw()
		}
		time.Sleep(1 * time.Second)
	}
}

// helper function to make secret prompt
func authView(app *tview.Application, pages *tview.Pages, vault *Vault, interval int) (*tview.InputField, *tview.Frame) {
	input := tview.NewInputField()
	input.SetFieldWidth(50).
		SetMaskCharacter('*').
		SetDoneFunc(func(key tcell.Key) {
			secret := input.GetText()
			if initGrid && secret != vault.Secret {
				log.Println("wrong password")
				return
			}
			if !initGrid {
				vault.Secret = secret
				err := vault.Read()
				if err != nil {
					log.Fatal("unable to read vault, error ", err)
				}
				log.Printf("read %d vault records", len(vault.Records))
				grid := gridView(app, pages, vault)
				pages.AddPage("grid", grid, true, true)
				initGrid = true
			}
			log.Println("switch to grid view")
			pages.HidePage("auth")
			pages.ShowPage("grid")
			pages.SwitchToPage("grid")
		})
	frame := tview.NewFrame(input)
	frame.SetBorders(10, 1, 1, 1, 10, 1)
	frame.AddText("\U0001F512 Generic Password Manager (GPM)", true, tview.AlignLeft, TitleColor)
	frame.AddText("\u00A9 2021 - github.com/vkuznet - \U0001F510", false, tview.AlignLeft, TitleColor)
	return input, frame
}

// helper function to make secret prompt
func lockView(app *tview.Application, verbose int) (string, error) {
	var password string

	input := tview.NewInputField().
		SetFieldWidth(50).
		SetMaskCharacter('*').
		SetDoneFunc(func(key tcell.Key) {
			app.Stop()
		})
	frame := tview.NewFrame(input)
	frame.SetBorders(10, 1, 1, 1, 10, 1)
	frame.AddText("\U0001F512 Generic Password Manager (GPM)", true, tview.AlignLeft, TitleColor)
	frame.AddText("\u00A9 2021 - github.com/vkuznet - \U0001F510", false, tview.AlignLeft, TitleColor)

	//     if err := app.SetRoot(input, true).EnableMouse(true).Run(); err != nil {
	if err := app.SetRoot(frame, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
	password = input.GetText()
	return password, nil
}

// helper function to provide input form, returns vault record
func inputForm(app *tview.Application) VaultRecord {
	var vrec VaultRecord
	form := tview.NewForm()
	form.AddInputField("Name", "", 100, nil, nil)
	form.AddInputField("Login", "", 100, nil, nil)
	form.AddPasswordField("Password", "", 100, '*', nil)
	form.AddInputField("URL", "", 100, nil, nil)
	form.AddButton("Save", func() {
		vrec = saveForm(form)
		app.Stop()
	})
	form.AddButton("Quit", func() {
		app.Stop()
	})
	form.SetBorder(true).SetTitle("Record Form").SetTitleAlign(tview.AlignCenter)
	if err := app.SetRoot(form, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
	return vrec
}

// helper function to save input form
func saveForm(form *tview.Form) VaultRecord {
	uid := uuid.NewString()
	rmap := make(Record)
	for i := 0; i < form.GetFormItemCount(); i++ {
		item := form.GetFormItem(i)
		key := item.GetLabel()
		val := form.GetFormItemByLabel(key).(*tview.InputField).GetText()
		rmap[key] = val
	}
	rec := VaultRecord{ID: uid, Map: rmap, ModificationTime: time.Now()}
	return rec
}

// helper function to list vault records
func listForm(list *tview.List, records []VaultRecord) *tview.List {
	list.Clear()
	for _, rec := range records {
		name, _ := rec.Map["Name"]
		list.AddItem(name, rec.ID, rune('-'), nil)
	}
	return list
}

// helper function to present recordForm
func recordForm(app *tview.Application, form *tview.Form, list *tview.List, info *tview.TextView, index int, vault *Vault) *tview.Form {
	var rec VaultRecord
	if len(vault.Records) > index {
		rec = vault.Records[index]
	}
	form.Clear(true) // clear the form
	for _, key := range rec.Keys() {
		val, _ := rec.Map[key]
		if strings.ToLower(key) == "password" {
			form.AddPasswordField(key, val, 100, '*', nil)
		} else {
			form.AddInputField(key, val, 100, nil, nil)
		}
	}
	form.SetBorder(true).SetTitle("Record form").SetTitleAlign(tview.AlignCenter)
	if len(vault.Records) == 0 {
		return form
	}
	form.AddButton("Save", func() {
		uid := rec.ID
		rmap := make(Record)
		for i := 0; i < form.GetFormItemCount(); i++ {
			item := form.GetFormItem(i)
			key := item.GetLabel()
			val := form.GetFormItemByLabel(key).(*tview.InputField).GetText()
			rmap[key] = val
		}
		rec := VaultRecord{ID: uid, Map: rmap, ModificationTime: time.Now()}
		vault.Update(rec)
		vault.Write()

		// update our list and info
		log.Println("form is updated")
		// we should update our list view too
		list.Clear()
		for _, rec := range vault.Records {
			name, _ := rec.Map["Name"]
			list.AddItem(name, rec.ID, rune('-'), nil)
		}
		// update info bar
		msg := fmt.Sprintf("Record %s is updated", uid)
		info = info.SetText(msg + helpKey())
	})
	return form
}

// helper function to build our application grid view
func gridView(app *tview.Application, pages *tview.Pages, vault *Vault) *tview.Grid {
	info := tview.NewTextView()
	list := tview.NewList()
	field := tview.NewInputField()
	find := tview.NewFrame(field)
	form := tview.NewForm()
	focusIndex := 1 // defaul focus index points to list view

	// add new frame for search bar
	input := tview.NewInputField()
	input.SetFieldWidth(50)
	input.SetDoneFunc(func(key tcell.Key) {
		pat := input.GetText()
		records := vault.Find(pat)
		msg := fmt.Sprintf("found %d records", len(records))
		if vault.Verbose > 0 {
			log.Println(msg)
		}
		if info != nil {
			info = info.SetText(msg + helpKey())
		}
		if list != nil {
			list = listForm(list, records)
		}
		// find index of record to display
		rec := records[0]
		index := 0
		for idx, r := range vault.Records {
			if r.ID == rec.ID {
				index = idx
				break
			}
		}
		// update form with proper record
		if form != nil {
			form = recordForm(app, form, list, info, index, vault)
		}
		app.SetFocus(list)
		focusIndex = 1
	})
	frame := tview.NewFrame(input)
	frame.SetBorders(2, 1, 1, 1, 10, 1)
	frame.AddText("\U0001F50D Search within the vault", true, tview.AlignLeft, TitleColor)
	find = frame

	// add search bar

	// set current record form view
	form = recordForm(app, form, list, info, 0, vault)

	// info bar
	info.SetTextAlign(tview.AlignCenter).SetText(vault.Info() + helpKey())
	info.SetDynamicColors(true)
	info.SetBorder(true).SetTitle("Info")
	info.SetTextAlign(tview.AlignLeft)
	info.SetTitleAlign(tview.AlignLeft)

	// set record list
	list = listForm(list, vault.Records)
	list.SetBorder(true).SetTitle("Record list")
	list.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		if index < len(vault.Records) {
			form = recordForm(app, form, list, info, index, vault)
			app.SetFocus(list)
			focusIndex = 1
		}
	})

	// construct grid view
	grid := tview.NewGrid()
	grid.SetBorders(false)

	// Layout for screens wider than 100 cells.
	grid.AddItem(find, 0, 0, 1, 2, 0, 0, false)
	grid.AddItem(list, 1, 0, 2, 1, 0, 0, true) // default focus, index 1
	grid.AddItem(form, 1, 1, 2, 1, 0, 0, false)
	grid.AddItem(info, 3, 0, 1, 2, 0, 0, false)

	grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		key := event.Key()
		switch key {
		case tcell.KeyCtrlA:
			idx := vault.AddRecord("login")
			list = listForm(list, vault.Records)
			list.SetCurrentItem(idx)
			app.SetFocus(form)
			focusIndex = 2
			return event
		case tcell.KeyCtrlR:
			list = listForm(list, vault.Records)
			info.SetText(helpKey())
			app.SetFocus(list)
			focusIndex = 1
		case tcell.KeyCtrlN:
			if focusIndex == 0 {
				app.SetFocus(list)
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
				app.SetFocus(list)
				focusIndex = 1
			}
			return event
		case tcell.KeyCtrlE:
			app.SetFocus(form)
			focusIndex = 2
			return event
		case tcell.KeyCtrlL:
			app.SetFocus(list)
			focusIndex = 1
			return event
		case tcell.KeyCtrlF:
			app.SetFocus(find)
			focusIndex = 0
			return event
		case tcell.KeyCtrlH:
			app.SetFocus(list)
			focusIndex = 0
			info.SetText(helpKeys())
			return event
		case tcell.KeyCtrlG:
			password := createPassword(24, true, true)
			info.SetText(password)
			return event
		case tcell.KeyCtrlP:
			app.SetFocus(form)
			copyToClipboard("Password", form, vault.Verbose)
			// return to previous view
			if focusIndex == 0 {
				app.SetFocus(find)
			} else if focusIndex == 1 {
				app.SetFocus(list)
			} else if focusIndex == 2 {
				app.SetFocus(form)
			}
			return event
		case tcell.KeyCtrlQ:
			app.Stop()
			return event
		case tcell.KeyHome:
			app.SetFocus(list)
			return event
		}
		return event
	})
	//     if err := app.SetRoot(grid, true).EnableMouse(true).Run(); err != nil {
	//         panic(err)
	//     }
	return grid
}

// helper function to copy key content from the form to clipboard
func copyToClipboard(key string, form *tview.Form, verbose int) {
	val := form.GetFormItemByLabel(key).(*tview.InputField).GetText()
	if err := clipboard.WriteAll(val); err != nil {
		log.Println("unable to copy to clipboard, error", err)
	}
	//     text, err := clipboard.ReadAll()
	//     if err != nil {
	//         log.Println("unable to read from clipboard", err)
	//     }
}
